package model

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Rule{},
		&CoinBalanceChangeEvent{},
		&DeleteObjectEvent{},
		&MoveEvent{},
		&MutateObjectEvent{},
		&NewObjectEvent{},
		&PublishEvent{},
		&TransferObjectEvent{},
	)
}

type Model interface {
	TableName() string
}
