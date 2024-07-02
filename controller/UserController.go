package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := c.Query("name")
	telephone := c.Query("telephone")
	password := c.Query("password")
	fmt.Println(name, telephone, password)
	//数据验证
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "手机号必须为11位",
		})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码不可少于6位",
		})
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}
	//判断手机号是否存在
	log.Println(name, telephone, password)
	if isTelephoneExist(DB, telephone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户已经存在",
		})
		return
	}
	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加密错误"})
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)
	//返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

func Login(c *gin.Context) {
	//获取参数
	telephone := c.Query("telephone")
	password := c.Query("password")
	//数据验证
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "手机号必须为11位",
		})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码不可少于6位",
		})
		return
	}

	//判断手机号是否存在
	var user model.User
	DB := common.GetDB()
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户不存在",
		})
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发放token
	token := "11"

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
