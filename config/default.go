package config

const (
	DevNetRpcUrl = "wss://fullnode.devnet.sui.io"
)

var DefaultConfig = Config{
	Debug: false,

	Sui: SuiConfig{
		Endpoint:   DevNetRpcUrl,
		EventTypes: []string{"MoveEvent", "Publish", "CoinBalanceChange", "TransferObject", "NewObject"},
	},

	Database: DatabaseConfig{
		Postgres: PostgresConfig{
			SchemaName: "public",
			PoolSize:   50,
		},
	},
}
