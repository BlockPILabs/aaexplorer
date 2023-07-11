package config

import (
	"entgo.io/ent/dialect"
	"fmt"
	"github.com/BlockPILabs/aa-scan/version"
	"strconv"
	"strings"
	"time"
)

type DbConfig struct {
	Group           string `mapstructure:"group" toml:"group"`
	Type            string `mapstructure:"type" toml:"type"`
	Host            string `mapstructure:"host" toml:"host"`
	Port            int    `mapstructure:"port" toml:"port"`
	User            string `mapstructure:"user" toml:"user"`
	Pass            string `mapstructure:"pass" toml:"pass"`
	Name            string `mapstructure:"name" toml:"name"`
	ApplicationName string `mapstructure:"applicationName" toml:"applicationName"`
	MaxIdleCount    int    `mapstructure:"maxIdleCount" toml:"maxIdleCount"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns" toml:"maxOpenConns"`
	MaxLifetime     int64  `mapstructure:"maxLifetime" toml:"maxLifetime"`
	Debug           bool   `mapstructure:"debug" toml:"debug"`
}

func DefaultDatabaseConfig() []*DbConfig {
	return []*DbConfig{
		{
			Group:           Default,
			Type:            dialect.Postgres,
			Host:            "127.0.0.1",
			Port:            5432,
			User:            "postgres",
			Pass:            "123456",
			Name:            "postgres",
			ApplicationName: version.Name,
			MaxIdleCount:    50,
			MaxOpenConns:    100,
			MaxLifetime:     int64(time.Hour),
		},
	}
}

func (cfg DbConfig) BuildDsn() (string, error) {
	//	postgresql://[user[:password]@][netloc][:port][/dbname][?params]
	switch cfg.Type {
	case dialect.Postgres:
		return cfg.buildPostgresqlDsn()
	default:
		return "", fmt.Errorf("unsupported driver: %q", cfg.Type)
	}

}

func (cfg DbConfig) buildPostgresqlDsn() (string, error) {
	//	postgresql://[user[:password]@][netloc][:port][/dbname][?params]
	dsn := strings.Builder{}
	dsn.WriteString("postgresql://")
	if len(cfg.User) > 0 {
		dsn.WriteString(cfg.User)
		dsn.WriteByte(':')
		dsn.WriteString(cfg.Pass)
		dsn.WriteByte('@')
	}
	if len(cfg.Host) > 0 {
		dsn.WriteString(cfg.Host)
		dsn.WriteByte(':')
		dsn.WriteString(strconv.Itoa(cfg.Port))
	}
	if len(cfg.Name) > 0 {
		dsn.WriteByte('/')
		dsn.WriteString(cfg.Name)
	}
	if len(cfg.ApplicationName) > 0 {
		dsn.WriteByte('?')
		dsn.WriteString("application_name=")
		dsn.WriteString(cfg.ApplicationName)
	}
	return dsn.String(), nil
}

func (cfg DbConfig) ValidateBasic() error {
	_, err := cfg.BuildDsn()
	return err
}
