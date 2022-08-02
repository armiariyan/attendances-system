package helper

import (
	"armiariyan/attendances-system/entity"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}

func CreateAttendanceResponse(data entity.Attendance) ResponseAttendance {
	stringDate := UnixMilliToString(data.Date, "date")
	stringTime := UnixMilliToString(data.Time, "time")

	attendanceResponse := ResponseAttendance{
		Id:     data.Id,
		UserId: data.UserId,
		Label:  data.Label,
		Date:   stringDate,
		Time:   stringTime,
	}

	return attendanceResponse
}

func CreateAttendanceResponses(tmpResponse []entity.Attendance) []ResponseAttendance {
	// Create activity response
	var response []ResponseAttendance

	for _, data := range tmpResponse {
		CreateAttendanceResponse(data)
		response = append(response, CreateAttendanceResponse(data))
	}

	return response
}

func CreateActivityResponse(data entity.Activity) ResponseActivity {
	stringDataCreated := UnixMilliToString(data.DateCreated, "date")
	stringTimeCreated := UnixMilliToString(data.TimeCreated, "time")

	activityResponse := ResponseActivity{
		Id:          data.Id,
		UserId:      data.UserId,
		Description: data.Description,
		DateCreated: stringDataCreated,
		TimeCreated: stringTimeCreated,
	}

	return activityResponse
}

func CreateActivityResponses(tmpResponse []entity.Activity) []ResponseActivity {
	// Create activity response
	var response []ResponseActivity

	for _, data := range tmpResponse {
		CreateActivityResponse(data)
		response = append(response, CreateActivityResponse(data))
	}

	return response
}

func UnixMilliToString(data int64, kind string) string {
	if kind == "date" {
		return time.UnixMilli(data).Format("2006-01-02")
	} else {
		return time.UnixMilli(data).Format("15:04:05")
	}
}

func StringToUnixMilli(str string) int64 {
	// Change to time
	str = str + " 23:59:59"
	timeDate, _ := time.Parse("2006-01-02 15:04:05", str)

	// Change to unix mili
	return timeDate.UnixMilli()
}

// Today start at 00.00, today end at 23.59
func GenerateTodayUnixMilli() (result []int64) {
	strTimeTodayStart := time.Now().Format("2006-01-02")
	tmpTimeTodayStart, _ := time.Parse("2006-01-02", strTimeTodayStart)
	result = append(result, tmpTimeTodayStart.UnixMilli())

	strTimeTodayEnd := time.Now().Format("2006-01-02 15:04:05")
	tmpTimeTodayEnd, _ := time.Parse("2006-01-02 15:04:05", strTimeTodayEnd)
	result = append(result, tmpTimeTodayEnd.UnixMilli())

	fmt.Println(strTimeTodayStart, strTimeTodayEnd)
	return result
}

func GenerateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateIdAttendance() (id string) {
	return "ATD-" + GenerateRandomString(5)
}

func GenerateIdActivity() (id string) {
	return "ACT-" + GenerateRandomString(5)
}

func IsUserEmpty(data entity.User) bool {
	if cmp.Equal(data, entity.User{}) {
		return true
	} else {
		return false
	}
}

func IsActivityEmpty(data entity.Activity) bool {
	if cmp.Equal(data, entity.Activity{}) {
		return true
	} else {
		return false
	}
}

func IsActivitiesEmpty(activities []entity.Activity) bool {
	if cmp.Equal(activities, []entity.Activity{}) {
		return true
	} else {
		return false
	}
}

func IsLogin(status interface{}) bool {
	if status != true {
		return false
	} else {
		return true
	}
}

func IsAuthorize(param_id interface{}, session_id int) bool {
	if param_id != session_id {
		return false
	} else {
		return true
	}
}

func IsCheckIn(attendances []entity.Attendance) bool {
	result := GenerateTodayUnixMilli()

	// result[0] is start, [1] is end
	for _, attendance := range attendances {
		if attendance.Date >= result[0] && attendance.Date <= result[1] && attendance.Label == "check in" {
			return true
		}
	}
	return false
}

func ComparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// func IsAuthorize(token *jwt.Token, user_id int) bool {
// 	claims := token.Claims.(jwt.MapClaims)
// 	fmt.Println("claims=", claims)

// 	id, err := strconv.ParseInt(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	if id == int64(user_id) {
// 		return true
// 	} else {
// 		return false
// 	}
// }
