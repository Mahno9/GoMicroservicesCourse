package part

import (
	"context"
	def "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ def.MongoCollection = (*mongoCollection)(nil)

type mongoCollection struct {
	// collection def.MongoCollection
}

// Find implements repository.MongoCollection.
func (m *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	panic("unimplemented")
}

// FindOne implements repository.MongoCollection.
func (m *mongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	panic("unimplemented")
}

// Indexes implements repository.MongoCollection.
func (m *mongoCollection) Indexes() mongo.IndexView {
	panic("unimplemented")
}

// InsertMany implements repository.MongoCollection.
func (m *mongoCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	panic("unimplemented")
}
