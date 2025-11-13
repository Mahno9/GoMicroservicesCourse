package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCollectionAdapter wraps mongo.Collection to implement MongoCollection interface
type MongoCollectionAdapter struct {
	*mongo.Collection
}

func (a *MongoCollectionAdapter) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) MongoSingleResult {
	return &MongoSingleResultAdapter{a.Collection.FindOne(ctx, filter, opts...)}
}

func (a *MongoCollectionAdapter) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursor, error) {
	cursor, err := a.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &MongoCursorAdapter{cursor}, nil
}

func (a *MongoCollectionAdapter) Indexes() MongoIndexView {
	indexView := a.Collection.Indexes()
	return &MongoIndexViewAdapter{&indexView}
}

// MongoCursorAdapter wraps mongo.Cursor to implement MongoCursor interface
type MongoCursorAdapter struct {
	*mongo.Cursor
}

// MongoSingleResultAdapter wraps mongo.SingleResult to implement MongoSingleResult interface
type MongoSingleResultAdapter struct {
	*mongo.SingleResult
}

// MongoIndexViewAdapter wraps mongo.IndexView to implement MongoIndexView interface
type MongoIndexViewAdapter struct {
	*mongo.IndexView
}

// MongoDatabaseAdapter wraps mongo.Database to implement MongoDatabase interface
type MongoDatabaseAdapter struct {
	*mongo.Database
}

func (a *MongoDatabaseAdapter) Collection(name string, opts ...*options.CollectionOptions) MongoCollection {
	return &MongoCollectionAdapter{a.Database.Collection(name, opts...)}
}

var _ MongoCursor = (*MongoCursorAdapter)(nil)

type MongoCursor interface {
	All(ctx context.Context, results interface{}) error
	Close(ctx context.Context) error
}

var _ MongoCollection = (*MongoCollectionAdapter)(nil)

type MongoCollection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) MongoSingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursor, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	Indexes() MongoIndexView
}

var _ MongoIndexView = (*MongoIndexViewAdapter)(nil)

type MongoIndexView interface {
	CreateOne(ctx context.Context, indexModel mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
}

var _ MongoDatabase = (*MongoDatabaseAdapter)(nil)

type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) MongoCollection
}

var _ MongoSingleResult = (*MongoSingleResultAdapter)(nil)

type MongoSingleResult interface {
	Decode(v interface{}) error
}
