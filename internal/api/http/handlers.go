package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*Server) liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (*Server) readiness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *Server) showAddItemForm(c *gin.Context) {
	data := map[string]any{
		"Error": c.Query("error"),
	}

	if cartId, exists := c.Get(CartKey); exists {
		cartItems, err := s.cartsService.AllCartItemsByCartId(c.Request.Context(), cartId.(int))
		if err != nil {
			slog.Error("error while retrieving cart-items for the user", "Err", err)
		} else {
			result := make([]map[string]any, len(cartItems))
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

// var itemPriceMapping = map[string]float64{
// 	"shoe":  100,
// 	"purse": 200,
// 	"bag":   300,
// 	"watch": 300,
// }

// type CartItemForm struct {
// 	Product  string `form:"product"   binding:"required"`
// 	Quantity string `form:"quantity"  binding:"required"`
// }

// func (s *Server) addItem(c *gin.Context) {
// 	cookieRaw, exists := c.Get(CookieName)
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, nil) // TODO
// 		return
// 	}
// 	cookie := cookieRaw.(*http.Cookie)

// 	db := db.GetDatabase()

// 	var isCartNew bool
// 	var cartEntity entity.CartEntity
// 	result := db.Where(fmt.Sprintf("status = '%s' AND session_id = '%s'", entity.CartOpen, cookie.Value)).First(&cartEntity)

// 	if result.Error != nil {
// 		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			c.Redirect(302, "/")
// 			return
// 		}
// 		isCartNew = true
// 		cartEntity = entity.CartEntity{
// 			SessionID: cookie.Value,
// 			Status:    entity.CartOpen,
// 		}
// 		db.Create(&cartEntity)
// 	}

// 	addItemForm, err := getCartItemForm(c)
// 	if err != nil {
// 		c.Redirect(302, "/?error="+err.Error())
// 		return
// 	}

// 	item, ok := itemPriceMapping[addItemForm.Product]
// 	if !ok {
// 		c.Redirect(302, "/?error=invalid item name")
// 		return
// 	}

// 	quantity, err := strconv.ParseInt(addItemForm.Quantity, 10, 0)
// 	if err != nil {
// 		c.Redirect(302, "/?error=invalid quantity")
// 		return
// 	}

// 	var cartItemEntity entity.CartItem
// 	if isCartNew {
// 		cartItemEntity = entity.CartItem{
// 			CartID:      cartEntity.ID,
// 			ProductName: addItemForm.Product,
// 			Quantity:    int(quantity),
// 			Price:       item * float64(quantity),
// 		}
// 		db.Create(&cartItemEntity)
// 	} else {
// 		result = db.Where(" cart_id = ? and product_name  = ?", cartEntity.ID, addItemForm.Product).First(&cartItemEntity)

// 		if result.Error != nil {
// 			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 				c.Redirect(302, "/")
// 				return
// 			}
// 			cartItemEntity = entity.CartItem{
// 				CartID:      cartEntity.ID,
// 				ProductName: addItemForm.Product,
// 				Quantity:    int(quantity),
// 				Price:       item * float64(quantity),
// 			}
// 			db.Create(&cartItemEntity)

// 		} else {
// 			cartItemEntity.Quantity += int(quantity)
// 			cartItemEntity.Price += item * float64(quantity)
// 			db.Save(&cartItemEntity)
// 		}
// 	}

// 	c.Redirect(302, "/")
// }

// func getCartItemForm(c *gin.Context) (*CartItemForm, error) {
// 	if c.Request.Body == nil {
// 		return nil, fmt.Errorf("body cannot be nil")
// 	}

// 	form := &CartItemForm{}

// 	if err := binding.FormPost.Bind(c.Request, form); err != nil {
// 		log.Println(err.Error())
// 		return nil, err
// 	}

// 	return form, nil
// }

// -----------------------------------------------------

// func (s *Server) deleteCartItem(c *gin.Context) {
// 	cookieRaw, exists := c.Get(CookieName)
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, nil) // TODO
// 		return
// 	}
// 	cookie := cookieRaw.(*http.Cookie)

// 	cartItemIDString := c.Query("cart_item_id")
// 	if cartItemIDString == "" {
// 		c.Redirect(302, "/")
// 		return
// 	}

// 	db := db.GetDatabase()

// 	var cartEntity entity.CartEntity
// 	result := db.Where(fmt.Sprintf("status = '%s' AND session_id = '%s'", entity.CartOpen, cookie.Value)).First(&cartEntity)
// 	if result.Error != nil {
// 		c.Redirect(302, "/")
// 		return
// 	}

// 	if cartEntity.Status == entity.CartClosed {
// 		c.Redirect(302, "/")
// 		return
// 	}

// 	cartItemID, err := strconv.Atoi(cartItemIDString)
// 	if err != nil {
// 		c.Redirect(302, "/")
// 		return
// 	}

// 	var cartItemEntity entity.CartItem

// 	result = db.Where(" ID  = ?", cartItemID).First(&cartItemEntity)
// 	if result.Error != nil {
// 		c.Redirect(302, "/")
// 		return
// 	}

// 	db.Delete(&cartItemEntity)
// 	c.Redirect(302, "/")
// }
