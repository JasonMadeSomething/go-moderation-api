package utils

import (
	"context"
	"time"

	"github.com/go-moderation-api/config"
	"github.com/go-moderation-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents a MongoDB client connection
type MongoDB struct {
	client     *mongo.Client
	database   string
	collection string
}

// NewMongoDB creates a new MongoDB client
func NewMongoDB(cfg *config.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		client:     client,
		database:   cfg.MongoDatabase,
		collection: cfg.MongoCollection,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// GetCollection returns the MongoDB collection
func (m *MongoDB) GetCollection() *mongo.Collection {
	return m.client.Database(m.database).Collection(m.collection)
}

// FindModerationResult looks for a cached moderation result by content
func (m *MongoDB) FindModerationResult(ctx context.Context, content string) (*models.ModerationResult, error) {
	collection := m.GetCollection()
	
	filter := bson.M{"content": content}
	result := &models.ModerationResult{}
	
	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found, not an error
		}
		return nil, err
	}
	
	return result, nil
}

// SaveModerationResult saves a moderation result to the cache
func (m *MongoDB) SaveModerationResult(ctx context.Context, result *models.ModerationResult) error {
	collection := m.GetCollection()
	
	// Set timestamps
	now := time.Now()
	result.CreatedAt = now
	result.UpdatedAt = now
	
	// Check if a document with this content already exists
	filter := bson.M{"content": result.Content}
	update := bson.M{
		"$set": bson.M{
			"source_system": result.SourceSystem,
			"allowed":       result.Allowed,
			"openai_result": result.OpenAIResult,
			"updated_at":    now,
		},
	}
	
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}
