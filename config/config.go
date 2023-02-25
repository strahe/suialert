package config

type Config struct {
	// Enable Debug model
	Debug bool `yaml:"debug" json:"debug" mapstructure:"debug"`

	// Event type to subscribe
	EventTypes []string `yaml:"event_types" json:"event_types" mapstructure:"event_types"`

	Bots BotsConfig `yaml:"bots" json:"bots" mapstructure:"bots"`
}

type BotsConfig struct {
	Discord DiscordBotConfig `yaml:"discord" json:"discord" mapstructure:"discord"`
}

type DiscordBotConfig struct {
	Enable bool   `yaml:"enable" json:"enable" mapstructure:"enable"`
	AppID  string `yaml:"app_id" json:"app_id" mapstructure:"app_id"`
	Token  string `yaml:"token" json:"token" mapstructure:"token"`
}
