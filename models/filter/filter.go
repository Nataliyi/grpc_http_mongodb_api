package models

type IFilter interface {
	Filter(map[string][]interface{}) interface{}
}
