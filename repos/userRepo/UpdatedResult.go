package userRepo

type UpdatedResult struct {
	MatchedCount   int
	ModifiedCount  int
	UpsertedCount  int
	UpsertedUserID interface{}
}
