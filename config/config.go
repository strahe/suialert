package config

type Config struct {
	// Enable Debug model
	Debug bool `yaml:"debug" json:"debug" mapstructure:"debug"`

	Sui SuiConfig `yaml:"sui" json:"sui" mapstructure:"sui"`

	Bots BotsConfig `yaml:"bots" json:"bots" mapstructure:"bots"`

	Database DatabaseConfig `yaml:"database" json:"database" mapstructure:"database"`
}

type SuiConfig struct {
	Endpoint string `yaml:"endpoint" json:"endpoint" mapstructure:"endpoint"`
	// Event type to subscribe
	EventTypes []string `yaml:"event_types" json:"event_types" mapstructure:"event_types"`
}

type BotsConfig struct {
	Discord DiscordBotConfig `yaml:"discord" json:"discord" mapstructure:"discord"`
}

type DiscordBotConfig struct {
	Enable bool   `yaml:"enable" json:"enable" mapstructure:"enable"`
	AppID  string `yaml:"app_id" json:"app_id" mapstructure:"app_id"`
	Token  string `yaml:"token" json:"token" mapstructure:"token"`
}

// DatabaseConfig
// https://gorm.io/docs/connecting_to_the_database.html
type DatabaseConfig struct {
	// Driver name, supported: mysql, sqlite3, postgres
	Driver string `yaml:"driver" json:"driver" mapstructure:"driver"`
	// Database connection string
	DSN string `yaml:"dsn" json:"dsn" mapstructure:"dsn"`
}
