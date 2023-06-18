package db

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)
var ctx = context.Background()

func ConnectDb() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Silent)

	log.Println("running migrations")
	db.AutoMigrate(User{}, Privilege{}, Workspace{}, UserPrivilege{}, Desk{}, Column{}, Card{}, Label{}, Attachment{}, CardsLabel{})

	DB = db
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	usrModel := DB.Model(User{})
	hashedPassword := sha256.Sum256([]byte("123456" + "loneliness"))
	user := User{
		ID:           10000,
		Name:         "test",
		Email:        "test@test.test",
		Password:     base64.URLEncoding.EncodeToString(hashedPassword[:]),
		IsAuthorized: true,
	}

	usrModel.Create(&user)
	result := RDB.RPush(ctx, fmt.Sprint(10000), "loneliness")
	if _, err := result.Result(); err != nil {
		log.Fatal(err)
	}

	privileges := []Privilege{
		{ID: 1, Name: "admin"},
		{ID: 2, Name: "user"},
		{ID: 3, Name: "guest"},
	}

	for _, privilege := range privileges {
		var priv = privilege
		DB.Model(Privilege{}).Create(&priv)
	}
}
