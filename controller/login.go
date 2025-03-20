package controller

import (
	"java2go/entity"
	"java2go/mapper"
	"java2go/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 处理 /permission_denied 请求
func PermissionDenied(c *gin.Context) {
	utils.Make403Resp("Permission denied")
}

// 处理 /getLoginStatus 请求
func GetLoginStatus(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	userRole := session.Get("user_role")
	lastLogin := session.Get("last_login")
	if username != nil && userRole != nil && lastLogin != nil {
		retData := map[string]interface{}{
			"username":   username,
			"user_role":  userRole,
			"last_login": lastLogin,
		}
		response := utils.Make200Resp("Success", retData)
		c.JSON(http.StatusOK, response)
	} else {
		str := utils.Make403Resp("Permission denied")
		c.JSON(http.StatusOK, str)
	}
}

// 处理 /Login 请求
func Login(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.LastLogin = time.Now()
	var userList []entity.User
	mapper.DB.Where("username = ? AND password = ?", user.Username, user.Password).Find(&userList)
	if len(userList) > 0 {
		mapper.DB.Model(&user).Where("username = ?", user.Username).Update("last_login", user.LastLogin)
		session := sessions.Default(c)
		session.Set("username", userList[0].Username)
		session.Set("user_role", userList[0].UserRole)
		session.Set("last_login", userList[0].LastLogin)
		session.Save()
		retData := map[string]interface{}{
			"username":   userList[0].Username,
			"user_role":  userList[0].UserRole,
			"last_login": userList[0].LastLogin,
		}
		response := map[string]interface{}{
			"status":  200,
			"message": "Success",
			"data":    retData,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := map[string]interface{}{
			"status":  403,
			"message": "Permission denied",
		}
		c.JSON(http.StatusForbidden, response)
	}
}

// 处理 /Logout 请求
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    "退出登陆成功",
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /Registered 请求
func Registered(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userByUsername []entity.User
	mapper.DB.Where("username = ?", user.Username).Find(&userByUsername)
	if len(userByUsername) > 0 {
		response := map[string]interface{}{
			"status":  500,
			"message": "用户名重复",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	user.LastLogin = time.Now()
	user.Enable = 0
	result := mapper.DB.Create(&user)
	if result.Error != nil {
		response := map[string]interface{}{
			"status":  500,
			"message": "注册失败",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    result.RowsAffected,
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /GetApplyUser 请求
func GetApplyUser(c *gin.Context) {
	session := sessions.Default(c)
	userRole := session.Get("user_role")
	if userRole == nil || userRole.(string) != "admin" {
		response := map[string]interface{}{
			"status":  403,
			"message": "Permission denied",
		}
		c.JSON(http.StatusForbidden, response)
		return
	}
	var applyUser []entity.User
	mapper.DB.Where("enable = 0").Find(&applyUser)
	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    applyUser,
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /GetAllUser 请求
func GetAllUser(c *gin.Context) {
	var allUser []entity.User
	mapper.DB.Find(&allUser)
	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    allUser,
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /DeleteUser 请求
func DeleteUser(c *gin.Context) {
	username := c.Query("username")
	result := mapper.DB.Where("username = ?", username).Delete(&entity.User{})
	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    result.RowsAffected,
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /PassApply 请求
func PassApply(c *gin.Context) {
	username := c.Query("username")
	result := mapper.DB.Model(&entity.User{}).Where("username = ?", username).Update("enable", 1)
	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    result.RowsAffected,
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /DeleteApply 请求
func DeleteApply(c *gin.Context) {
	username := c.Query("username")
	result := mapper.DB.Where("username = ? AND enable = 0", username).Delete(&entity.User{})
	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    result.RowsAffected,
	}
	c.JSON(http.StatusOK, response)
}
