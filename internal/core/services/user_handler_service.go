package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang-rest/internal/core/domain"
	"golang-rest/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type UserHandlerService struct {
	userRepository ports.UserRepositoryInterface
}

func NewUserHandlerService(userRepository ports.UserRepositoryInterface) ports.UserHandlerInterface {
	return &UserHandlerService{userRepository: userRepository}
}

func (u UserHandlerService) RegisterUser(ctx *fiber.Ctx) error {
	var user domain.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body!"})
	}
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields!"})
	}

	// Check if the email already exists
	existingUser, err := u.userRepository.GetUserByEmail(user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error!"})
	}
	if existingUser != nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered!"})
	}

	// Proceed to create the user
	err = u.userRepository.CreateUser(&user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create user!"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully!",
	})
}

func (u UserHandlerService) LoginUser(ctx *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format!"})
	}

	// Find the user by email
	user, err := u.userRepository.GetUserLoginByEmail(input.Email)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Find not found user!"})
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password!"})
	}

	// JWT creation
	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"issuedAt": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate token!"})
	}

	return ctx.JSON(fiber.Map{
		"token": signedToken,
	})
}

func (u UserHandlerService) GetAllUsers(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")

	users, err := u.userRepository.GetAllUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get users!"})
	}

	return ctx.JSON(fiber.Map{
		"user_id": userID,
		"users":   users,
	})
}

func (u UserHandlerService) GetUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(user)
}

func (u UserHandlerService) UpdateUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var payload map[string]string
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updates := bson.M{}
	if name, ok := payload["name"]; ok {
		updates["name"] = name
	}
	if email, ok := payload["email"]; ok {
		updates["email"] = email
	}

	user, err := u.userRepository.UpdateUserByID(id, updates)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return ctx.JSON(user)
}

func (u UserHandlerService) DeleteUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := u.userRepository.DeleteUserByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return ctx.JSON(fiber.Map{"message": "User deleted successfully"})
}
