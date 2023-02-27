package config

type Config struct {
	// Enable Debug model
	Debug bool `yaml:"debug" json:"debug" mapstructure:"debug"`

	// Event type to subscribe
	EventTypes []string `yaml:"event_types" json:"event_types" mapstructure:"event_types"`

	Bots BotsConfig `yaml:"bots" json:"bots" mapstructure:"bots"`

	Database DatabaseConfig `yaml:"database" json:"database" mapstructure:"database"`
}

type BotsConfig struct {
	Discord DiscordBotConfig `yaml:"discord" json:"discord" mapstructure:"discord"`
}

type DiscordBotConfig struct {
	Enable bool   `yaml:"enable" json:"enable" mapstructure:"enable"`
	AppID  string `yaml:"app_id" json:"app_id" mapstructure:"app_id"`
	Token  string `yaml:"token" json:"token" mapstructure:"token"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres" json:"postgres" mapstructure:"postgres"`
}

type PostgresConfig struct {
	URL        string `yaml:"url" json:"url" mapstructure:"url"`
	PoolSize   int    `yaml:"pool_size" json:"pool_size" mapstructure:"pool_size"`
	SchemaName string `yaml:"schema_name" json:"schema_name" mapstructure:"schema_name"`
}
