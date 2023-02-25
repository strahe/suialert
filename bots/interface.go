package bots

type Bot interface {
	Run() error
	Close() error
}
