package controller

import (
	"armiariyan/attendances-system/dto"
	"armiariyan/attendances-system/entity"
	"armiariyan/attendances-system/helper"
	"armiariyan/attendances-system/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Index(context *gin.Context)
	Healthcheck(context *gin.Context)
	Register(context *gin.Context)
	Login(context *gin.Context)
	Logout(context *gin.Context)
	CheckIn(context *gin.Context)
	CheckOut(context *gin.Context)
	CreateActivity(context *gin.Context)
	UpdateActivity(context *gin.Context)
	DeleteActivity(context *gin.Context)
	GetActivityHistoryByDate(context *gin.Context)
	GetAttendancesHistory(context *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(user service.UserService) UserController {
	return &userController{
		userService: user,
	}
}

func (c *userController) Index(context *gin.Context) {
	context.Redirect(http.StatusFound, "/api/check/health")
}

func (c *userController) Healthcheck(context *gin.Context) {
	//Build response if success
	result := helper.Response{
		Status:  true,
		Message: "ok! check documentation at https://github.com/armiariyan/attendances_system",
		Errors:  "null",
		Data:    "null",
	}

	response := helper.BuildResponse(true, "healthcheck successfull", result)
	context.JSON(http.StatusOK, response)
}

func (c *userController) CheckIn(context *gin.Context) {

	// Take id from parameter and convert to int
	user_id, errConv := strconv.Atoi(context.Param("id"))
	if errConv != nil {
		response := helper.BuildErrorResponse("Failed to process request", errConv.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Check if user exist and logged in using session
	session := sessions.Default(context)
	if !helper.IsLogin(session.Get("loggedIn")) {
		response := helper.BuildErrorResponse("Failed to process request", "Please login first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Check if user authorized to access data
	if !helper.IsAuthorize(session.Get("user_id"), user_id) {
		response := helper.BuildErrorResponse("Failed to process request", "Unauthorized!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// // Take user id from session
	// i, ok := session.Get("user_id").(int)

	// Make Checkin Data
	checkInData := entity.Attendance{
		Id:     helper.GenerateIdAttendance(),
		UserId: user_id,
		Label:  "check in",
		Date:   time.Now().UnixMilli(),
		Time:   time.Now().UnixMilli(),
	}

	// Checkin
	result := helper.CreateAttendanceResponse(c.userService.CheckIn(checkInData))

	//Build response if success
	response := helper.BuildResponse(true, "Successfully Check In!", result)
	context.JSON(http.StatusOK, response)
}

func (c *userController) Login(context *gin.Context) {
	var loginDTO dto.LoginDTO

	// Fill loginDTO variable
	errDTO := context.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Verify the data exist
	tmpResult := c.userService.VerifyCredential(loginDTO.Email)
	if tmpResult == nil {
		// Build response error
		response := helper.BuildErrorResponse("Failed to process request", "Invalid email or password", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Take the entity result because tmpResult is interface{}
	entityResult := tmpResult.(entity.User)

	// Check if password match
	if !helper.ComparePassword(entityResult.Password, []byte(loginDTO.Password)) {
		response := helper.BuildErrorResponse("Failed to process request", "Invalid email or password", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Build sessions
	session := sessions.Default(context)
	// Set session variable
	session.Set("loggedIn", true)
	session.Set("user_id", entityResult.Id)
	session.Set("name", entityResult.Name)
	session.Set("email", entityResult.Email)
	session.Options(sessions.Options{MaxAge: 86400}) // Set session for one day (value in seconds)
	// Save session
	session.Save()

	// Build response if success
	response := helper.BuildResponse(true, "Successfully Logged In!", entityResult)
	context.JSON(http.StatusOK, response)
}

func (c *userController) Register(context *gin.Context) {
	// Fill RegisterDTO variable for validation
	var registerDTO dto.RegisterDTO

	errDTO := context.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Check Duplicate Email
	if c.userService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Email has been used", helper.EmptyObj{})
		context.JSON(http.StatusConflict, response)
		return
	}

	// Hash Password
	registerDTO.Password = helper.HashAndSalt([]byte(registerDTO.Password))

	// Create User
	createdUser := c.userService.CreateUser(registerDTO)

	//Build Response
	response := helper.BuildResponse(true, "User Registered! Please Login", createdUser)
	context.JSON(http.StatusCreated, response)
}

func (c *userController) CheckOut(context *gin.Context) {

	// Take id from parameter and convert to int
	user_id, errConv := strconv.Atoi(context.Param("id"))
	if errConv != nil {
		response := helper.BuildErrorResponse("Failed to process request", errConv.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	// Check if user exist and logged in using session
	session := sessions.Default(context)
	if !helper.IsLogin(session.Get("loggedIn")) {
		response := helper.BuildErrorResponse("Failed to process request", "Please login first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Check if user authorized to access data
	if !helper.IsAuthorize(session.Get("user_id"), user_id) {
		response := helper.BuildErrorResponse("Failed to process request", "Unauthorized!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Cek if user already check in today
	userAtd := c.userService.GetAttendancesHistory(user_id)
	if !helper.IsCheckIn(userAtd) {
		//Build response error because user not check in today
		response := helper.BuildErrorResponse("Failed to process request", "You should check in first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Make Checkout Data
	checkOutData := entity.Attendance{
		Id:     helper.GenerateIdAttendance(),
		UserId: user_id,
		Label:  "check out",
		Date:   time.Now().UnixMilli(),
		Time:   time.Now().UnixMilli(),
	}

	// Checkout
	result := helper.CreateAttendanceResponse(c.userService.CheckIn(checkOutData))

	//Build response if success
	response := helper.BuildResponse(true, "Successfully Check Out!", result)
	context.JSON(http.StatusOK, response)
}

func (c *userController) CreateActivity(context *gin.Context) {

	// Take id from parameter and convert to int
	user_id, errConv := strconv.Atoi(context.Param("id"))
	if errConv != nil {
		response := helper.BuildErrorResponse("Failed to process request", errConv.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	// Check if user exist and logged in using session
	session := sessions.Default(context)
	if !helper.IsLogin(session.Get("loggedIn")) {
		response := helper.BuildErrorResponse("Failed to process request", "Please login first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Check if user authorized to access data
	if !helper.IsAuthorize(session.Get("user_id"), user_id) {
		response := helper.BuildErrorResponse("Failed to process request", "Unauthorized!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Cek if user already check in today
	userAtd := c.userService.GetAttendancesHistory(user_id)
	if !helper.IsCheckIn(userAtd) {
		//Build response error because user not check in today
		response := helper.BuildErrorResponse("Failed to process request", "You should check in first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	var createActivityData entity.Activity
	// Fill the createActivityData
	errDTO := context.ShouldBind(&createActivityData)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Create activity data
	createActivityData = entity.Activity{
		Id:          helper.GenerateIdActivity(),
		UserId:      user_id,
		Description: createActivityData.Description,
		DateCreated: time.Now().UnixMilli(),
		TimeCreated: time.Now().UnixMilli(),
	}

	// Create activity
	result := helper.CreateActivityResponse(c.userService.CreateActivity(createActivityData))

	//Build response if success
	response := helper.BuildResponse(true, "Successfully Created Activity!", result)
	context.JSON(http.StatusCreated, response)
}

func (c *userController) UpdateActivity(context *gin.Context) {

	// Take id from parameter and convert to int
	user_id, errConv := strconv.Atoi(context.Param("id"))
	if errConv != nil {
		response := helper.BuildErrorResponse("Failed to process request", errConv.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	// Check if user exist and logged in using session
	session := sessions.Default(context)
	if !helper.IsLogin(session.Get("loggedIn")) {
		response := helper.BuildErrorResponse("Failed to process request", "Please login first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Check if user authorized to access data
	if !helper.IsAuthorize(session.Get("user_id"), user_id) {
		response := helper.BuildErrorResponse("Failed to process request", "Unauthorized!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Cek if user already check in today
	userAtd := c.userService.GetAttendancesHistory(user_id)
	if !helper.IsCheckIn(userAtd) {
		//Build response error because user not check in today
		response := helper.BuildErrorResponse("Failed to process request", "You should check in first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Get activity data
	act_id := context.Param("id_activity")
	actData := c.userService.GetActivityById(act_id)

	// Check if activity data empty
	if helper.IsActivityEmpty(actData) {
		//Build response error because activity data empty
		response := helper.BuildErrorResponse("Failed to process request", "Activity not found", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	var updateActivityData entity.Activity
	// Fill the updateActivityData
	errDTO := context.ShouldBind(&updateActivityData)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Create new activity data
	updateActivityData = entity.Activity{
		Id:          actData.Id,
		UserId:      actData.UserId,
		Description: updateActivityData.Description,
		DateCreated: actData.DateCreated,
		TimeCreated: actData.TimeCreated,
	}

	// Create activity response
	result := helper.CreateActivityResponse(c.userService.UpdateActivity(updateActivityData))

	//Build response if success
	response := helper.BuildResponse(true, "Successfully Update Activity!", result)
	context.JSON(http.StatusCreated, response)
}

func (c *userController) DeleteActivity(context *gin.Context) {
	// Take id from parameter and convert to int
	user_id, errConv := strconv.Atoi(context.Param("id"))
	if errConv != nil {
		response := helper.BuildErrorResponse("Failed to process request", errConv.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	// Check if user exist and logged in using session
	session := sessions.Default(context)
	if !helper.IsLogin(session.Get("loggedIn")) {
		response := helper.BuildErrorResponse("Failed to process request", "Please login first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Check if user authorized to access data
	if !helper.IsAuthorize(session.Get("user_id"), user_id) {
		response := helper.BuildErrorResponse("Failed to process request", "Unauthorized!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Cek if user already check in today
	userAtd := c.userService.GetAttendancesHistory(user_id)
	if !helper.IsCheckIn(userAtd) {
		//Build response error because user not check in today
		response := helper.BuildErrorResponse("Failed to process request", "You should check in first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Get activity data
	act_id := context.Param("id_activity")
	actData := c.userService.GetActivityById(act_id)

	// Check if activity data empty
	if helper.IsActivityEmpty(actData) {
		//Build response error because activity data empty
		response := helper.BuildErrorResponse("Failed to process request", "Activity not found", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	// Delete
	c.userService.DeleteActivity(actData)

	// Build response if success
	res := helper.BuildResponse(true, "Activity deleted!", helper.EmptyObj{})
	context.JSON(http.StatusOK, res)
}

func (c *userController) GetAttendancesHistory(context *gin.Context) {

	// Take id from parameter and convert to int
	user_id, errConv := strconv.Atoi(context.Param("id"))
	if errConv != nil {
		response := helper.BuildErrorResponse("Failed to process request", errConv.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	// Check if user exist and logged in using session
	session := sessions.Default(context)
	if !helper.IsLogin(session.Get("loggedIn")) {
		response := helper.BuildErrorResponse("Failed to process request", "Please login first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Check if user authorized to access data
	if !helper.IsAuthorize(session.Get("user_id"), user_id) {
		response := helper.BuildErrorResponse("Failed to process request", "Unauthorized!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Get attendances history
	response := helper.CreateAttendanceResponses(c.userService.GetAttendancesHistory(user_id))
	// If attendances history empty
	if response == nil {
		res := helper.BuildResponse(true, "Successfully get attendance history!", "attendances history is empty")
		context.JSON(http.StatusOK, res)
		return
	}

	// Build response if success
	res := helper.BuildResponse(true, "Successfully get attendance history!", response)
	context.JSON(http.StatusOK, res)
}

func (c *userController) GetActivityHistoryByDate(context *gin.Context) {
	// Take id from parameter and convert to int
	user_id, errConv := strconv.Atoi(context.Param("id"))
	if errConv != nil {
		response := helper.BuildErrorResponse("Failed to process request", errConv.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	// Check if user exist and logged in using session
	session := sessions.Default(context)
	if !helper.IsLogin(session.Get("loggedIn")) {
		response := helper.BuildErrorResponse("Failed to process request", "Please login first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Check if user authorized to access data
	if !helper.IsAuthorize(session.Get("user_id"), user_id) {
		response := helper.BuildErrorResponse("Failed to process request", "Unauthorized!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	// Take start date and end date from querry
	startDate := helper.StringToUnixMilli(context.Query("startDate"))
	endDate := helper.StringToUnixMilli(context.Query("endDate"))

	// Get activity history
	activities := c.userService.GetActivityHistoryByDate(startDate, endDate)

	// Check if activity in range date input empty
	// But this will return status No Content and no response were made
	if helper.IsActivitiesEmpty(activities) {
		fmt.Println("request succes but status no content, no response were made")
		res := helper.BuildResponse(true, "Activities in that range date is empty", activities)
		context.JSON(http.StatusNoContent, res)
		return
	}

	// Create activity response
	response := helper.CreateActivityResponses(activities)

	// Build response if success
	res := helper.BuildResponse(true, "Successfully get activity history!", response)
	context.JSON(http.StatusOK, res)
}

func (c *userController) Logout(context *gin.Context) {

	// Check if user exist and logged in using session
	session := sessions.Default(context)
	if !helper.IsLogin(session.Get("loggedIn")) {
		response := helper.BuildErrorResponse("Failed to process request", "Please login first!", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	session.Set("user_id", "") // this will mark the session as "written" and hopefully remove the username
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1}) // this sets the cookie with a MaxAge of 0
	session.Save()

	//Build response if success
	response := helper.BuildResponse(true, "Successfully Logged Out!", nil)
	context.JSON(http.StatusOK, response)
}
