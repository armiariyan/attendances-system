package repository

import (
	"armiariyan/attendances-system/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(data entity.User) entity.User
	VerifyCredential(email string) interface{}
	GetDataByEmail(email string) entity.User
	ChangeStatusLogin(data entity.User) entity.User
	GetUserById(user_id int) entity.User
	GetActivityById(act_id string) entity.Activity
	CheckIn(data entity.Attendance) entity.Attendance
	CreateActivity(data entity.Activity) entity.Activity
	UpdateActivity(data entity.Activity) entity.Activity
	DeleteActivity(activity entity.Activity)
	GetActivityHistoryByDate(startDate, endDate int64) []entity.Activity
	GetAttendancesHistory(user_id int) []entity.Attendance
}

type userConnection struct {
	connection *gorm.DB
}

// Construct
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) GetAttendancesHistory(user_id int) []entity.Attendance {
	var attendances []entity.Attendance
	db.connection.Find(&attendances, "user_id = ?", user_id)
	return attendances
}

func (db *userConnection) RegisterUser(user entity.User) entity.User {
	// user.Password = hashAndSalt([]byte(user.Password))
	db.connection.Create(&user)
	// db.connection.Preload("User").Find(&b)
	return user
}

func (db *userConnection) GetDataByEmail(email string) entity.User {
	var user entity.User
	db.connection.Where("email = ?", email).Take(&user)
	return user
}

func (db *userConnection) VerifyCredential(email string) interface{} {
	var user entity.User

	res := db.connection.Where("email = ?", email).Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (db *userConnection) ChangeStatusLogin(data entity.User) entity.User {
	// db.connection.Where("id = ?", data.Id).Update("status_login", data.StatusLogin).Take(&data)
	db.connection.Updates(&data)
	db.connection.Find(&data)
	return data
}

func (db *userConnection) GetUserById(user_id int) entity.User {
	var user entity.User
	db.connection.First(&user, "id = ?", user_id)
	return user
}

func (db *userConnection) GetActivityById(act_id string) entity.Activity {
	var activity entity.Activity
	db.connection.First(&activity, "id = ?", act_id)
	return activity
}

func (db *userConnection) CheckIn(data entity.Attendance) entity.Attendance {
	db.connection.Create(&data)
	db.connection.Find(&data)
	return data
}

func (db *userConnection) CreateActivity(data entity.Activity) entity.Activity {
	db.connection.Create(&data)
	db.connection.Find(&data)
	return data
}

func (db *userConnection) UpdateActivity(activity entity.Activity) entity.Activity {
	db.connection.Where("id = ?", activity.Id).Updates(&activity).Take(&activity)
	// db.connection.Updates(&activity)
	db.connection.Find(&activity)
	return activity
}

func (db *userConnection) DeleteActivity(activity entity.Activity) {
	db.connection.Delete(&activity)
}

func (db *userConnection) GetActivityHistoryByDate(startDate, endDate int64) []entity.Activity {
	var activities []entity.Activity
	db.connection.Where("date_created >= ? AND date_created <= ?", startDate, endDate).Find(&activities)
	return activities
}
