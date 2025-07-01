package handler

import (
	"papslab/item"
	"papslab/user"
)

type Storage interface {
	InsertItem(i item.Item) error
	SelectAllItems() ([]item.Item, error)
	SelectAnyItems(i item.Item) ([]item.Item, error)
	DeleteItem(virtualID int) error
	InsertUser(in *user.User) error
	CheckUser(in *user.User) (exists bool, isPriv bool, err error)
	IsLoginAvailable(login string) (bool, error)
}
