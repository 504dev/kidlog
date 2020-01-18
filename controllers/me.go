package controllers

import (
	"fmt"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/user"
	"github.com/gin-gonic/gin"
)

type MeController struct{}

func (_ MeController) Me(c *gin.Context) {
	usr, _ := user.GetById(1)
	c.JSON(200, usr)
}
func (_ MeController) Dashboards(c *gin.Context) {
	jwt, _ := c.Get("jwt")
	fmt.Println("jwt", jwt)
	dashboards, _ := dashboard.GetAll()
	c.JSON(200, dashboards)
}
