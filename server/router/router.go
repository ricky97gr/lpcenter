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
			public.GET("/products", controller.GetAllProducts)
			public.GET("/products/:id", controller.GetProduct)
			public.GET("/license-types", controller.GetAllLicenseTypes)
			public.GET("/license-types/:id", controller.GetLicenseType)
			public.POST("/plugins", controller.CreatePlugin)
			public.POST("/plugins/upload", controller.UploadPlugin)
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

		licenseTypes := api.Group("/license-types")
		{
			licenseTypes.GET("", controller.GetAllLicenseTypes)
			licenseTypes.GET("/:id", controller.GetLicenseType)
		}

		licenseTypesAuth := api.Group("/license-types")
		licenseTypesAuth.Use(middleware.Auth())
		{
			licenseTypesAuth.POST("", controller.CreateLicenseType)
			licenseTypesAuth.PUT("/:id", controller.UpdateLicenseType)
			licenseTypesAuth.DELETE("/:id", controller.DeleteLicenseType)
		}

		products := api.Group("/products")
		{
			products.GET("", controller.GetAllProducts)
			products.GET("/:id", controller.GetProduct)
		}

		productsAuth := api.Group("/products")
		productsAuth.Use(middleware.Auth())
		{
			productsAuth.POST("", controller.CreateProduct)
			productsAuth.PUT("/:id", controller.UpdateProduct)
			productsAuth.DELETE("/:id", controller.DeleteProduct)
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
			plugins.PUT("/:id/sign", controller.SignPlugin)
			plugins.PUT("/:id/publish", controller.PublishPlugin)
			plugins.PUT("/:id/disable", controller.DisablePlugin)
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

	go engine.Run(":8081")
}

func StartDownloadServer() {
	gin.SetMode(gin.DebugMode)

	engine := gin.New()
	engine.Use(middleware.Logger(), middleware.Recovery(), middleware.CORS())

	engine.POST("/download/task", controller.CreateDownloadTask)
	engine.GET("/download/status/:taskId", controller.GetDownloadStatus)
	engine.GET("/download/file", controller.DownloadFile)

	engine.Run(":8082")
}
