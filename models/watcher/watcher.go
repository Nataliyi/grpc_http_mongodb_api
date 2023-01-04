package models

type WatcherChannel chan string

type IWatcher interface {
	Listen()
	Close()
	GetChannel() WatcherChannel
}
