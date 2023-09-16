package entity

import (
	"context"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"strings"
	"time"

	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"sync"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgres driver
	_ "github.com/mattn/go-sqlite3"    // sqlite3
)

type client struct {
	Client  *ent.Client
	Config  *config.DbConfig
	Dialect string
}

var clients = &sync.Map{}

func Start(logger log.Logger, cfg *config.Config) error {
	for i, database := range cfg.Databases {
		dsn, err := database.BuildDsn()
		if err != nil {
			return err
		}

		drv, err := entsql.Open(database.Type, dsn)
		if err != nil {
			return err
		}
		if database.MaxIdleCount > 0 {
			drv.DB().SetMaxIdleConns(database.MaxIdleCount)
		}
		if database.MaxOpenConns > 0 {
			drv.DB().SetMaxOpenConns(database.MaxOpenConns)
		}
		if database.MaxLifetime > 0 {
			leftTime := time.Duration(database.MaxLifetime)
			if leftTime < time.Millisecond {
				leftTime *= time.Second
			}
			drv.DB().SetConnMaxLifetime(leftTime)
		}

		opts := []ent.Option{
			ent.Driver(drv),
			ent.Log(func(a ...any) {

			}),
		}
		if database.Debug {
			opts = append(opts, ent.Debug(), ent.Log(func(a ...any) {
				if len(a) == 1 {
					msg := fmt.Sprint(a[0])
					logger.Debug(msg)
				} else if len(a) > 1 {
					msg := fmt.Sprint(a[0])
					logger.Debug(msg, "args", a[1:])
				}
			}))
		}

		if database.Schema != nil {
			opts = append(opts, ent.AlternateSchema(*database.Schema))
		}
		// connect
		c := ent.NewClient(opts...)
		if err != nil {
			return err
		}

		_c := &client{
			Client:  c,
			Config:  database,
			Dialect: database.Type,
		}
		if i == 0 {
			clients.Store(config.Default, _c)
		}
		if len(database.Group) > 0 {
			clients.Store(database.Group, _c)
		}

	}
	return nil
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
