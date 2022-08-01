package file

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	v *viper.Viper
}

// PoolConfig Pool Settings
type PoolConfig struct {
	ConnMaxLifetime int `json:"conn_max_lifetime"`
	MaxIdleConn     int `json:"max_idle_conn"`
	MaxOpenConn     int `json:"max_open_conn"`
}

// MysqlConfig Mysql Settings
type MysqlConfig struct {
	Dsn                       string `json:"'dsn'"`
	SkipInitializeWithVersion bool   `json:"skip_initialize_with_version"`
	DefaultStringSize         uint   `json:"default_string_size"`
	DisableDatetimePrecision  bool   `json:"disable_datetime_precision"`
	DontSupportRenameIndex    bool   `json:"dont_support_rename_index"`
	DontSupportRenameColumn   bool   `json:"dont_support_rename_column"`
}

// GormConfig Gorm Settings
type GormConfig struct {
	SkipDefaultTransaction   bool `json:"skip_default_transaction"`
	DisableNestedTransaction bool `json:"disable_nested_transaction"`
	AllowGlobalUpdate        bool `json:"allow_global_update"`
}

// ServerMysqlConfig Server Mysql Config
type ServerMysqlConfig struct {
	MysqlConfig MysqlConfig `json:"mysql_config"`
	GormConfig  GormConfig  `json:"gorm_config"`
	PoolConfig  PoolConfig  `json:"pool_config"`
}

func NewConfig() *Config {
	
	v := viper.New()
	
	v.SetConfigFile("./config/database.json")
	
	return &Config{
		v: v,
	}
}

// Get Connect Config
func (c *Config) Get(clusterName string) (res *ConnectConfig, err error) {
	
	var config *ServerMysqlConfig
	
	// Load Config
	if err = c.v.ReadInConfig(); nil != err {
		
		return nil, err
		
	}
	
	// Unmarshal To Struct .
	if err = c.v.UnmarshalKey(clusterName, &config); nil != err {
		
		return nil, err
		
	}
	
	return &ConnectConfig{
		MysqlConfig: mysql.Config{
			DSN:                       config.MysqlConfig.Dsn,
			SkipInitializeWithVersion: config.MysqlConfig.SkipInitializeWithVersion,
			DefaultStringSize:         config.MysqlConfig.DefaultStringSize,
			DisableDatetimePrecision:  config.MysqlConfig.DisableDatetimePrecision,
			DontSupportRenameIndex:    config.MysqlConfig.DontSupportRenameIndex,
			DontSupportRenameColumn:   config.MysqlConfig.DontSupportRenameColumn,
		},
		GormConfig: &gorm.Config{
			SkipDefaultTransaction:   config.GormConfig.SkipDefaultTransaction,
			DisableNestedTransaction: config.GormConfig.DisableNestedTransaction,
			AllowGlobalUpdate:        config.GormConfig.AllowGlobalUpdate,
		},
		MysqlPoolConfig: PoolConfig{
			ConnMaxLifetime: config.PoolConfig.ConnMaxLifetime,
			MaxIdleConn:     config.PoolConfig.MaxIdleConn,
			MaxOpenConn:     config.PoolConfig.MaxOpenConn,
		},
	}, nil
}
