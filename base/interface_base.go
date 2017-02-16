package base

import (
//"errors"
)

type PconnBase interface{}

type Hub interface {
	AddPconn(cuid string, c *PconnBase) error
	DelPconn(cuid string) error
	GetPconn(cuid string) *PconnBase
	IsPconnExist(cuid string) bool
}
