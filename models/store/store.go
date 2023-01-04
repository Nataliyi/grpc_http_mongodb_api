package models

type IStore interface {
	DoOne(DoID, IStoreDoRequest) error
	Get(GetID, interface{}) ([]IStoreGetResponse, error)
}
