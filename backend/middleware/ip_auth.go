package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hwg1999/go_gateway/backend/golang_common/lib"
)

func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		for _, host := range lib.GetStringSliceConf("base.http.allow_ip") {
			if c.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			ResponseError(c, InternalErrorCode, fmt.Errorf("%v, not in iplist", c.ClientIP()))
			c.Abort()
			return
		}
		c.Next()
	}
}
