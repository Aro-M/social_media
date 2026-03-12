package router

import (
	"social_media/internal/config"
	"social_media/internal/controller"
	"social_media/internal/router/middleware"

	_ "social_media/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Router {
	return &Router{ctrl: ctrl}
}

func (r *Router) Register(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", r.ctrl.Health)

	v1 := e.Group("/api/v1")
	{
		v1.GET("/ping", r.ctrl.Ping)

		users := v1.Group("/users")
		{
			users.POST("/register", r.ctrl.RegisterUser)
			users.POST("/login", r.ctrl.LoginUser)
			users.GET("/search", r.ctrl.SearchUsers, middleware.JWTMiddleware(config.JWTSecret()))
			users.PUT("/profile", r.ctrl.UpdateProfile, middleware.JWTMiddleware(config.JWTSecret()))
		}

		protected := v1.Group("", middleware.JWTMiddleware(config.JWTSecret()))
		{

			friends := protected.Group("/friends")
			{
				friends.GET("", r.ctrl.ListFriends)
				friends.POST("/request", r.ctrl.SendFriendRequest)
				friends.GET("/requests", r.ctrl.ListFriendRequests)
				friends.PUT("/requests/:id/accept", r.ctrl.AcceptFriendRequest)
				friends.PUT("/requests/:id/decline", r.ctrl.DeclineFriendRequest)
				friends.DELETE("/:id", r.ctrl.RemoveFriend)
			}
		}
	}
}
