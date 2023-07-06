package entity

import (
	entsql "entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"time"

	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq" // postgres driver
	_ "github.com/mattn/go-sqlite3"
)

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
		// 从db变量中构造一个ent.Driver对象。
		client := ent.NewClient(ent.Driver(drv))
		if err != nil {
			return err
		}

		if i == 0 {
			clients.Store(config.Default, client)
		}
		if len(database.Group) > 0 {
			clients.Store(database.Group, client)
		}

	}
	return nil
}

func Client(group ...string) (*ent.Client, error) {
	g := config.Default
	if len(group) > 0 && len(group[0]) > 0 {
		g = group[0]
	}
	c, ok := clients.Load(g)
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found group %s", g))
	}
	client, ok := c.(*ent.Client)
	if !ok {
		return nil, errors.New(fmt.Sprintf("group error %s", g))
	}
	return client, nil
}
