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
	_ "github.com/lib/pq"              // postgres driver
	_ "github.com/mattn/go-sqlite3"    // sqlite3
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
