package controllers

import "github.com/gin-gonic/gin"

func RegisterUser(c *gin.Context) {
	res := struct {
		Message string
	}{
		Message: "ok",
	}
	c.JSON(200, res)
}
