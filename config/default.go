package config

const (
	DevNetRpcUrl = "wss://fullnode.devnet.sui.io"
)

var DefaultConfig = Config{
	Debug: false,

	Sui: SuiConfig{
		Endpoint: DevNetRpcUrl,
	},

	Database: DatabaseConfig{
		Driver: "sqlite3",
		DSN:    "db.sqlite3",
	},
}
