package repos

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		return err
	}

	ur.client = *mongoClient

	return nil
}

func (ur *UserRepoImpl) InitRepository(dbUri string, dbName string, repoName string) error {
	err := ur.connect(dbUri)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUserRepoInit, err)
	}
	ur.store = *ur.client.Database(dbName).Collection(repoName)
	return nil
}

func (ur UserRepoImpl) DeinitRepository() error {
	err := ur.client.Disconnect(ur.ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUserRepoDeinit, err)
	}
	return nil
}

func (ur UserRepoImpl) CreateNewUser(user *models.SignUpInput) (string, error) {
	insertResult, err := ur.store.InsertOne(ur.ctx, &user)

	//Catch Errs
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", fmt.Errorf("%w: %w", ErrDuplicateEmail, err)
		}
		return "", fmt.Errorf("%w: %w", ErrUserInsertion, err)
	}

	// Assert InsertedID to ObjectID
	idObj, isObjID := insertResult.InsertedID.(primitive.ObjectID)
	if !isObjID {
		return "", ErrUserIDAssertion
	}

	return idObj.Hex(), nil
}

func (ur UserRepoImpl) FindUserByID(id string) (*models.DBResponse, error) {
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidIDHex, err)
	}

	user := &models.DBResponse{}
	filter := bson.M{"_id": objID}
	err = ur.store.FindOne(ur.ctx, filter).Decode(user)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	return user, nil
}

func (ur UserRepoImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	user := &models.DBResponse{}
	filter := bson.M{"email": strings.ToLower(email)}
	err := ur.store.FindOne(ur.ctx, filter).Decode(user)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	return user, nil
}

func (ur UserRepoImpl) UpdateUserById(id string, update *models.UpdateInput) error {
	// Convert String to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidIDHex, err)
	}

	doc, err := utils.ToDoc(update)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidUpdateInput, err)
	}

	filter := bson.M{"_id": objectID}
	updatePrimitive := bson.D{{Key: "$set", Value: doc}}
	res, err := ur.store.UpdateOne(ur.ctx, filter, updatePrimitive)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrUserUpdate, err)
	}

	if res.MatchedCount < 1 {
		return fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	return nil
}

func (ur UserRepoImpl) UpdateUserByEmail(email string, update *models.UpdateInput) error {
	doc, err := utils.ToDoc(update)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidUpdateInput, err)
	}

	filter := bson.M{"email": email}
	updatePrimitive := bson.D{{Key: "$set", Value: doc}}
	res, err := ur.store.UpdateOne(ur.ctx, filter, updatePrimitive)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrUserUpdate, err)
	}

	if res.MatchedCount < 1 {
		return fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	return nil
}

func (ur UserRepoImpl) VerifyUserEmail(verificationCode string) error {
	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}

	res, err := ur.store.UpdateOne(ur.ctx, query, update)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrUserVerification, err)
	}

	if res.MatchedCount < 1 {
		return fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	return nil
}

func (ur UserRepoImpl) StorePasswordResetToken(userEmail string, passwordResetToken string) error {
	query := bson.D{{Key: "email", Value: strings.ToLower(userEmail)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	res, err := ur.store.UpdateOne(ur.ctx, query, update)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrStorePasswordResetToken, err)
	}

	if res.MatchedCount < 1 {
		return fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	return nil
}

func (ur UserRepoImpl) ResetUserPassword(passwordResetToken string, newPassword string) error {
	query := bson.D{{Key: "passwordResetToken", Value: passwordResetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: newPassword}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	res, err := ur.store.UpdateOne(ur.ctx, query, update)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrResetPassword, err)
	}

	if res.MatchedCount < 1 {
		return fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	return nil
}
