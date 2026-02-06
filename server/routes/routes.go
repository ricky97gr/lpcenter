package routes

import (
	"lpcenter/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		plugins := api.Group("/plugins")
		{
			plugins.POST("", handlers.CreatePlugin)
			plugins.GET("", handlers.GetPlugins)
			plugins.GET("/:id", handlers.GetPlugin)
			plugins.PUT("/:id", handlers.UpdatePlugin)
			plugins.DELETE("/:id", handlers.DeletePlugin)
		}

		licenses := api.Group("/licenses")
		{
			licenses.POST("", handlers.CreateLicenseRequest)
			licenses.GET("", handlers.GetLicenseRequests)
			licenses.GET("/:id", handlers.GetLicenseRequest)
			licenses.PUT("/:id/approve", handlers.ApproveLicense)
			licenses.PUT("/:id/reject", handlers.RejectLicense)
			licenses.GET("/user/:userId", handlers.GetUserLicenses)
		}
	}
}
