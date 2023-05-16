package repo

import (
	"context"
	"fmt"

	"github.com/AmadoJunior/Gipitty/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepoImpl struct {
	ctx    context.Context
	client mongo.Client
	repo   mongo.Collection
}

func NewMongoRepo() *MongoRepoImpl {
	return &MongoRepoImpl{}
}

func (r *MongoRepoImpl) Connect(ctx context.Context, dbUri string) error {
	//Connect to MongoDB
	mongoConn := options.Client().ApplyURI(dbUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)

	if err != nil {
		return err
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	r.client = *mongoClient

	fmt.Println("MongoDB Successfully Connected...")
	return nil
}

func (r *MongoRepoImpl) SelectRepository(dbName string, repoName string) error {
	r.repo = *r.client.Database(dbName).Collection(repoName)
	return nil
}

func (r MongoRepoImpl) InsertOne(data interface{}) (InsertedResult, error) {
	result, err := r.repo.InsertOne(r.ctx, &data)

	if err != nil {
		return InsertedResult{}, err
	}

	return InsertedResult{InsertedID: result.InsertedID}, nil
}

func (r MongoRepoImpl) FindOne(result *models.DBResponse, filter primitive.M) (*models.DBResponse, error) {
	//var result *models.DBResponse

	err := r.repo.FindOne(r.ctx, filter).Decode(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r MongoRepoImpl) UpdateOne(filter primitive.M, update primitive.M) (UpdatedResult, error) {
	res, err := r.repo.UpdateOne(r.ctx, filter, update)

	if err != nil {
		return UpdatedResult{}, err
	}

	return UpdatedResult{
		MatchedCount:  int(res.MatchedCount),
		ModifiedCount: int(res.ModifiedCount),
		UpsertedCount: int(res.UpsertedCount),
		UpsertedID:    res.UpsertedID}, nil
}
