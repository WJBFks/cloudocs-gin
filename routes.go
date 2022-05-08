package main

import (
	"Cloudocs/controller"
	"Cloudocs/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	// test
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	r.POST("/db", controller.DbInsert)
	r.GET("/db", controller.DbFinds)
	r.GET("/db/:id", controller.DbFindId)
	// user
	r.GET("/users", controller.GetUsers)
	r.GET("/users/token", middleware.AuthMiddleware(), controller.GetUsersToken)
	r.GET("/users/:id", controller.GetUsersId)
	r.POST("/users", controller.PostUsers)
	// doc
	r.POST("/docs", middleware.AuthMiddleware(), controller.PostDocs)
	r.GET("/docs", middleware.AuthMiddleware(), controller.GetDocs)
	r.GET("/docs/:id", middleware.AuthMiddleware(), controller.GetDocsId)
	r.PUT("/docs/:id", middleware.AuthMiddleware(), controller.PutDocsId)
	r.DELETE("/docs/:id", middleware.AuthMiddleware(), controller.DelDoc)
	// doc/share
	r.PUT("/docs/share/:id", middleware.AuthMiddleware(), controller.AddShareDoc)
	r.GET("/docs/share/:id", middleware.AuthMiddleware(), controller.GetDocShare)
	r.DELETE("/docs/share/:id", middleware.AuthMiddleware(), controller.DelDocShareTel)
	r.GET("/docs/share", middleware.AuthMiddleware(), controller.GetShareDocs)
	// other
	r.POST("/login", controller.Login)
	r.PUT("/bind", controller.Bind)
	return r
}
