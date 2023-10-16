package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/vanghuynh/fiber-api/database"
	"github.com/vanghuynh/fiber-api/models"
)

type ProductDto struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateProductResponse(product models.Product) ProductDto {
	return ProductDto{
		ID:           product.ID,
		Name:         product.Name,
		SerialNumber: product.SerialNumber,
	}
}

// CreateProduct godoc
// @Summary      Create product
// @Description  create product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        data body models.Product true "The input product struct"
// @Success      200  {object}  ProductDto
// @Router       /api/product [post]
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.Database.Db.Create(&product)
	responseProduct := CreateProductResponse(product)
	return c.Status(200).JSON(responseProduct)
}

// GetProducts godoc
// @Summary      List products
// @Description  get products
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}   ProductDto
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/product [get]
func GetProducts(c *fiber.Ctx) error {
	products := []models.Product{}
	database.Database.Db.Find(&products)
	responseProducts := []ProductDto{}
	for _, product := range products {
		responseProduct := CreateProductResponse(product)
		responseProducts = append(responseProducts, responseProduct)
	}
	return c.Status(200).JSON(responseProducts)
}

func findProductById(id int, product *models.Product) error {
	database.Database.Db.Find(&product, "id = ?", id)
	if product.ID == 0 {
		return errors.New("Product not found")
	}
	return nil
}
