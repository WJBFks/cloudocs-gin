package controller

import (
	"Cloudocs/common"
	"Cloudocs/db"
	"Cloudocs/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func GetUsers(ctx *gin.Context) {
	skipStr := ctx.DefaultQuery("skip", "0")
	limitStr := ctx.DefaultQuery("limit", "20")
	skip, err := strconv.Atoi(skipStr)
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "数据非法",
			"err": err.Error(),
		})
		return
	}
	// 查找用户
	users := db.Users.FindAll(skip, limit)
	// 返回用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func GetUsersToken(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	// 查找用户
	user, err := db.Users.FindId(fmt.Sprint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
			"id":  id,
		})
		return
	}
	// 返回用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func GetUsersId(ctx *gin.Context) {
	id := ctx.Param("id")
	// 查找用户
	user, err := db.Users.FindId(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
			"id":  id,
		})
		return
	}
	// 返回用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func PostUsers(ctx *gin.Context) {
	// 获取参数
	tel := ctx.PostForm("tel")
	pass := ctx.PostForm("pass")
	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	genderStr := ctx.DefaultPostForm("gender", "0")
	gender, err := strconv.Atoi(genderStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "数据非法",
			"err": err.Error(),
		})
		return
	}

	if tel == "" || pass == "" || name == "" {
		ctx.JSON(http.StatusPreconditionFailed, gin.H{
			"msg": "手机号,密码,用户名不能为空,创建用户失败",
		})
		return
	}
	// 读取数据库
	_, err = db.Users.Find("tel", tel)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"msg": "手机号已存在,创建用户失败",
		})
		return
	}
	passBcrypt, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Bcrypt错误,创建用户失败",
			"err": err.Error(),
		})
	}

	objectId := bson.NewObjectId()
	user := model.User{
		Id:      objectId,
		Name:    name,
		Pass:    string(passBcrypt),
		Tel:     tel,
		Email:   email,
		Gender:  int8(gender),
		Created: objectId.Time().Unix(),
		Last:    objectId.Time().Unix(),
	}
	err = db.Users.Insert(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "MongoDB错误,创建用户失败",
			"err": err.Error(),
		})
		return
	}
	// 获取Token
	token, err := common.ReleaseToken(objectId.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "创建用户成功,但获取Token失败",
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"msg":   "创建用户成功",
		"token": token,
		"user":  user,
	})
}
