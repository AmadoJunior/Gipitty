package userRepo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type UserRepoImpl struct {
	ctx    context.Context
	client mongo.Client
	store  mongo.Collection
}

func NewUserRepo(ctx context.Context) *UserRepoImpl {
	return &UserRepoImpl{ctx: ctx}
}

func (ur *UserRepoImpl) connect(dbUri string) error {
	//Connect to MongoDB
	mongoConn := options.Client().ApplyURI(dbUri)
	mongoClient, err := mongo.Connect(ur.ctx, mongoConn)

	if err != nil {
		return errors.New("failed to connect to user repo")
	}

	if err := mongoClient.Ping(ur.ctx, readpref.Primary()); err != nil {
		return errors.New("failed to ping user repo")
	}

	ur.client = *mongoClient

	return nil
}

func (ur *UserRepoImpl) InitRepository(dbUri string, dbName string, repoName string) error {
	err := ur.connect(dbUri)
	if err != nil {
		return errors.New("failed to initiate user repo")
	}
	ur.store = *ur.client.Database(dbName).Collection(repoName)
	return nil
}

func (ur UserRepoImpl) DeinitRepository() {
	ur.client.Disconnect(ur.ctx)
}

func (ur UserRepoImpl) InsertUser(user interface{}) (*InsertedResult, error) {
	result, err := ur.store.InsertOne(ur.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, errors.New("failed to insert user")
	}

	return &InsertedResult{InsertedUserID: result.InsertedID}, nil
}

func (ur UserRepoImpl) FindUser(result interface{}, filter interface{}) error {
	err := ur.store.FindOne(ur.ctx, filter).Decode(result)

	if err != nil {
		return errors.New("failed to find user")
	}

	return nil
}

func (ur UserRepoImpl) UpdateUser(filter interface{}, update interface{}, upsert bool) (*UpdatedResult, error) {
	res, err := ur.store.UpdateOne(ur.ctx, filter, update)

	if err != nil {
		return nil, errors.New("failed to update user")
	}

	return &UpdatedResult{
		MatchedCount:   int(res.MatchedCount),
		ModifiedCount:  int(res.ModifiedCount),
		UpsertedCount:  int(res.UpsertedCount),
		UpsertedUserID: res.UpsertedID}, nil
}

func (ur UserRepoImpl) CreateUserIndex(key string, unique bool) (string, error) {
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{key: 1}, Options: opt}
	_, err := ur.store.Indexes().CreateOne(ur.ctx, index)

	if err != nil {
		return "", errors.New("failed to create user index: " + key)
	}

	return key, nil
}
