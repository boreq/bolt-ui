package application

import "errors"

type Key struct {
	b []byte
}

func NewKey(b []byte) (Key, error) {
	if len(b) == 0 {
		return Key{}, errors.New("key can not be empty")
	}

	tmp := make([]byte, len(b))
	copy(tmp, b)
	return Key{tmp}, nil
}

func MustNewKey(b []byte) Key {
	v, err := NewKey(b)
	if err != nil {
		panic(err)
	}
	return v
}

func (k Key) Bytes() []byte {
	tmp := make([]byte, len(k.b))
	copy(tmp, k.b)
	return tmp
}

type Value struct {
	b []byte
}

func NewValue(b []byte) (Value, error) {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return Value{tmp}, nil
}

func MustNewValue(b []byte) Value {
	v, err := NewValue(b)
	if err != nil {
		panic(err)
	}
	return v
}

type Database interface {
	Browse(path []Key, before *Key, after *Key) ([]Entry, error)
}

type Entry struct {
	Key   Key
	Value Value
}

type Application struct {
	Browse *BrowseHandler
}

type TransactionProvider interface {
	Read(handler TransactionHandler) error
	Write(handler TransactionHandler) error
}

type TransactionHandler func(adapters *TransactableAdapters) error

type TransactableAdapters struct {
	Database Database
}
