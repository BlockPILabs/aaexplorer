package entity

import (
	"context"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgres driver
	_ "github.com/mattn/go-sqlite3"    // sqlite3
)

type client struct {
	Client  *ent.Client
	Config  *config.DbConfig
	Dialect string
	Network *ent.Network
	logger  log.Logger
	lck     sync.Mutex
}

func (c *client) start() error {
	if c.Client != nil {
		return nil
	}

	c.lck.Lock()
	defer c.lck.Unlock()
	if c.Client != nil {
		return nil
	}

	database := c.Config
	logger := c.logger
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
				logger.Debug(log.MaskMsg(msg))
			} else if len(a) > 1 {
				msg := fmt.Sprint(a[0])
				logger.Debug(log.MaskMsg(msg), "args", a[1:])
			}
		}))
	}

	if database.Schema != nil {
		opts = append(opts, ent.AlternateSchema(*database.Schema))
	}
	// connect
	entClient := ent.NewClient(opts...)
	if err != nil {
		return err
	}
	c.Client = entClient
	return nil
}

func loadAllNetworksClients(ctx context.Context, logger log.Logger) (err error) {
	entClient, err := Client(ctx)
	if err != nil {
		return err
	}

	networks, err := entClient.Network.Query().All(ctx)
	if err != nil {
		return err
	}

	for _, network := range networks {
		err = networkClientInit(ctx, logger, network)
	}
	return nil
}

func networkClientInit(ctx context.Context, logger log.Logger, network *ent.Network) (err error) {

	database := &config.DbConfig{}

	if network.DbConfig == nil {
		return
	}

	err = network.DbConfig.AssignTo(database)
	if err != nil {
		logger.Warn("network db config error", "network", network.ID, "err", err)
		return
	}

	err = database.ValidateBasic()
	if err != nil {
		logger.Warn("network db config Validate error", "network", network.ID, "err", err)
		return
	}

	c, ok := clients.Load(network.ID)
	if ok {
		_cc, ok := c.(*client)
		if ok {
			if _cc.Config.Host == database.Host &&
				_cc.Config.Port == database.Port &&
				_cc.Config.Name == database.Name &&
				_cc.Config.User == database.User &&
				_cc.Config.Pass == database.Pass {
				_cc.Network = network
				return
			}
			defer func() {
				if _cc.Client != nil {
					logger.Warn("close db config", "network", _cc.Network.ID)
					_cc.Client.Close()
				}
			}()
		}
	}

	cc := &client{
		logger:  logger,
		Config:  database,
		Dialect: database.Type,
		Network: network,
	}
	if database.AutoStart {
		err = cc.start()
		if err != nil {
			logger.Warn("network db config start error", "network", network.ID, "err", err)
			return err
		}
	}

	clients.Store(network.ID, cc)
	return nil
}
