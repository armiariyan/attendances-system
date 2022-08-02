package seeder

import (
	"armiariyan/attendances-system/entity"

	"gorm.io/gorm"
)

type Seeder struct {
	Seeder interface{}
}

func User(db *gorm.DB) []entity.User {
	return []entity.User{
		{
			Id:       1,
			Name:     "User 1",
			Email:    "user1@gmail.com",
			Password: "password",
		},
		{
			Id:       2,
			Name:     "User 2",
			Email:    "user2@gmail.com",
			Password: "password",
		},
	}
}

func RegisterSeeders(db *gorm.DB) []Seeder {
	return []Seeder{
		{Seeder: User(db)},
	}
}

func DBSeed(db *gorm.DB) error {
	for _, seeder := range RegisterSeeders(db) {
		err := db.Debug().Create(seeder.Seeder).Error
		if err != nil {
			return err
		}
	}
	return nil
}
