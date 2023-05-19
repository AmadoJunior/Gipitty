package repos

import (
	"context"
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
		return utils.GenerateError(ErrUserRepoInit, err)
	}
	ur.store = *ur.client.Database(dbName).Collection(repoName)
	return nil
}

func (ur UserRepoImpl) DeinitRepository() error {
	err := ur.client.Disconnect(ur.ctx)
	if err != nil {
		return utils.GenerateError(ErrUserRepoDeinit, err)
	}
	return nil
}

func (ur UserRepoImpl) CreateNewUser(user *models.SignUpInput) (string, error) {
	insertResult, err := ur.store.InsertOne(ur.ctx, &user)

	//Catch Errs
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			//return "", utils.GenerateError(ErrDuplicateEmail, err)
			return "", utils.GenerateError(ErrDuplicateEmail, err)
		}
		return "", utils.GenerateError(ErrUserInsertion, err)
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
		return nil, utils.GenerateError(ErrInvalidIDHex, err)
	}

	user := &models.DBResponse{}
	filter := bson.M{"_id": objID}
	err = ur.store.FindOne(ur.ctx, filter).Decode(user)

	if err != nil {
		return nil, utils.GenerateError(ErrUserNotFound, err)
	}

	return user, nil
}

func (ur UserRepoImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	user := &models.DBResponse{}
	filter := bson.M{"email": strings.ToLower(email)}
	err := ur.store.FindOne(ur.ctx, filter).Decode(user)

	if err != nil {
		return nil, utils.GenerateError(ErrUserNotFound, err)
	}

	return user, nil
}

func (ur UserRepoImpl) FindAndUpdateUserByID(id string, data *models.UpdateInput) (*models.DBResponse, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return nil, utils.GenerateError(ErrInvalidUpdateInput, err)
	}

	// Convert String to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.GenerateError(ErrInvalidIDHex, err)
	}

	query := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$set", Value: doc}}
	result := ur.store.FindOneAndUpdate(ur.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedUser *models.DBResponse
	if err := result.Decode(&updatedUser); err != nil {
		return nil, utils.GenerateError(ErrUserNotFound, err)
	}

	return updatedUser, nil
}

func (ur UserRepoImpl) UpdateUserById(id string, update *models.UpdateInput) error {
	// Convert String to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return utils.GenerateError(ErrInvalidIDHex, err)
	}

	doc, err := utils.ToDoc(update)
	if err != nil {
		return utils.GenerateError(ErrInvalidUpdateInput, err)
	}

	filter := bson.M{"_id": objectID}
	updatePrimitive := bson.D{{Key: "$set", Value: doc}}
	res, err := ur.store.UpdateOne(ur.ctx, filter, updatePrimitive)

	if err != nil {
		return utils.GenerateError(ErrUserUpdate, err)
	}

	if res.MatchedCount < 1 {
		return utils.GenerateError(ErrUserNotFound, err)
	}

	return nil
}

func (ur UserRepoImpl) UpdateUserByEmail(email string, update *models.UpdateInput) error {
	doc, err := utils.ToDoc(update)
	if err != nil {
		return utils.GenerateError(ErrInvalidUpdateInput, err)
	}

	filter := bson.M{"email": email}
	updatePrimitive := bson.D{{Key: "$set", Value: doc}}
	res, err := ur.store.UpdateOne(ur.ctx, filter, updatePrimitive)

	if err != nil {
		return utils.GenerateError(ErrUserUpdate, err)
	}

	if res.MatchedCount < 1 {
		return utils.GenerateError(ErrUserNotFound, err)
	}

	return nil
}

func (ur UserRepoImpl) VerifyUserEmail(verificationCode string) error {
	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}

	res, err := ur.store.UpdateOne(ur.ctx, query, update)

	if err != nil {
		return utils.GenerateError(ErrUserVerification, err)
	}

	if res.MatchedCount < 1 {
		return utils.GenerateError(ErrUserNotFound, err)
	}

	return nil
}

func (ur UserRepoImpl) StorePasswordResetToken(userEmail string, passwordResetToken string) error {
	query := bson.D{{Key: "email", Value: strings.ToLower(userEmail)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	res, err := ur.store.UpdateOne(ur.ctx, query, update)

	if err != nil {
		return utils.GenerateError(ErrStorePasswordResetToken, err)
	}

	if res.MatchedCount < 1 {
		return utils.GenerateError(ErrUserNotFound, err)
	}

	return nil
}

func (ur UserRepoImpl) ResetUserPassword(passwordResetToken string, newPassword string) error {
	query := bson.D{{Key: "passwordResetToken", Value: passwordResetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: newPassword}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	res, err := ur.store.UpdateOne(ur.ctx, query, update)

	if err != nil {
		return utils.GenerateError(ErrResetPassword, err)
	}

	if res.MatchedCount < 1 {
		return utils.GenerateError(ErrUserNotFound, err)
	}

	return nil
}
