package config

import (
	"os"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
)

func InitWithSession() (r *gin.Engine) {
	r = gin.Default()

	ss := os.Getenv("session_secret")
	store := gormsessions.NewStore(SetupDatabaseConnection(), true, []byte(ss))
	r.Use(sessions.Sessions("session_id", store)) // set session name

	return
}
