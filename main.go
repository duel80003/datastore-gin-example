package main

import (
	"datastore-gin-example/common"
	_ "datastore-gin-example/common"
	"datastore-gin-example/users"

	"github.com/gin-gonic/gin"
)

func main() {
	defer common.DatastoreClient.Close()
	r := gin.Default()
	auth := r.Group("/")
	users.LoginRegister(auth)
	v1 := r.Group("/api")
	v1.Use(common.AuthenticationToken())
	users.UserRegister(v1.Group("user"))
	r.Run("localhost:8080")
}
