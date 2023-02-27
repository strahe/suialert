package config

var DefaultConfig = Config{
	Debug:      false,
	EventTypes: []string{"MoveEvent", "Publish", "CoinBalanceChange", "TransferObject", "NewObject", "EpochChange", "Checkpoint"},

	Database: DatabaseConfig{
		Postgres: PostgresConfig{
			SchemaName: "public",
			PoolSize:   50,
		},
	},
}
