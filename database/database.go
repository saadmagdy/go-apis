package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
type MongoDB struct {
	Ctx    context.Context
	DBName string
	DBUri  string
}
func NewMongoDb(ctx context.Context, dbname, uri string) *MongoDB {
	return &MongoDB{
		Ctx:    ctx,
		DBName: dbname,
		DBUri:  uri,
	}
}
func (db *MongoDB) CreateClient() (*mongo.Client, error) {
	client, err := mongo.Connect(db.Ctx, options.Client().ApplyURI(db.DBUri))
	if err != nil {
		return nil, err
	}
	// Check the connection
	if err = client.Ping(db.Ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}
func (db *MongoDB) CloseConnection(c *mongo.Client) error {
	err := c.Disconnect(db.Ctx)
	if err != nil {
		return err
	}
	return nil
}
func (db *MongoDB) CreateDataBase(c *mongo.Client) *mongo.Database {
	database := c.Database(db.DBName)
	return database
}
func (db *MongoDB) CreateCollection(d *mongo.Database, collName string) *mongo.Collection {
	coll := d.Collection(collName)
	return coll
}
