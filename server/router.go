package server

import (
	"github.com/504dev/logr/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:    []string{"Authorization", "Content-Type"},
		AllowAllOrigins: true,
	}))

	// oauth
	auth := controllers.AuthController{}
	auth.Init()
	{
		r.GET("/oauth/authorize", auth.Authorize)
		r.GET("/oauth/callback", auth.Callback)
	}

	// me
	me := controllers.MeController{}
	{
		r.GET("/me", auth.EnsureJWT, me.Me)
		r.GET("/me/dashboards", auth.EnsureJWT, me.Dashboards)
		r.POST("/me/dashboard", auth.EnsureJWT, me.AddDashboard)
		r.POST("/me/dashboard/share/:dash_id/to/:username", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDash, me.ShareDashboard)
		r.PUT("/me/dashboard/:dash_id", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDash, me.EditDashboard)
		r.DELETE("/me/dashboard/:dash_id", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDash, me.DeleteDashboard)
	}

	logsController := controllers.LogsController{}
	{
		r.GET("/logs", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, logsController.Find)
		r.GET("/logs/stats", auth.EnsureJWT, logsController.Stats)
	}

	countsController := controllers.CountsController{}
	{
		r.GET("/counts", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.Find)
		r.GET("/counts/snippet", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.FindSnippet)
		r.GET("/counts/stats", auth.EnsureJWT, countsController.Stats)

	}

	adminController := controllers.AdminController{}
	{
		r.GET("/dashboards", auth.EnsureJWT, auth.EnsureAdmin, adminController.Dashboards)
		r.GET("/dashboard/:id", auth.EnsureJWT, auth.EnsureAdmin, adminController.DashboardById)
		r.GET("/users", auth.EnsureJWT, auth.EnsureAdmin, adminController.Users)
		r.GET("/user/:id", auth.EnsureJWT, auth.EnsureAdmin, adminController.UserById)
	}

	wsController := controllers.WsController{}
	r.GET("/ws", wsController.Index)

	return r
}
