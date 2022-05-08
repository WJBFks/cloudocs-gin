package controller

import (
	"Cloudocs/db"
	"Cloudocs/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func PostDocs(ctx *gin.Context) {
	// 获取参数
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")
	Iid, isExist := ctx.Get("id")
	id := fmt.Sprint(Iid)
	if !isExist {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
			"id":  id,
		})
		return
	}
	// 查找用户
	user, err := db.Users.FindId(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
			"id":  id,
		})
		return
	}
	objectId := bson.NewObjectId()
	doc := model.Docs{
		Id:          objectId,
		Title:       title,
		Content:     content,
		Creator:     id,
		CreatorName: user.Name,
		Created:     objectId.Time().Unix(),
		Last:        objectId.Time().Unix(),
		Openid:      user.Openid,
	}
	err = db.Docs.Insert(doc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统异常，创建失败",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"doc":     doc,
		"content": doc.Content,
	})
}

func GetDocs(ctx *gin.Context) {
	// 获取参数
	Iid, isExist := ctx.Get("id")
	id := fmt.Sprint(Iid)
	if !isExist {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
			"id":  id,
		})
		return
	}
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
	_, err = db.Users.FindId(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
			"id":  id,
		})
		return
	}
	docs := db.Docs.Finds(id, skip, limit)
	for index, doc := range docs {
		user, err := db.Users.FindId(doc.Creator)
		if err == nil {
			docs[index].CreatorName = user.Name
			docs[index].Openid = user.Openid
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"docs": docs,
	})
}

func GetDocsId(ctx *gin.Context) {
	// 获取参数
	docId := ctx.Param("id")
	// 查找用户
	userId, isExist := ctx.Get("id")
	if !isExist {
		return
	}
	doc, err := db.Docs.FindId(docId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":   "文档不存在",
			"docId": docId,
		})
		return
	}
	if doc.Creator != userId {
		user, err := db.Users.FindId(fmt.Sprint(userId))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"msg":   "用户不存在",
				"docId": docId,
			})
			return
		}
		authority := false
		for _, value := range user.ShareDocs {
			if value == docId {
				authority = true
				break
			}
		}
		if !authority {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg":   "权限不足",
				"docId": docId,
			})
			return
		}
	}
	user, err := db.Users.FindId(doc.Creator)
	if err == nil {
		doc.CreatorName = user.Name
		doc.Openid = user.Openid
	}
	ctx.JSON(http.StatusOK, gin.H{
		"doc":     doc,
		"content": doc.Content,
	})
}

func PutDocsId(ctx *gin.Context) {
	// 获取参数
	docId := ctx.Param("id")
	content := ctx.PostForm("content")
	title := ctx.PostForm("title")
	// 查找用户
	userId, isExist := ctx.Get("id")
	if !isExist {
		return
	}
	doc, err := db.Docs.FindId(docId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":   "文档不存在",
			"docId": docId,
		})
		return
	}
	if doc.Creator != userId {
		user, err := db.Users.FindId(fmt.Sprint(userId))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"msg":   "用户不存在",
				"docId": docId,
			})
			return
		}
		authority := false
		for _, value := range user.ShareDocs {
			if value == docId {
				authority = true
				break
			}
		}
		if !authority {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg":   "权限不足",
				"docId": docId,
			})
			return
		}
	}
	if content != "" {
		err = db.Docs.UpdateContent(docId, content)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "文档内容修改异常",
			})
			return
		}
	}
	if title != "" {
		err = db.Docs.UpdateTitle(docId, title)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "文档内容修改异常",
			})
			return
		}
	}
	doc, err = db.Docs.FindId(docId)
	if err != nil {
		ctx.JSON(900, gin.H{
			"msg": "文档内容修改成功，但获取更新后文档数据异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"doc":     doc,
		"content": doc.Content,
	})
}

func DelDoc(ctx *gin.Context) {
	// 获取参数
	id := ctx.Param("id")
	// 查找用户
	user, isExist := ctx.Get("id")
	if !isExist {
		return
	}
	doc, err := db.Docs.FindId(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "文档不存在",
			"id":  id,
		})
		return
	}
	if doc.Creator != user {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "权限不足",
			"id":  id,
		})
		return
	}
	err = db.Docs.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "后端错误，删除失败",
			"id":  id,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
		"id":  id,
	})
}

func AddShareDoc(ctx *gin.Context) {
	// 获取参数
	id := ctx.Param("id")
	tel := ctx.PostForm("tel")
	// 查找用户
	if tel == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "tel不能为空",
			"id":  id,
			"tel": tel,
		})
		return
	}
	creator, isExist := ctx.Get("id")
	if !isExist {
		return
	}
	doc, err := db.Docs.FindId(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":     "文档不存在",
			"creator": creator,
			"id":      id,
			"tel":     tel,
		})
		return
	}
	if doc.Creator != creator {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg":     "权限不足",
			"creator": creator,
			"id":      id,
			"tel":     tel,
		})
		return
	}
	var user model.User
	user, err = db.Users.Find("tel", tel)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":     "用户不存在",
			"creator": creator,
			"id":      id,
			"tel":     tel,
		})
		return
	}
	if user.Id.Hex() == creator {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":     "不能给自己分享",
			"creator": creator,
			"id":      id,
			"tel":     tel,
		})
		return
	}
	err = db.Users.AddShareDocs(user.Id, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "后端错误，分享失败",
			"id":  id,
			"tel": tel,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "分享添加成功",
		"id":  id,
		"tel": tel,
	})
}

func GetShareDocs(ctx *gin.Context) {
	id, isExist := ctx.Get("id")
	if !isExist {
		return
	}
	user, err := db.Users.FindId(fmt.Sprint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
			"id":  id,
		})
		return
	}
	var docs []model.Docs = make([]model.Docs, 0)
	var newShareDocs []string = make([]string, 0)
	for _, value := range user.ShareDocs {
		doc, err := db.Docs.FindId(value)
		if err == nil {
			docs = append(docs, doc)
			newShareDocs = append(newShareDocs, value)
		}
	}
	db.Users.UpdateShareDocs(user.Id, newShareDocs)
	for index, doc := range docs {
		user, err := db.Users.FindId(doc.Creator)
		if err == nil {
			docs[index].CreatorName = user.Name
			docs[index].Openid = user.Openid
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"docs": docs,
	})
}

func GetDocShare(ctx *gin.Context) {
	docId := ctx.Param("id")
	userId, isExist := ctx.Get("id")
	if !isExist {
		return
	}
	_, err := db.Users.FindId(fmt.Sprint(userId))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":    "用户不存在",
			"userId": userId,
		})
		return
	}
	doc, err := db.Docs.FindId(fmt.Sprint(docId))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":    "文档不存在",
			"userId": userId,
			"docId":  docId,
		})
		return
	}
	if doc.Creator != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg":    "权限不足",
			"userId": userId,
			"docId":  docId,
		})
		return
	}
	type Response struct {
		Tel    string `json:"tel"`
		Name   string `json:"name"`
		Openid string `json:"openid,omitempty"`
	}
	responses := make([]Response, 0)
	for _, value := range doc.Share {
		user, err := db.Users.Find("tel", value)
		if err == nil {
			response := Response{}
			response.Name = user.Name
			response.Tel = user.Tel
			response.Openid = user.Openid
			responses = append(responses, response)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"users": responses,
	})
}

func DelDocShareTel(ctx *gin.Context) {
	docId := ctx.Param("id")
	userId, isExist := ctx.Get("id")
	tel := ctx.Query("tel")
	// 查找用户
	if tel == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":    "tel不能为空",
			"userId": userId,
			"tel":    tel,
		})
		return
	}
	if !isExist {
		return
	}
	_, err := db.Users.FindId(fmt.Sprint(userId))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":    "用户不存在",
			"userId": userId,
			"tel":    tel,
		})
		return
	}
	doc, err := db.Docs.FindId(fmt.Sprint(docId))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":    "文档不存在",
			"userId": userId,
			"docId":  docId,
			"tel":    tel,
		})
		return
	}
	if doc.Creator != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg":    "权限不足",
			"userId": userId,
			"docId":  docId,
			"tel":    tel,
		})
		return
	}
	err = db.Users.DelShareDocs(docId, tel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "后端错误，删除失败",
			"userId": userId,
			"docId":  docId,
			"tel":    tel,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "删除成功",
		"userId": userId,
		"docId":  docId,
		"tel":    tel,
	})
	return
}
