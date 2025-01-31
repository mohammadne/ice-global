package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CookieName = "ice_session_id"
	UserKey    = "user_id"
)

func (s *Server) OptionalCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie(CookieName)
	if errors.Is(err, http.ErrNoCookie) {
		cookieValue := strconv.FormatInt(time.Now().UnixNano(), 10)
		c.SetCookie(CookieName, cookieValue, 3600, "/", "localhost", false, true)
	} else {
		// find the user and if exists, put the user-id
		user, err := s.usersService.RetrieveUserOptional(c.Request.Context(), cookie.Value)
		if err != nil {
			slog.Error("error retrieving user", "Err", err)
		} else if user != nil {
			c.Set(UserKey, user.Id)
		}
	}
	c.Next()
}

func (s *Server) RequiredCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie(CookieName)
	if err != nil || errors.Is(err, http.ErrNoCookie) || (cookie != nil && cookie.Value == "") {
		c.Redirect(302, "/")
		return
	}

	user, err := s.usersService.RetrieveUserRequired(c.Request.Context(), cookie.Value)
	if err != nil {
		c.Redirect(302, "/")
		return
	}

	c.Set(UserKey, user.Id)
	c.Next()
}
