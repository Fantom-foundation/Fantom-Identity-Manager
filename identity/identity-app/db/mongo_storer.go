package db

import (
	"context"
	"github.com/op/go-logging"
	"github.com/volatiletech/authboss"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"identity-app/config"
	applogging "identity-app/logging"
	"identity-app/model"
	"strings"
)

const (
	mongoStorerIdentifier = "mongodb"
)

var (
	mongoStorer = &MongoStorer{}

	_ StorerBase                    = mongoStorer
	_ authboss.CreatingServerStorer = mongoStorer
	//_ authboss.ConfirmingServerStorer  = mongoStorer
	_ authboss.RecoveringServerStorer = mongoStorer
	//_ authboss.RememberingServerStorer = mongoStorer
)

const (
	userCollection = "users"
	keyPID         = "username"
	keyEmail       = "email"
)

type MongoStorer struct {
	client   *mongo.Client
	database string
	log      *logging.Logger
}

func init() {
	RegisterStorer(func() StorerBase {
		return NewMongoStorer()
	})
}
func NewMongoStorer() *MongoStorer {
	return &MongoStorer{}
}

func (storer *MongoStorer) Close() {
	// do we have a client?
	if storer.client != nil {
		// try to disconnect
		err := storer.client.Disconnect(context.Background())
		if err != nil {
			storer.log.Errorf("error on closing database connection; %s", err.Error())
		}
	}
}

func (storer *MongoStorer) CanHandle(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	return scheme == mongoStorerIdentifier
}

func (storer *MongoStorer) FromConfig(cfg *config.Config, log *applogging.Logger) (StorerBase, error) {
	// Empty unrestricted context
	ctx := context.Background()

	// Connect the database
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DSN))
	if err != nil {
		return nil, err
	}
	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	storer.client = client
	storer.database = "identity"
	storer.log = &log.Logger
	return storer, nil
}

// Load the user
func (storer *MongoStorer) Load(_ context.Context, key string) (authboss.User, error) {
	collection := storer.getUserCollection()

	result := collection.FindOne(context.Background(), bson.D{{keyPID, key}})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, authboss.ErrUserNotFound
		} else {
			storer.log.Errorf("User load failed: %s", result.Err().Error())
			return nil, result.Err()
		}
	}

	var user model.User

	if err := result.Decode(&user); err != nil {
		storer.log.Errorf("User decode failed: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func (storer *MongoStorer) Save(_ context.Context, user authboss.User) error {
	collection := storer.getUserCollection()
	opts := options.Update().SetUpsert(true)
	u := user.(*model.User)

	value, err := bson.Marshal(u)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.D{{keyPID, user.GetPID()}},
		bson.D{{"$set", value}},
		opts)
	if err != nil {
		storer.log.Errorf("User save failed: %s", err.Error())
		return err
	}
	storer.log.Debugf("user saved [%s]", user.GetPID())
	return nil
}

func (storer *MongoStorer) New(_ context.Context) authboss.User {
	return &model.User{}
}

func (storer *MongoStorer) Create(_ context.Context, user authboss.User) error {
	collection := storer.getUserCollection()
	u := user.(*model.User)

	value, err := bson.Marshal(u)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.Background(), value)
	if err != nil {
		storer.log.Errorf("User creation failed: %s", err.Error())
		return authboss.ErrUserFound
	}
	return nil
}

func (storer *MongoStorer) LoadByRecoverSelector(_ context.Context, selector string) (authboss.RecoverableUser, error) {
	collection := storer.getUserCollection()
	result := collection.FindOne(context.Background(), bson.D{
		{"$or", []interface{}{
			bson.D{{keyEmail, selector}},
			bson.D{{keyPID, selector}},
		}},
	})
	if result.Err() != nil {
		if result.Err() != mongo.ErrNoDocuments {
			storer.log.Errorf("User search failed: %s", result.Err().Error())
		}
		return nil, authboss.ErrUserNotFound
	}

	var user model.User

	if err := result.Decode(&user); err != nil {
		storer.log.Errorf("User decode failed: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func (storer *MongoStorer) getUserCollection() *mongo.Collection {
	return storer.client.Database(storer.database).Collection(userCollection)
}
