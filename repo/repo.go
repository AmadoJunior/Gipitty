package repo

import "context"

type IRepo interface {
	Connect(ctx context.Context, dbUri string) error
	SelectRepository(dbName string, repoName string) error
	InsertOne(data interface{}) (InsertedResult, error)
	FindOne(result interface{}, filter interface{}) (interface{}, error)
	UpdateOne(filter interface{}, data interface{}) (UpdatedResult, error)
}
