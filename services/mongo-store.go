package services

import (
	"api/consts"
	store_models "api/models/store"
	"api/util"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

// MongoStore contains mongo.Client
type MongoStore struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

// NewMongoStore creates a new Client and then initializes it using the Connect method.
func NewMongoStore(ctx context.Context,
	cfg *store_models.StoreConfig) (ms *MongoStore, err error) {

	uri := fmt.Sprintf("mongodb://%s:%s", cfg.Addr, cfg.Port)

	// set auth if needed
	clientOpts := options.Client().ApplyURI(uri)
	if cfg.Login != "" && cfg.Password != "" {
		// set credentials
		credential := options.Credential{
			Username: cfg.Login,
			Password: cfg.Password,
		}
		// set db client options
		clientOpts = options.Client().ApplyURI(
			uri).SetAuth(credential).SetRegistry(mongoRegistry)
	}

	// init client
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return
	}

	// ping db
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}
	// set db
	db := client.Database(cfg.DB)

	// set table
	collection := db.Collection(cfg.Table)

	return &MongoStore{
		Client:     client,
		Database:   db,
		Collection: collection,
	}, nil
}

// Performs a specific action on the database according to the received DoID
func (ms *MongoStore) DoOne(act store_models.DoID, req store_models.IStoreDoRequest) (err error) {

	if req == nil {
		return fmt.Errorf("the request couldn't be empty")
	}
	switch act {
	case store_models.ADD:
		err = ms.InsertOne(req)
	case store_models.MODIFY:
		err = ms.UpdateOne(req)
	case store_models.DELETE:
		err = ms.DeleteOne(req)
	default:
		err = fmt.Errorf("wrong DoID type")
	}
	return
}

// Performs a specific getting on the database according to the received GetID
func (ms *MongoStore) Get(act store_models.GetID,
	filter interface{}) (results []store_models.IStoreGetResponse,
	err error) {

	var errs []error
	switch act {
	case store_models.GET_ALL:
		results, errs = ms.GetAll()
	case store_models.GET_FILTERED:
		if filter == nil {
			return nil, fmt.Errorf("the filter couldn't be empty")
		}
		results, errs = ms.GetFiltered(filter)
	default:
		return nil, fmt.Errorf("wrong DoID type")
	}
	// join all errors to the one
	err = util.JoinErrors(errs)
	return
}

// Insert one document to the DB
func (ms *MongoStore) InsertOne(req store_models.IStoreDoRequest) (err error) {

	_, err = ms.Collection.InsertOne(ctx, req)
	return
}

func (ms *MongoStore) UpdateOne(req store_models.IStoreDoRequest) (err error) {
	filter := bson.D{{"_id", req.GetID()}}
	// _, err = ms.collection.UpdateOne(ctx, filter, req)
	pByte, err := bson.Marshal(req)
	if err != nil {
		return err
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		return
	}

	res, err := ms.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})

	if res.MatchedCount == 0 {
		return fmt.Errorf(consts.STORE_KEY_NOT_FOUND, req.GetID())
	}

	return err
}

func (ms *MongoStore) DeleteOne(req store_models.IStoreDoRequest) (err error) {
	filter := bson.D{{"_id", req.GetID()}}

	res, err := ms.Collection.DeleteOne(ctx, filter)

	if res.DeletedCount == 0 {
		return fmt.Errorf(consts.STORE_KEY_NOT_FOUND, req.GetID())
	}
	return err
}

func (ms *MongoStore) GetAll() (results []store_models.IStoreGetResponse,
	errs []error) {

	// send empty filter
	return ms.GetFiltered(bson.D{{}})
}

func (ms *MongoStore) GetFiltered(filter interface{}) (results []store_models.IStoreGetResponse,
	errs []error) {

	// d.Shared.BsonToJSONPrint(filter)

	cur, err := ms.Collection.Find(ctx, filter)
	if err != nil {
		errs = append(errs, err)
		return
	}
	defer cur.Close(ctx)

	// Loop through the cursor
	for cur.Next(ctx) {
		res := make(store_models.IStoreGetResponse, 0)
		err := cur.Decode(res)
		if err != nil {
			errs = append(errs, err)
		} else {
			results = append(results, res)
		}
	}
	if err := cur.Err(); err != nil {
		errs = append(errs, err)
	}

	return
}
