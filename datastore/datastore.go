package datastore

import "context"

// Repository is an interface for the storage implementation of the auth service
type Repository interface {
	Create(ctx context.Context, user *User) error
	//DeleteUser(ctx context.Context, email string, clientId string) error
}
