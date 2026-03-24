package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	repoerrors "github.com/thealish/go-microservice/internal/domain/errors"
	"github.com/thealish/go-microservice/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// MongoUserRepository implements UserRepository using the MongoDB driver.
type MongoUserRepository struct {
	col      *mongo.Collection
	counters *mongo.Collection
}

// NewUserMongo returns a new MongoUserRepository.
func NewUserMongo(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		col:      db.Collection("users"),
		counters: db.Collection("counters"),
	}
}

// nextID atomically increments and returns the next sequence value for the given collection.
func (r *MongoUserRepository) nextID(ctx context.Context) (uint, error) {
	filter := bson.M{"_id": "users"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result struct {
		Seq uint `bson:"seq"`
	}
	if err := r.counters.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result); err != nil {
		return 0, fmt.Errorf("next id: %w", err)
	}
	return result.Seq, nil
}

// notDeletedFilter returns a filter that excludes soft-deleted documents.
func notDeletedFilter() bson.M {
	return bson.M{"deleted_at": bson.M{"$eq": nil}}
}

func (r *MongoUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	filter := bson.M{"id": id}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}

	var user models.User
	if err := r.col.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("get user by id: %w", repoerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.M{"email": strings.ToLower(email)}
	for k, v := range notDeletedFilter() {
		filter[k] = v
	}

	var user models.User
	if err := r.col.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("get user by email: %w", repoerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *MongoUserRepository) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.col.Find(ctx, notDeletedFilter(), opts)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	return users, nil
}

func (r *MongoUserRepository) Create(ctx context.Context, user *models.User) error {
	id, err := r.nextID(ctx)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	user.ID = id
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if _, err := r.col.InsertOne(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("create user: %w", repoerrors.ErrCannotCreate)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *models.User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.UpdatedAt = time.Now()

	filter := bson.M{"id": user.ID}
	update := bson.M{"$set": user}

	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update user: %w", repoerrors.ErrCannotUpdate)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update user: %w", repoerrors.ErrNotFound)
	}
	return nil
}

func (r *MongoUserRepository) Delete(ctx context.Context, id uint) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true}}}

	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("delete user: %w", repoerrors.ErrCannotDelete)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("delete user: %w", repoerrors.ErrNotFound)
	}
	return nil
}

func (r *MongoUserRepository) Restore(ctx context.Context, id uint) error {
	filter := bson.M{"id": id}
	update := bson.M{"$unset": bson.M{"deleted_at": ""}}

	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("restore user: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("restore user: %w", repoerrors.ErrNotFound)
	}
	return nil
}

func (r *MongoUserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.col.CountDocuments(ctx, notDeletedFilter())
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}
