package data

import (
	"WebService/config"
	"WebService/log"
)

type DataStore struct {
	MongoConfig *config.MongoDBConfig
	HttpConfig  *config.HttpConfig

	Mongo *MongoDB
}

func NewDataStore() *DataStore {
	ds := &DataStore{}
	ds.MongoConfig = config.GetMongoDBConfig()
	ds.HttpConfig = config.GetHttpConfig()
	return ds
}
func (ds *DataStore) StorageConnect() {
	if ds.MongoConfig != nil {
		log.LogInfo("Trying MongoDB storage")
		s := CreateMongoDB(config.GetMongoDBConfig())
		if s == nil {
			log.LogInfo("MongoDB storage unavailable")
		} else {
			log.LogInfo("Using MongoDB storage")
			ds.Mongo = s
		}
		s.initDB()
	}
}

func (ds *DataStore) StorageDisconnect() {
	if ds.Mongo != nil {
		ds.Mongo.Close()
	}
}
