package app

import (
	"github.com/gin-gonic/gin"
)

func BindAndValid(c *gin.Context, form interface{}) error {
	if err := c.ShouldBind(form); err != nil {
		return err
	}
	return nil
}
