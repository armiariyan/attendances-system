package service

import (
	"armiariyan/attendances-system/dto"
	"armiariyan/attendances-system/entity"
	"armiariyan/attendances-system/helper"
	"armiariyan/attendances-system/repository"
	"log"

	"github.com/mashingan/smapping"
)

type UserService interface {
	CreateUser(user dto.RegisterDTO) entity.User
	VerifyCredential(email string) interface{}
	ChangeStatusLogin(data entity.User) entity.User
	GetUserById(user_id int) entity.User
	GetActivityById(act_id string) entity.Activity
	CheckIn(data entity.Attendance) entity.Attendance
	CreateActivity(data entity.Activity) entity.Activity
	UpdateActivity(data entity.Activity) entity.Activity
	DeleteActivity(data entity.Activity)
	GetActivityHistoryByDate(startDate, endDate int64) []entity.Activity
	GetAttendancesHistory(user_id int) []entity.Attendance
	IsDuplicateEmail(email string) bool
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		userRepository: repository,
	}
}

func (service *userService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.userRepository.RegisterUser(userToCreate)
	return res
}

func (service *userService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.GetDataByEmail(email)
	// Return false if data empty, true if there is one
	return !helper.IsUserEmpty(res)
}

func (service *userService) VerifyCredential(email string) interface{} {
	res := service.userRepository.VerifyCredential(email)
	if res == nil {
		return nil
	} else {
		return res
	}
}

func (service *userService) ChangeStatusLogin(data entity.User) entity.User {
	res := service.userRepository.ChangeStatusLogin(data)
	return res
}

func (service *userService) GetUserById(user_id int) entity.User {
	return service.userRepository.GetUserById(user_id)
}

func (service *userService) GetActivityById(act_id string) entity.Activity {
	return service.userRepository.GetActivityById(act_id)
}

func (service *userService) CheckIn(data entity.Attendance) entity.Attendance {
	return service.userRepository.CheckIn(data)
}

func (service *userService) CreateActivity(data entity.Activity) entity.Activity {
	return service.userRepository.CreateActivity(data)
}

func (service *userService) UpdateActivity(data entity.Activity) entity.Activity {
	return service.userRepository.UpdateActivity(data)
}

func (service *userService) DeleteActivity(data entity.Activity) {
	service.userRepository.DeleteActivity(data)
}

func (service *userService) GetActivityHistoryByDate(startDate, endDate int64) []entity.Activity {
	return service.userRepository.GetActivityHistoryByDate(startDate, endDate)
}

func (service *userService) GetAttendancesHistory(user_id int) []entity.Attendance {
	return service.userRepository.GetAttendancesHistory(user_id)
}
