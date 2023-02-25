package config

var DefaultConfig = Config{
	Debug:      false,
	EventTypes: []string{"MoveEvent", "Publish", "CoinBalanceChange", "TransferObject", "NewObject", "EpochChange", "Checkpoint"},
}
