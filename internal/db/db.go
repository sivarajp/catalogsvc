package db

import (
	"fmt"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/sivarajp/catalogsvc/pkg/logger"
)

var (

	// Mongo stores the mongodb connection string information
	mongo *mgo.DialInfo

	db *mgo.Database

	Collection *mgo.Collection

)

// GetEnv accepts the ENV as key and a default string
// If the lookup returns false then it uses the default string else it leverages the value set in ENV variable
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	logger.Logger.Info("Setting default values for ENV variable " + key)
	return fallback
}

// ConnectDB accepts name of database and collection as a string
func ConnectDB(dbName string, collectionName string) *mgo.Session {

	dbUsername := GetEnv("CATALOG_DB_USERNAME", "mongoadmin")
	dbSecret := GetEnv("CATALOG_DB_PASSWORD", "Vmware2020!")

	// Get ENV variable or set to default value
	dbIP := GetEnv("CATALOG_DB_HOST", "catalog-mongo.acme.svc.cluste.local")
	dbPort := GetEnv("CATALOG_DB_PORT", "27017")

	mongoDBUrl := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", dbUsername, dbSecret, dbIP, dbPort)

	Session, error := mgo.DialWithTimeout(mongoDBUrl, time.Duration(15 * time.Second))

	if error != nil {
		fmt.Printf(error.Error())
		logger.Logger.Fatalf(error.Error())
		os.Exit(1)

	}

	db = Session.DB(dbName)

	error = db.Session.Ping()
	if error != nil {
		logger.Logger.Errorf("Unable to connect to database %s", dbName)
	}

	Collection = db.C(collectionName)

	logger.Logger.Info("Connected to database and the collection")

	return Session
}

// CloseDB accepst Session as input to close Connection to the database
func CloseDB(s *mgo.Session) {

	defer s.Close()
	logger.Logger.Info("Closed connection to db")
}
