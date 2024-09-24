package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DBTrxUserKey             string = "userTrx"
	DBTrxProgramKey          string = "programTrx"
	DBTrxUserActivityLogKey  string = "userActLogTrx"
	DBTrxSubscriptionKey     string = "subscriptionTrx"
	DBTrxPromoKey            string = "promoTrx"
	DBTrxUserPointHistoryKey string = "userPointHistoryTrx"
)

type DataSources struct {
	DB *gorm.DB
}

type dbConnConfig struct {
	port     int
	user     string
	password string
	dbName   string
	sslmode  string
	host     string
}

func InitDS() (*DataSources, error) {
	var err error
	var db *DataSources
	host := os.Getenv("PG_HOST")
	port, err := strconv.Atoi(os.Getenv("PG_PORT"))
	if err != nil {
		log.Printf("failed to parse port from env: %v", err)
	}
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DB")
	sslmode := "disable"

	dbConnConf := dbConnConfig{
		port:     port,
		user:     user,
		password: password,
		dbName:   dbname,
		sslmode:  sslmode,
		host:     host,
	}

	log.Printf("Initializing databases\n")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=%s",
		dbConnConf.host, dbConnConf.user, dbConnConf.dbName, dbConnConf.password, dbConnConf.port, dbConnConf.sslmode)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
		return nil, fmt.Errorf("error connecting to db: %v", err)
	}
	log.Println("psql connected..")

	// Initialize redis connection
	// redisHost := os.Getenv("REDIS_HOST")
	// redisPort := os.Getenv("REDIS_PORT")
	// redisPassword := os.Getenv("REDIS_PASSWORD")

	// log.Printf("Connecting to Redis \n")
	// redisDB := redis.NewClient(&redis.Options{
	// 	Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
	// 	Password: redisPassword,
	// 	DB:       0,
	// })

	// // verift redis connection
	// _, err = redisDB.Ping(context.Background()).Result()
	// if err != nil {
	// 	log.Panic(err)
	// 	return nil, fmt.Errorf("error connecting to redis: %v", err)
	// }
	// log.Println("redis connected..")

	db = &DataSources{
		DB: gormDB,
	}

	return db, nil
}

func (d *DataSources) Close() error {
	// gorm doesnt provide closing connection

	// closing redisdb

	// if err := d.RedisClient.Close(); err != nil {
	// 	return fmt.Errorf("error closing Redis Client: %v", err)
	// }
	return nil
}
