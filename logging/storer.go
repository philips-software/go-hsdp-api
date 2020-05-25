package logging

// Storer defines the store operations for logging
type Storer interface {
	StoreResources(msgs []Resource, count int) (*StoreResponse, error)
}

var _ Storer = &Client{}
