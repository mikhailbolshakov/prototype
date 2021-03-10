package storage

import (
	"fmt"
	kitLog "gitlab.medzdrav.ru/prototype/kit/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type Storage struct {
	Instance *gorm.DB
	DBName   string
	logger   kitLog.CLoggerFunc
}

type Params struct {
	UserName string
	Password string
	DBName   string
	Port     string
	Host     string
}

func Open(params *Params, logger kitLog.CLoggerFunc) (*Storage, error) {

	s := &Storage{
		DBName: params.DBName,
		logger: logger,
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable TimeZone=Europe/Moscow",
		params.UserName,
		params.Password,
		params.DBName,
		params.Port,
		params.Host,
	)

	// uncomment to log all queries
	cfg := &gorm.Config{
		Logger: gormLogger.New(
			logger(),
			//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold: time.Second * 10,  // Slow SQL threshold
				LogLevel:      gormLogger.Silent, // Log level
				Colorful:      false,              // Disable color
			},
		),
	}

	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}

	logger().Pr("db").Cmp(params.UserName).Inf("ok")

	s.Instance = db

	return s, nil

}

func (s *Storage) Close() {
	db, _ := s.Instance.DB()
	_ = db.Close()
}
