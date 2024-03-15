package routes

import (
	"github.com/daalvand/go-api/actions"
	"github.com/daalvand/go-api/db"
	"github.com/daalvand/go-api/middlewares"
	"github.com/daalvand/go-api/repositories"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *gin.Engine, db *db.Database) {
	registerUserRoutes(server, db)
	registerEventRoutes(server, db)
	server.GET("", func(context *gin.Context) {
		context.JSON(200, map[string]string{"message": "Hello World!"})
	})
}

func registerUserRoutes(server *gin.Engine, db *db.Database) {
	userRepo := repositories.NewUserRepository(db.DB)
	userActions := actions.NewUserActions(userRepo)
	userRoutes := server.Group("users")
	{
		userRoutes.POST("signup", userActions.Signup)
		userRoutes.POST("login", userActions.Login)
	}
}

func registerEventRoutes(server *gin.Engine, db *db.Database) {
	eventRepo := repositories.NewEventRepository(db.DB)
	eventActions := actions.NewEventActions(eventRepo)
	eventRoutes := server.Group("events")
	{
		eventRoutes.GET("", eventActions.GetEvents)
		eventRoutes.GET(":id", eventActions.GetEvent)
		authenticatedRoutes := eventRoutes.Group("")
		authenticatedRoutes.Use(middlewares.Authenticate)
		{
			authenticatedRoutes.POST("", eventActions.CreateEvent)
			authenticatedRoutes.PUT(":id", eventActions.UpdateEvent)
			authenticatedRoutes.DELETE(":id", eventActions.DeleteEvent)
			authenticatedRoutes.POST(":id/register", eventActions.RegisterForEvent)
			authenticatedRoutes.DELETE(":id/register", eventActions.CancelRegistration)
		}
	}
}
