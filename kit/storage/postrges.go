package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Storage struct {
	Instance *gorm.DB
	DBName   string
}

type Params struct {
	UserName string
	Password string
	DBName   string
	Port     string
	Host     string
}

func Open(params *Params) (*Storage, error) {

	s := &Storage{
		DBName: params.DBName,
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable TimeZone=Europe/Moscow",
		params.UserName,
		params.Password,
		params.DBName,
		params.Port,
		params.Host,
	)

	cfg := &gorm.Config{
		//Logger: logger.New(
		//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		//	logger.Config{
		//		SlowThreshold: time.Second * 10, // Slow SQL threshold
		//		LogLevel:      logger.Info,      // Log level
		//		Colorful:      true,             // Disable color
		//	},
		//),
	}

	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to database %s", params.DBName)

	s.Instance = db

	return s, nil

}

