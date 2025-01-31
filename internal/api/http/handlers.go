package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mohammadne/ice-global/internal/calculator"
	"github.com/mohammadne/ice-global/internal/db"
	"github.com/mohammadne/ice-global/internal/entity"
)

func (*Server) liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (*Server) readiness(c *gin.Context) {
	c.Status(http.StatusOK)
}

// ------------------------------------------------------

func (s *Server) showAddItemForm(c *gin.Context) {
	data := map[string]any{
		"Error": c.Query("error"),
	}

	cookie, exists := c.Get(CookieName)
	if exists {
		data["CartItems"] = getCartItemData(cookie.(*http.Cookie).Value)
	}

	c.HTML(http.StatusOK, "index.html", data)
}

func getCartItemData(sessionID string) (items []map[string]interface{}) {
	db := db.GetDatabase()
	var cartEntity entity.CartEntity
	result := db.Where(fmt.Sprintf("status = '%s' AND session_id = '%s'", entity.CartOpen, sessionID)).First(&cartEntity)

	if result.Error != nil {
		return
	}

	var cartItems []entity.CartItem
	result = db.Where(fmt.Sprintf("cart_id = %d", cartEntity.ID)).Find(&cartItems)
	if result.Error != nil {
		return
	}

	for _, cartItem := range cartItems {
		item := map[string]interface{}{
			"ID":       cartItem.ID,
			"Quantity": cartItem.Quantity,
			"Price":    cartItem.Price,
			"Product":  cartItem.ProductName,
		}

		items = append(items, item)
	}
	return items
}

// -----------------------------------------------------

func (s *Server) addItem(c *gin.Context) {
	calculator.AddItemToCart(c)
}

// -----------------------------------------------------

func (s *Server) deleteCartItem(c *gin.Context) {
	calculator.DeleteCartItem(c)
}
