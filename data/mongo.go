package data

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"WebService/config"
	"WebService/log"
	"WebService/model"
)

const (
	USERS    = "Users"
	EXAMPLES = "Examples"
	COUNTERS = "Counters"
)
const mversion = "0.0.1"

type MongoDB struct {
	Config  *config.MongoDBConfig
	Session *mgo.Session
}

func CreateMongoDB(c *config.MongoDBConfig) *MongoDB {
	log.LogTrace("Connecting to MongoDB: %s\n", c.MongoUri)

	session, err := mgo.Dial(c.MongoUri)
	if err != nil {
		log.LogError("Error connecting to MongoDB: %s", err)
		return nil
	}

	return &MongoDB{
		Config:  c,
		Session: session,
	}
}
func (mongo *MongoDB) initDB() {
	session := mongo.getSession()
	defer session.Close()
	cdv := session.DB(mongo.Config.MongoDb).C("data_version")
	dv := bson.M{"_id": "data_version", "version": mversion}
	dvd := bson.M{}
	cdv.FindId("data_version").One(&dvd)
	if dv["version"] != dvd["version"] {
		session.DB(mongo.Config.MongoDb).DropDatabase()
		cdv.Insert(dv)
	}
	index := func() {

	}
	index()
}
func (mongo *MongoDB) getNextSequence(name string) (int, error) {
	session := mongo.getSession()
	defer session.Close()

	counters := session.DB(mongo.Config.MongoDb).C(COUNTERS)
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	doc := bson.M{}
	_, err := counters.FindId(name).Apply(change, &doc)
	return doc["seq"].(int), err
}

func (mongo *MongoDB) getSession() *mgo.Session {
	if mongo.Session == nil {
		var err error
		mgoSession, err := mgo.Dial(mongo.Config.MongoUri)
		if err != nil {
			log.LogError("Session Error connecting to MongoDB: %s", err)
			return nil
		}
		mongo.Session = mgoSession
	}
	return mongo.Session.Clone()
	//return mongo.Session.Copy()
}
func (mongo *MongoDB) Close() {
	mongo.Session.Close()
}

func (mongo *MongoDB) FindExample(id string) (*model.Example, error) {
	session := mongo.getSession()
	defer session.Close()
	c := session.DB(mongo.Config.MongoDb).C(EXAMPLES)
	if !bson.IsObjectIdHex(id) {
		return nil, mgo.ErrNotFound
	}
	var example *model.Example
	err := c.FindId(bson.ObjectIdHex(id)).One(&example)
	return example, err
}
func (mongo *MongoDB) InsertExample(example *model.Example) (*model.Example, error) {
	session := mongo.getSession()
	defer session.Close()
	c := session.DB(mongo.Config.MongoDb).C(EXAMPLES)

	err := c.Insert(example)
	return example, err
}

//todo mew function
func (mongo *MongoDB) CreateUser(user *model.User) (*model.User, error) {
	session := mongo.getSession()
	defer session.Close()
	c := session.DB(mongo.Config.MongoDb).C(USERS)

	user.Id = bson.NewObjectId()
	err := c.Insert(user)
	return user, err
}
