package db

// PoolConfig Pool Settings
type PoolConfig struct {
	ConnMaxLifeTime int `json:"conn_max_life_time"`
	MaxIdleConn     int `json:"max_idle_conn"`
	MaxOpenConn     int `json:"max_open_conn"`
}

// MysqlConfig Mysql Settings
type MysqlConfig struct {
	Dsn                       string `json:"dsn"`
	SkipInitializeWithVersion bool   `json:"skip_initialize_with_version"`
	DefaultStringSize         uint   `json:"default_string_size"`
	DisableDatetimePrecision  bool   `json:"disable_datetime_precision"`
	DontSupportRenameIndex    bool   `json:"dont_support_rename_index"`
	DontSupportRenameColumn   bool   `json:"dont_support_rename_column"`
}

// GormConfig Gorm Settings
type GormConfig struct {
	SkipDefaultTransaction    bool `json:"skip_default_transaction"`
	DisableNestedTransaction  bool `json:"disable_nested_transaction"`
	AllowGlobalUpdate         bool `json:"allow_global_update"`
	LogLevel                  uint `json:"log_level"`
	IgnoreRecordNotFoundError bool `json:"ignore_record_not_found_error"`
	SlowThreshold             uint `json:"slow_threshold"`
}

// LoggerConfig Settings
type LoggerConfig struct {
	Level uint `json:"level"`
}

// ServerMysqlConfig Server Mysql Config
type ServerMysqlConfig struct {
	MysqlConfig MysqlConfig `json:"mysql_config"`
	GormConfig  GormConfig  `json:"gorm_config"`
	PoolConfig  PoolConfig  `json:"pool_config"`
}
