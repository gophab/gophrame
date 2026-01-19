package mongo

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/mongo/config"
	"github.com/gophab/gophrame/core/sensitive"
	"github.com/gophab/gophrame/core/starter"
)

type Collection struct {
	*mongo.Collection
}

type SingleResult struct {
	*mongo.SingleResult
}
type Cursor struct {
	*mongo.Cursor
}

func (r *Collection) InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	sensitive.Translate(document, true)
	return r.Collection.InsertOne(ctx, document, opts...)
}

func (r *Collection) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *SingleResult {
	return &SingleResult{
		SingleResult: r.Collection.FindOne(ctx, filter, opts...),
	}
}

func (r *Collection) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (cur *Cursor, err error) {
	cursor, err := r.Collection.Find(ctx, filter, opts...)
	if err == nil {
		return &Cursor{
			Cursor: cursor,
		}, nil
	}
	return nil, err
}

func (r *SingleResult) Decode(v any) error {
	err := r.SingleResult.Decode(v)
	if err == nil {
		sensitive.Translate(v, false)
	}
	return err
}

func (r *Cursor) Decode(v any) error {
	err := r.Cursor.Decode(v)
	if err == nil {
		sensitive.Translate(v, false)
	}
	return err
}

func (r *Cursor) All(ctx context.Context, results any) error {
	err := r.Cursor.All(ctx, results)
	if err == nil {
		sensitive.Translate(results, false)
	}
	return err
}

type MongoDB struct {
	Client *mongo.Client
}

func (m *MongoDB) GetDatabase(dbName string) *mongo.Database {
	return m.Client.Database(dbName)
}

func (m *MongoDB) DB() *mongo.Database {
	return m.GetDatabase(config.Setting.Database)
}

func (m *MongoDB) Close() {
	if m.Client != nil {
		if err := m.Client.Disconnect(context.TODO()); err != nil {
			logger.Error("Error disconnecting from MongoDB: ", err)
		}
		logger.Info("Disconnected from MongoDB")
	}
}

func init() {
	starter.RegisterStarter(Start)
}

func InitDB() (*MongoDB, error) {
	Setting := config.Setting
	logger.Debug("Mongo Settings: ", json.String(Setting))

	clientOpts := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d/?connect=direct", config.Setting.Host, config.Setting.Port)).
		SetMaxPoolSize(config.Setting.MaxPoolSize).
		SetMinPoolSize(config.Setting.MinPoolSize).
		SetMaxConnIdleTime(config.Setting.MaxConnIdleTime).
		SetConnectTimeout(config.Setting.ConnectTimeout).
		SetServerSelectionTimeout(config.Setting.ServerSelectionTimeout).
		SetAuth(options.Credential{
			Username:   config.Setting.User,
			Password:   config.Setting.Password,
			AuthSource: config.Setting.AdminDatabase,
		})
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}
	return &MongoDB{
		Client: client,
	}, nil
}

var Mongo *MongoDB

var (
	once sync.Once
)

func Start() {
	if config.Setting.Enabled {
		once.Do(func() {
			if db, err := InitDB(); err == nil {
				Mongo = db
				logger.Info("Init Mongo Database: ", config.Setting.Host, config.Setting.Port)
				inject.InjectValue("mongo", Mongo)
			} else {
				logger.Error("Init Mongo Database Error: ", err.Error())
			}
		})

	}
}
