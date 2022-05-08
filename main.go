package main

import (
	"Cloudocs/db"
	"Cloudocs/middleware"

	"github.com/gin-gonic/gin"
)

func ginInit() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r = CollectRoute(r)
	panic(r.Run("10.0.4.16:9000"))
}

func main() {
	// 打开数据库
	db.MongoDB.Open()
	defer db.MongoDB.Close()
	// Gin初始化
	ginInit()
}
