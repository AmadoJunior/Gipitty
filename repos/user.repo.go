package repos

type IUserRepo interface {
	connect(dbUri string) error
	InitRepository(dbUri string, dbName string, repoName string) error
	DeinitRepository()
	InsertUser(user interface{}) (*InsertedResult, error)
	FindUser(result interface{}, filter interface{}) error
	UpdateUser(filter interface{}, update interface{}, upsert bool) (*UpdatedResult, error)
	CreateUserIndex(key string, unique bool) (string, error)
}
