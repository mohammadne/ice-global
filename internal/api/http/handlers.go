package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mohammadne/ice-global/internal/calculator"
)

func (*Server) liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (*Server) readiness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *Server) ShowAddItemForm(c *gin.Context) {
	_, err := c.Request.Cookie("ice_session_id")
	if errors.Is(err, http.ErrNoCookie) {
		c.SetCookie("ice_session_id", time.Now().String(), 3600, "/", "localhost", false, true)
	}

	calculator.GetCartData(c)
}

func (s *Server) AddItem(c *gin.Context) {
	cookie, err := c.Request.Cookie("ice_session_id")

	if err != nil || errors.Is(err, http.ErrNoCookie) || (cookie != nil && cookie.Value == "") {
		c.Redirect(302, "/")
		return
	}

	calculator.AddItemToCart(c)
}

func (s *Server) DeleteCartItem(c *gin.Context) {
	cookie, err := c.Request.Cookie("ice_session_id")

	if err != nil || errors.Is(err, http.ErrNoCookie) || (cookie != nil && cookie.Value == "") {
		c.Redirect(302, "/")
		return
	}

	calculator.DeleteCartItem(c)
}
