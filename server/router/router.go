package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ricky97gr/lpcenter/server/controller"
	"github.com/ricky97gr/lpcenter/server/middleware"
)

func Start() {
	gin.SetMode(gin.DebugMode)

	engine := gin.New()
	engine.Use(middleware.Logger(), middleware.Recovery(), middleware.CORS())

	api := engine.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", controller.Login)
			auth.POST("/logout", controller.Logout)
		}

		public := api.Group("/public")
		{
			public.GET("/plugins", controller.GetPublicPlugins)
		}

		users := api.Group("/users")
		users.Use(middleware.Auth())
		{
			users.GET("", controller.GetAllUsers)
			users.GET("/me", controller.GetCurrentUser)
			users.POST("", controller.CreateUser)
			users.GET("/:id", controller.GetUser)
			users.PUT("/:id", controller.UpdateUser)
			users.DELETE("/:id", controller.DeleteUser)
			users.POST("/:id/reset-password", controller.ResetPassword)
			users.POST("/change-password", controller.ChangePassword)
		}

		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.Auth())
		{
			dashboard.GET("/stats", controller.GetDashboardStats)
			dashboard.GET("/recent-licenses", controller.GetRecentLicenses)
			dashboard.GET("/recent-plugins", controller.GetRecentPlugins)
		}

		versions := api.Group("/versions")
		versions.Use(middleware.Auth())
		{
			versions.GET("", controller.GetAllVersions)
			versions.POST("", controller.CreateVersion)
			versions.GET("/:id", controller.GetVersion)
			versions.PUT("/:id", controller.UpdateVersion)
			versions.DELETE("/:id", controller.DeleteVersion)
		}

		products := api.Group("/products")
		products.Use(middleware.Auth())
		{
			products.GET("", controller.GetAllProducts)
			products.POST("", controller.CreateProduct)
			products.GET("/:id", controller.GetProduct)
			products.PUT("/:id", controller.UpdateProduct)
			products.DELETE("/:id", controller.DeleteProduct)
		}

		plugins := api.Group("/plugins")
		plugins.Use(middleware.Auth())
		{
			plugins.GET("", controller.GetAllPlugins)
			plugins.POST("", controller.CreatePlugin)
			plugins.GET("/:id", controller.GetPlugin)
			plugins.PUT("/:id", controller.UpdatePlugin)
			plugins.DELETE("/:id", controller.DeletePlugin)
			plugins.PUT("/:id/status", controller.UpdatePluginStatus)
		}

		downloads := api.Group("/downloads")
		downloads.Use(middleware.Auth())
		{
			downloads.POST("/record", controller.RecordDownload)
			downloads.GET("/plugin/:uuid", controller.GetPluginDownloadStats)
			downloads.GET("/stats", controller.GetAllPluginDownloadStats)
		}

		licenses := api.Group("/licenses")
		licenses.Use(middleware.Auth())
		{
			licenses.GET("", controller.GetAllLicenses)
			licenses.POST("", controller.CreateLicense)
			licenses.GET("/:id", controller.GetLicense)
			licenses.PUT("/:id", controller.UpdateLicense)
			licenses.DELETE("/:id", controller.DeleteLicense)
			licenses.PUT("/:id/approve", controller.ApproveLicense)
			licenses.PUT("/:id/reject", controller.RejectLicense)
			licenses.GET("/:id/download", controller.DownloadLicenseFile)
		}
	}

	engine.Run(":8081")
}
