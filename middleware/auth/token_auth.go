package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nacos-service-discovery-controller/pkg/app"
	"nacos-service-discovery-controller/pkg/errcode"
)

// curl -X GET "http://127.0.0.1:8000/api/v1/tags?pageNumber=1&pageSize=10"
// curl -X GET "http://127.0.0.1:8000/api/v1/tags?pageNumber=1&pageSize=10" -H 'token:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBLZXkiOiJ6cXlhbmdjaG4iLCJhcHBTZWNyZXQiOiJnby1wcm9tZ3JhbW1pbmctdG91ci1ib29rIiwiZXhwIjoxNTk4NTAwMzIwLCJpc3MiOiJodHRwLXNlcnZpY2UifQ.HxA0waRpdCpeiSK2P1qyFgBcz4O3kP_chrbF2UJ1oLY'

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token   string
			errCode = errcode.Success
		)

		if t, exist := c.GetQuery("token"); exist { // query token
			token = t
		}
		if t := c.GetHeader("token"); t != "" { // header token
			token = t
		}
		if t := c.GetHeader("X-Gitlab-Token"); t != "" { // header X-Gitlab-Token for gitlab
			token = t
		}

		if token != "BPsGHoO3SyIviSwj" {
			errCode = errcode.TokenAuthError
		}

		if errCode.Code != errcode.Success.Code {
			appG := app.Gin{Context: c}
			appG.ResponseError(http.StatusUnauthorized, errCode)
			c.Abort()
			return
		}

		c.Next()
	}
}
