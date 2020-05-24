package users

import (
	"datastore-gin-example/common"
	"net/http"

	"github.com/fatih/structs"

	"github.com/gin-gonic/gin"
)

func UserRegister(router *gin.RouterGroup) {
	router.GET("/all", FetchAllUsers)
	router.POST("/save", SaveUser)
	router.PUT("/update/:userID", Update)
	router.GET("/get/:userID", GetUser)
	router.DELETE("/delete/:userID", Delete)
}

func LoginRegister(router *gin.RouterGroup) {
	router.POST("/login", Login)
}

func FetchAllUsers(c *gin.Context) {
	users, err := GetAllUsers()
	if err != nil {
		common.LogError("error: ", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": users,
	})
}

func GetUser(c *gin.Context) {
	var err error
	userID := c.Param("userID")
	common.LogInfo("find user by id: " + userID)
	user, err := GetUserByID(userID)
	if err != nil {
		common.LogError("error: ", err)
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": user,
	})
}

func SaveUser(c *gin.Context) {
	user := User{}
	if c.ShouldBind(&user) != nil {
		c.Status(http.StatusBadRequest)
	}
	common.LogInfo("Save user ", structs.Map(user))
	err := InsertUser(&user)
	if err != nil {
		common.LogError("Save user error: ", err)
		c.Status(http.StatusInternalServerError)
	} else {
		c.String(http.StatusOK, "Success")
	}
}

func Update(c *gin.Context) {
	user := User{}
	userID := c.Param("userID")
	if c.ShouldBind(&user) != nil {
		c.Status(http.StatusBadRequest)
	}

	common.LogInfo("update user id: " + userID)
	err := UpdateUser(userID, &user)
	if err != nil {
		common.LogError("Update user error: ", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, "Success")

}

func Delete(c *gin.Context) {
	userID := c.Param("userID")
	err := DeleteUser(userID)
	if err != nil {
		common.LogError("Delete user error: ", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, "Success")
}

func Login(c *gin.Context) {
	loginAccount := LoginAccount{}
	common.LogInfo("user login. ")
	if c.ShouldBind(&loginAccount) != nil {
		c.Status(http.StatusBadRequest)
	}
	user, account, err := UserLogin(&loginAccount)
	if user == nil || account == nil {
		c.Status(401)
		return
	}
	tokenString, err := common.CreateToken(user.Name, account.Account)
	if err != nil {
		c.Status(500)
		return
	}
	c.Writer.Header().Set("Authorization", "Bearer "+tokenString)
	c.Status(http.StatusOK)
}
