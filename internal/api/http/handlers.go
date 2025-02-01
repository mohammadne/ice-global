package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mohammadne/ice-global/internal/services"
)

func (*Server) liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (*Server) readiness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *Server) showAddItemForm(c *gin.Context) {
	items, err := s.itemsService.AllItems(c.Request.Context())
	if err != nil {
		slog.Error("error retrieving all-items")
		c.Status(http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Error": c.Query("error"),
		"Items": items,
	}

	if cartId, exists := c.Get(CartKey); exists {
		cartItems, err := s.cartsService.AllCartItemsByCartId(c.Request.Context(), cartId.(int))
		if err != nil {
			slog.Error("error while retrieving cart-items for the user", "Err", err)
		} else {
			result := make([]map[string]any, 0, len(cartItems))
			for _, cartItem := range cartItems {
				if cartItem.Quantity <= 0 {
					continue
				}
				resultItem := map[string]any{
					"ID":       cartItem.Id,
					"Quantity": cartItem.Quantity,
					"Price":    cartItem.Item.Price * cartItem.Quantity,
					"Product":  cartItem.Item.Name,
				}
				result = append(result, resultItem)
			}
			data["CartItems"] = result
		}
	}

	c.HTML(http.StatusOK, "index.html", data)
}

// -----------------------------------------------------

type CartItemForm struct {
	ItemId   int    `form:"product"   binding:"required"`
	Quantity string `form:"quantity"  binding:"required"`
}

func (s *Server) addItem(c *gin.Context) {
	cartId, exists := c.Get(CartKey)
	if !exists {
		c.Redirect(302, "/?error="+"user is not authorized")
		return
	}

	if c.Request.Body == nil {
		c.Redirect(302, "/?error="+"body can not be empty")
		return
	}

	addItemForm := &CartItemForm{}
	if err := binding.FormPost.Bind(c.Request, addItemForm); err != nil {
		c.Redirect(302, "/?error="+err.Error())
		return
	}

	quantity, err := strconv.Atoi(addItemForm.Quantity)
	if err != nil {
		c.Redirect(302, "/?error=invalid quantity")
		return
	}

	err = s.cartsService.AddItemToCart(c.Request.Context(), cartId.(int), addItemForm.ItemId, quantity)
	if err != nil {
		slog.Error("", "Err", err)
		c.Redirect(302, "/?error=internal error for adding item to the cart")
		return
	}

	c.Redirect(302, "/")
}

// -----------------------------------------------------

func (s *Server) deleteCartItem(c *gin.Context) {
	cartId, exists := c.Get(CartKey)
	if !exists {
		c.Redirect(302, "/?error="+"user is not authorized")
		return
	}

	cartItemIdString := c.Query("cart_item_id")
	if cartItemIdString == "" {
		c.Redirect(302, "/")
		return
	}

	cartItemId, err := strconv.Atoi(cartItemIdString)
	if err != nil {
		c.Redirect(302, "/")
		return
	}

	err = s.cartsService.DeleteCartItem(c.Request.Context(), cartId.(int), cartItemId)
	if err != nil {
		if err == services.ErrorCartHasBeenClosed {
			c.Redirect(302, "/?error="+"you can't remove from an closed cart")
			return
		}
		c.Redirect(302, "/?error="+"cart-item has not been removed")
		return
	}

	c.Redirect(302, "/")
}
