package repo

import "context"

type IRepo interface {
	Connect(ctx context.Context, dbUri string) error
	SelectRepository(dbName string, repoName string) error
	InsertOne(ctx context.Context, data interface{}) error
	FindOne(ctx context.Context, filter interface{}) (*interface{}, error)
	UpdateOne(ctx context.Context, filter interface{}, data interface{}) error
}
