package entity

import (
	"context"
	"entgo.io/ent/dialect"
	"errors"
	"fmt"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/google/uuid"
	_ "github.com/lib/pq"           // postgres driver
	_ "github.com/mattn/go-sqlite3" // sqlite3
	"github.com/procyon-projects/chrono"
	"strings"
	"sync"
	"time"
)

var clients = &sync.Map{}

var clientTick = &sync.Once{}

func Start(logger log.Logger, cfg *config.Config) error {
	for i, database := range cfg.Databases {
		c := &client{
			logger:  logger,
			Config:  database,
			Dialect: database.Type,
		}
		if database.AutoStart {
			err := c.start()
			if err != nil {
				return err
			}
		}
		if i == 0 {
			clients.Store(config.Default, c)
		}
		if len(database.Group) > 0 {
			clients.Store(database.Group, c)
		}
	}
	//ctx := context.Background()
	//err := loadAllNetworksClients(ctx, logger)
	//if err != nil {
	//	return err
	//}
	clientTick.Do(func() {
		loadAllNetworksClients(context.Background(), logger)
		chrono.NewDefaultTaskScheduler().ScheduleWithFixedDelay(func(ctx context.Context) {
			loadAllNetworksClients(ctx, logger)
		}, time.Minute, chrono.WithTime(time.Now().Add(time.Minute*5)))
	})

	return nil
}

func NetworkClient(ctx context.Context, network *ent.Network) (*ent.Client, error) {
	_ = networkClientInit(ctx, log.Context(ctx), network)
	return Client(ctx, network.ID)
}
func Client(ctx context.Context, group ...string) (*ent.Client, error) {
	g := config.Default
	if len(group) > 0 && len(group[0]) > 0 {
		g = group[0]
	}
	c, ok := clients.Load(g)
	if !ok {
		if li := strings.LastIndexByte(g, ':'); li >= 0 {
			return Client(ctx, g[0:li])
		}
		if g != config.Default {
			return Client(ctx)
		}
		log.Context(ctx).Error("not found group")
		return nil, errors.New(fmt.Sprintf("not found group %s", g))
	}
	cc, ok := c.(*client)
	if !ok {
		log.Context(ctx).Error("group error")
		return nil, errors.New(fmt.Sprintf("group connect error %s", g))
	}
	err := cc.start()
	if err != nil {
		return nil, err
	}
	return cc.Client, nil
}

func MustClient(group ...string) *ent.Client {
	c, err := Client(context.Background(), group...)
	if err != nil {
		panic(err)
		return nil
	}
	return c
}
func SetDialect(ctx context.Context, f func(string), group ...string) {
	g := config.Default
	if len(group) > 0 && len(group[0]) > 0 {
		g = group[0]
	}
	c, ok := clients.Load(g)
	if !ok {
		if li := strings.LastIndexByte(g, ':'); li >= 0 {
			SetDialect(ctx, f, g[0:li])
			return
		}
		if g != config.Default {
			SetDialect(ctx, f)
			return
		}
		log.Context(ctx).Error("not found group")
		f(dialect.Postgres)
		return
	}
	cc, ok := c.(*client)
	if !ok {
		log.Context(ctx).Error("group error")
		f(dialect.Postgres)
		return
	}
	f(cc.Dialect)
}

var TxErr = errors.New("ent: cannot start a transaction within a transaction")

func WithTx(ctx context.Context, client *ent.Client, fn func(db *ent.Client) error) error {

	tx, err := client.Tx(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "cannot start a transaction within a transaction") {
			return withSubTx(ctx, client, fn)
		}
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx.Client()); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func withSubTx(ctx context.Context, client *ent.Client, fn func(db *ent.Client) error) error {

	pointId := "p_" + strings.Replace(uuid.NewString(), "-", "", -1)

	_, err := client.ExecContext(ctx, fmt.Sprintf("SAVEPOINT %s", pointId))
	if err != nil {
		return err
	}

	defer func() {
		if v := recover(); v != nil {
			client.ExecContext(ctx, fmt.Sprintf("ROLLBACK TO  %s", pointId))
			panic(v)
		}
	}()
	if err := fn(client); err != nil {
		if _, rerr := client.ExecContext(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", pointId)); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if _, rerr := client.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT  %s", pointId)); rerr != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
