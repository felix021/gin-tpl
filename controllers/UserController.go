package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

// curl http://127.0.0.1:8080/user/
// curl http://127.0.0.1:8080/user/index
func (this UserController) GetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, fmt.Sprintf("Index\n"))
	}
}

// curl http://127.0.0.1:8080/user/detail/1
func (this UserController) GetDetail() (gin.HandlerFunc, string) {
	return func(c *gin.Context) {
		c.String(200, fmt.Sprintf("Detail: %s\n", c.Param("id")))
	}, "/:id"
}

// curl http://127.0.0.1:8080/user/register/felix/mypwd -i -d ''
func (this UserController) PostRegister() (gin.HandlerFunc, string) {
	return func(c *gin.Context) {
		c.String(200, fmt.Sprintf("Register: username=%s, password=%s\n", c.Param("username"), c.Param("password")))
	}, "/:username/:password"
}
