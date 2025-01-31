package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const CookieName = "ice_session_id"

func (*Server) GenerateCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie(CookieName)
	if errors.Is(err, http.ErrNoCookie) {
		c.SetCookie(CookieName, time.Now().String(), 3600, "/", "localhost", false, true)
	} else {
		c.Set(CookieName, cookie)
	}
	c.Next()
}

func (*Server) RequiredCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie(CookieName)
	if err != nil || errors.Is(err, http.ErrNoCookie) || (cookie != nil && cookie.Value == "") {
		c.Redirect(302, "/")
		return
	}
	c.Set(CookieName, cookie)
	c.Next()
}
