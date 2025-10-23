package listing

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user/app/entity"
	"user/app/repositories"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type event struct {
	Logger *logrus.Entry
	Db     *mongo.Client
}

func New(logger *logrus.Entry, db *mongo.Client) repositories.User {
	return &event{
		Logger: logger,
		Db:     db,
	}
}

func getNextSequence(ctx context.Context, db *mongo.Database, name string) (int64, error) {
	counters := db.Collection("counters")
	filter := bson.M{"_id": name}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		Seq int64 `bson:"seq"`
	}

	err := counters.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err = counters.InsertOne(ctx, bson.M{"_id": name, "seq": 1})
			if err != nil {
				return 0, err
			}
			return 1, nil
		}
		return 0, err
	}

	return result.Seq, nil
}

func (v *event) CreateUser(c echo.Context, params map[string]interface{}) (id int64, e error) {
	logger := c.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": params}).Info("repositories: CreateUser")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := v.Db.Database("mydb")
	collection := db.Collection("users")

	userID, err := getNextSequence(ctx, db, "user_id")
	if err != nil {
		logger.WithError(err).Error("failed to generate user_id")
		return 0, err
	}

	params["user_id"] = userID
	params["created_at"] = time.Now().UnixMicro()
	params["updated_at"] = time.Now().UnixMicro()

	_, err = collection.InsertOne(ctx, params)
	if err != nil {
		logger.WithError(err).Error("failed to insert user")
		return 0, err
	}

	logger.WithField("inserted_user_id", userID).Info("User created successfully")
	return userID, nil
}

func (v *event) GetAllUser(c echo.Context, meta entity.Meta) (data []entity.User, total int64, e error) {
	logger := c.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": meta}).Info("repositories: GetAllUser")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := v.Db.Database("mydb")
	collection := db.Collection("users")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(meta.Limit))
	findOptions.SetSkip(int64(meta.Offset))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		logger.WithError(err).Error("failed to query users")
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var u entity.User
		if err := cursor.Decode(&u); err != nil {
			logger.WithError(err).Error("failed to decode user")
			return nil, 0, err
		}
		data = append(data, u)
	}

	total, err = collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.WithError(err).Error("failed to count users")
		return nil, 0, err
	}

	return data, total, nil
}

func (v *event) GetUser(c echo.Context, id int64) (data entity.User, e error) {
	logger := c.Get("logger").(*logrus.Entry)
	logger.WithField("user_id", id).Info("repositories: GetUser")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := v.Db.Database("mydb")
	collection := db.Collection("users")

	filter := bson.M{"user_id": id}

	err := collection.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("user not found")
			return data, fmt.Errorf("user not found with id %d", id)
		}
		logger.WithError(err).Error("failed to get user")
		return data, err
	}

	return data, nil
}
