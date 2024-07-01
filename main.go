package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null"`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	defer db.Close()

	//创建一个服务
	ginServe := gin.Default()

	ginServe.POST("/user/add/:username/:password", func(c *gin.Context) {
		username := c.Param("username")
		password := c.Param("password")
		c.JSON(http.StatusOK, gin.H{
			"msg":      "ok",
			"username": username,
			"password": password,
		})
	})

	//restfulApi
	ginServe.POST("/api/auth/register", func(c *gin.Context) {
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
			name = RandomString(10)
		}
		//判断手机号是否存在
		log.Println(name, telephone, password)
		if isTelephoneExist(db, telephone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户已经存在",
			})
			return
		}
		//创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)
		//返回结果
		c.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	ginServe.Run()
}

// 随机字符串
func RandomString(n int) string {
	var letters = []byte("arfshghijklmnABCFADBREHIJKLMN")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 初始化MySQL数据库
func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "123456"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
