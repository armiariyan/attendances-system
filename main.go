package main

import (
	"armiariyan/attendances-system/config"
	"armiariyan/attendances-system/controller"
	"armiariyan/attendances-system/repository"
	"armiariyan/attendances-system/service"

	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	userService    service.UserService       = service.NewUserService(userRepository)
	userController controller.UserController = controller.NewUserController(userService)
)

func main() {
	defer config.CloseDatabaseConnection(db)

	r := config.InitWithSession()

	// seeder.DBSeed(db)

	r.GET("/", userController.Index)

	userRoutes := r.Group("api/")
	{
		userRoutes.GET("/check/health", userController.Healthcheck)
		userRoutes.POST("/register", userController.Register)
		userRoutes.POST("/login", userController.Login)
		userRoutes.POST("/logout", userController.Logout)

		userRoutes.POST("/checkin/:id", userController.CheckIn)
		userRoutes.POST("/checkout/:id", userController.CheckOut)

		userRoutes.POST("/activity/:id", userController.CreateActivity)
		userRoutes.PUT("/activity/:id/:id_activity", userController.UpdateActivity)
		userRoutes.DELETE("/activity/:id/:id_activity", userController.DeleteActivity)

		userRoutes.GET("/activity/:id", userController.GetActivityHistoryByDate)
		userRoutes.GET("/attendances/:id", userController.GetAttendancesHistory)
	}

	r.Run()
}
