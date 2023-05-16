package repo

type UpdatedResult struct {
	MatchedCount  int
	ModifiedCount int
	UpsertedCount int
	UpsertedID    interface{}
}
