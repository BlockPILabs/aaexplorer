package entity

import (
	"context"
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
	Client *ent.Client
	Config *config.DbConfig
}

var clients = &sync.Map{}

func Start(cfg *config.Config) error {
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
		}
		if database.Schema != nil {
			opts = append(opts, ent.AlternateSchema(*database.Schema))
		}
		// connect
		c := ent.NewClient(opts...)
		if err != nil {
			return err
		}
		if database.Debug {
			c = c.Debug()
		}
		_c := client{
			Client: c,
			Config: database,
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
		return nil, errors.New(fmt.Sprintf("group error %s", g))
	}
	return cc.Client, nil
}
