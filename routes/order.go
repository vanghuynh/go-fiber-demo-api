package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vanghuynh/fiber-api/database"
	"github.com/vanghuynh/fiber-api/models"
)

type OrderDto struct {
	ID        uint       `json:"id"`
	User      UserDto    `json:"user"`
	Product   ProductDto `json:"product"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateOrderDto struct {
	UserId    uint `json:"user_id"`
	ProductId uint `json:"product_id"`
}

func CreateOrderResponse(order models.Order, user models.User, product models.Product) OrderDto {
	return OrderDto{
		ID:        order.ID,
		User:      CreateResponseUser(user),
		Product:   CreateProductResponse(product),
		CreatedAt: order.CreatedAt,
	}
}

// CreateOrder godoc
// @Summary      Create order
// @Description  create order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        data body CreateOrderDto true "The input order struct"
// @Success      200  {object}  OrderDto
// @Router       /api/order [post]
func CreateOrder(c *fiber.Ctx) error {
	var createOrderDto CreateOrderDto
	if err := c.BodyParser(&createOrderDto); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	var user models.User
	if err := findUserById(int(createOrderDto.UserId), &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	var product models.Product
	if err := findProductById(int(createOrderDto.ProductId), &product); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	var saveOrder models.Order
	saveOrder.ProductRefer = int(product.ID)
	saveOrder.UserRefer = int(user.ID)
	database.Database.Db.Create(&saveOrder)
	responseOrder := CreateOrderResponse(saveOrder, user, product)
	return c.Status(200).JSON(responseOrder)
}

// GetOrders godoc
// @Summary      List orders
// @Description  get orders
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200  {array}   OrderDto
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/order [get]
func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order
	database.Database.Db.Find(&orders)
	responseOrders := []OrderDto{}
	for _, order := range orders {
		var user models.User
		var product models.Product
		database.Database.Db.Find(&user, "id = ?", order.UserRefer)
		database.Database.Db.Find(&product, "id = ?", order.ProductRefer)
		responseOrders = append(responseOrders, CreateOrderResponse(order, user, product))
	}
	return c.Status(200).JSON(responseOrders)
}
