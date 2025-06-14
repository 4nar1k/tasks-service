package database

import (
	"errors"
	"github.com/4nar1k/project-protos/proto/task"
	"github.com/4nar1k/project-protos/proto/user"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logrus.Error("DATABASE_URL is not set")
		return nil, errors.New("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Error("Failed to connect to database")
		return nil, err
	}

	if err := db.AutoMigrate(&user.User{}, &task.Task{}); err != nil {
		logrus.WithError(err).Error("Failed to auto-migrate User model")
		return nil, err
	}

	logrus.Info("Successfully connected to database")
	return db, nil
}
