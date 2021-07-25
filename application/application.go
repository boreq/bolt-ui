package application

type Application struct {
	Browse *BrowseHandler
}

type TransactionProvider interface {
	Read(handler TransactionHandler) error
	Write(handler TransactionHandler) error
}

type TransactionHandler func(adapters *TransactableAdapters) error

type TransactableAdapters struct {
}
