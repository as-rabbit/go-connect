package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Configer interface {
	Get(clusterName string) ClusterConfig
}

type defaultConfig struct {
}

type PoolConfig struct {
	ConnMaxLifetime int `json:"conn_max_lifetime"`
	MaxIdleConn     int `json:"max_idle_conn"`
	MaxOpenConn     int `json:"max_open_conn"`
}

type ClusterConfig struct {
	MysqlConfig     mysql.Config `json:"mysql_config"`
	GormConfig      *gorm.Config `json:"gorm_config"`
	MysqlPoolConfig PoolConfig   `json:"mysql_pool_config"`
}

func NewConfig() (d *defaultConfig, err error) {
	return &defaultConfig{}, nil
}

func (d *defaultConfig) Get(clusterName string) (config *ClusterConfig, err error) {

	return &ClusterConfig{
		MysqlConfig: mysql.Config{
			DSN:                       "",
			SkipInitializeWithVersion: true,
			DefaultStringSize:         256,  // string 类型字段的默认长度
			DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		},
		GormConfig: &gorm.Config{
			SkipDefaultTransaction:   true,  // 表，列命名策略 需要实现Namer 默认实现schema.NamingStrategy
			DisableNestedTransaction: true,  // 禁止嵌套事务// 全局更新
			AllowGlobalUpdate:        false, // 全局更新
		},
		MysqlPoolConfig: PoolConfig{
			ConnMaxLifetime: 600,
			MaxIdleConn:     10,
			MaxOpenConn:     50,
		},
	}, nil
}
