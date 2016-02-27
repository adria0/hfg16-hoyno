package main

import (
	"errors"
	"github.com/boltdb/bolt"
)

var (
	ErrNotExists = errors.New("Entry does not exist")
)

var (
	UserNameKey    = []byte{1}
	PublicNameKey  = []byte{2}
	EmailKey       = []byte{3}
	GroupNameKey   = []byte{4}
	GroupEmailsKey = []byte{5}
	ChatStatusKey  = []byte{6}

	ChatOFF byte = 0
	ChatON  byte = 1
)

type User struct {
	ID          string
	UserName    string
	PublicName  string
	Email       string
	GroupName   string
	GroupEmails string
	ChatStatus  bool
}

func save(ID string, user User) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(ID))
		if err != nil {
			return err

		}
        b.Put(UserNameKey, []byte(user.UserName))
		b.Put(PublicNameKey, []byte(user.PublicName))
		b.Put(EmailKey, []byte(user.Email))
		b.Put(GroupNameKey, []byte(user.GroupName))
		b.Put(GroupEmailsKey, []byte(user.GroupEmails))
		return nil

	})
}

func load(ID string) (User, error) {
	var user User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ID))
		if b == nil {
			return ErrNotExists
		}
		chatStatus := b.Get(ChatStatusKey)
		user = User{
			ID:          ID,
			UserName:    string(b.Get(UserNameKey)),
			PublicName:  string(b.Get(PublicNameKey)),
			Email:       string(b.Get(EmailKey)),
			GroupName:   string(b.Get(GroupNameKey)),
			GroupEmails: string(b.Get(GroupEmailsKey)),
			ChatStatus:  chatStatus[0] == ChatON,
		}
		return nil
	})
	return user, err
}

func setChatStatus(ID string, status bool) error {
	var isActive []byte
	if status {
		isActive = []byte{1}
	} else {
		isActive = []byte{0}
	}

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(ID))
		if err != nil {
			return nil
		}
		b.Put(ChatStatusKey, isActive)
		return nil
	})
}
