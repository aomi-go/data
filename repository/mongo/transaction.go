package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	instance                     *Transaction
	ErrTransactionNotInitialized = errors.New("transaction not initialized")
)

func InitTransaction(client *mongo.Client) *Transaction {
	instance = &Transaction{client: client}
	return instance
}

type Transaction struct {
	client *mongo.Client
}

// WithTransaction 执行一个事务
func WithTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	if instance == nil {
		return nil, ErrTransactionNotInitialized
	}

	session, err := instance.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, fn)
}
