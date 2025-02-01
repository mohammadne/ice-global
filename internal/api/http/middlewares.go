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
	CartKey    = "cart_id"
)

func (s *Server) OptionalCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie(CookieName)
	if errors.Is(err, http.ErrNoCookie) {
		cookieValue := strconv.FormatInt(time.Now().UnixNano(), 10)
		c.SetCookie(CookieName, cookieValue, 3600, "/", "localhost", false, true)
	} else {
		// find the cart and if exists, put the cart-id
		cart, err := s.cartsService.RetrieveCartOptional(c.Request.Context(), cookie.Value)
		if err != nil {
			slog.Error("error retrieving cart", "Err", err)
		} else if cart != nil {
			c.Set(CartKey, cart.Id)
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

	cart, err := s.cartsService.RetrieveCartRequired(c.Request.Context(), cookie.Value)
	if err != nil {
		c.Redirect(302, "/")
		return
	}

	c.Set(CartKey, cart.Id)
	c.Next()
}
