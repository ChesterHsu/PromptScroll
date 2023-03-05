package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"promptscroll/model"
	"promptscroll/service"
)

type UserController interface {
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

// 處理獲取使用者的請求
func (ctl *userController) GetUser(c *gin.Context) {
	// 從 URL 中獲取使用者 ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	// 使用 UserService 獲取使用者
	user, err := ctl.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user",
		})
		return
	}

	// 返回使用者資訊
	c.JSON(http.StatusOK, user)
}

// 處理建立新使用者的請求
func (ctl *userController) CreateUser(c *gin.Context) {
	// 從 HTTP 請求中解析出使用者資訊
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// 使用 UserService 創建新使用者
	newUser, err := ctl.userService.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	// 返回新建使用者的資訊
	c.JSON(http.StatusOK, newUser)
}
