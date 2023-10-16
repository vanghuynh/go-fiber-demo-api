package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/vanghuynh/fiber-api/database"
	"github.com/vanghuynh/fiber-api/models"
)

type UserDto struct {
	// this is not the model user, just serializer
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUserDto struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func CreateResponseUser(user models.User) UserDto {
	return UserDto{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

// CreateUser godoc
// @Summary      Create user
// @Description  create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        data body models.User true "The input user struct"
// @Success      200  {object}  UserDto
// @Router       /api/user [post]
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

// GetUsers godoc
// @Summary      List users
// @Description  get users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   UserDto
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/user [get]
func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	database.Database.Db.Find(&users)
	responseUsers := []UserDto{}
	for _, user := range users {
		responseUsers = append(responseUsers, CreateResponseUser(user))
	}
	return c.Status(200).JSON(responseUsers)
}

// GetUser godoc
// @Summary      Get user
// @Description  get user by id
// @Tags         users
// @Accept       json
// @Param        id path string true "User ID"
// @Produce      json
// @Success      200  {object}   UserDto
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/user/{id} [get]
func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User
	if err != nil {
		return c.Status(400).JSON("Invalid ID")
	}
	if err := findUserById(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  update user by id
// @Tags         users
// @Accept       json
// @Param        id path string true "User ID"
// @Param        data body UpdateUserDto true "The input user update"
// @Produce      json
// @Success      200  {object}   UserDto
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/user/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User
	if err != nil {
		return c.Status(400).JSON("Invalid ID")
	}
	if err := findUserById(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	var updateUserDto UpdateUserDto
	if err := c.BodyParser(&updateUserDto); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	user.FirstName = updateUserDto.FirstName
	user.LastName = updateUserDto.LastName
	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  delete user by id
// @Tags         users
// @Accept       json
// @Param        id path string true "User ID"
// @Produce      json
// @Success      200  {object}   UserDto
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/user/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User
	if err != nil {
		return c.Status(400).JSON("Invalid user id")
	}
	if err := findUserById(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	// check if error when delete user
	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON("Error delete user")
	}
	return c.Status(200).JSON("Deleted user")
}

func findUserById(id int, user *models.User) error {
	database.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("User not found")
	}
	return nil
}
