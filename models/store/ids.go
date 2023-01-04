package models

type DoID int

const (
	ADD    DoID = 1
	MODIFY DoID = 2
	DELETE DoID = 3
)

type GetID int

const (
	GET_ALL      GetID = 1
	GET_FILTERED GetID = 2
)
