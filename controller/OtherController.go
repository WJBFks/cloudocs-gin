package controller

import (
	"Cloudocs/common"
	"Cloudocs/db"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *gin.Context) {
	// 获取参数
	tel := ctx.PostForm("tel")
	pass := ctx.PostForm("pass")
	// 查找用户
	user, err := db.Users.Find("tel", tel)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "账号或密码错误，登录失败",
			"err": err.Error(),
		})
		return
	}
	// 如果密码错误
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "账号或密码错误，登录失败",
			"err": err.Error(),
		})
		return
	}
	// 登录成功，获取Token
	token, err := common.ReleaseToken(user.Id.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统异常，登录失败",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"token": token,
		"user":  user,
	})
}

func Bind(ctx *gin.Context) {
	// 获取参数
	openid := ctx.PostForm("openid")
	tel := ctx.PostForm("tel")
	pass := ctx.PostForm("pass")
	// ctx.JSON(http.StatusNotFound, gin.H{
	// 	"openid": openid,
	// 	"tel":    tel,
	// 	"pass":   pass,
	// })
	// return
	// 查找用户
	user, err := db.Users.Find("tel", tel)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "账号或密码错误，登录失败",
			"err": err.Error(),
		})
		return
	}
	// 如果密码错误
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "账号或密码错误，登录失败",
			"err": err.Error(),
		})
		return
	}
	// 登录成功，获取Token
	user.Openid = openid
	db.Users.Update(user)
	token, err := common.ReleaseToken(user.Id.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统异常，登录失败",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"token": token,
		"user":  user,
	})
}

func Bind2(ctx *gin.Context) {
	// 获取参数
	openid := ctx.PostForm("openid")
	tel := ctx.PostForm("tel")
	wxName := ctx.PostForm("wxName")
	// 查找用户
	user, err := db.Users.Find("tel", tel)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "该手机号暂未注册通用账号",
			"err": err.Error(),
		})
		return
	}
	user.Openid = "@" + wxName + "@" + openid
	err = db.Users.Update(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误，绑定失败",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "申请绑定成功",
	})
}

func DbInsert(ctx *gin.Context) {
	// 获取参数
	data := ctx.PostForm("data")
	D := &gin.H{}
	err := json.Unmarshal([]byte(data), D)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Unmarshal错误",
			"err": err.Error(),
		})
		return
	}
	err = db.TestInsert(D)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "MongoDB错误",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "添加成功",
		"data": D,
	})
}

func DbFinds(ctx *gin.Context) {
	// 获取参数
	data := db.TestFinds()
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func DbFindId(ctx *gin.Context) {
	id := ctx.Param("id")
	// 获取参数
	data, err := db.TestFind(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "MongoDB错误",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
