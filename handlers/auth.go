package handlers

import (
	"os"
	"time"

	"github.com/cmerin0/SimpleCarsApp/db"
	"github.com/cmerin0/SimpleCarsApp/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateToken creates a new JWT token for the given email address
// It sets the token expiration to 24 hours from creation
// Returns the signed token string and any error that occurred
func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(), // change to 24 hours, right now its only 2 hours
	})

	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT token string
// It takes a token string as input and returns a status message and any error
// The function verifies the token signature using the JwtKey
// Returns "Valid Token" if the token is valid, "Invalid Token" if not valid,
// or an empty string and error if parsing fails
func VerifyToken(c *fiber.Ctx) error {

	tokenString := c.Cookies("token") // Getting token from cookie

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Missing authorization token"})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	// Store the email in the context for later use
	c.Locals("email", email)

	return c.Next()
}

// Listing Users Protected
func GetUsers(c *fiber.Ctx) error {

	users := []models.User{} // Creating a list of users bucket
	db.DB.Db.Find(&users)    //Getting all users
	return c.Status(fiber.StatusOK).JSON(users)
}

func Register(c *fiber.Ctx) error {

	user := new(models.User) // Creating model of user

	// Getting the parameters sent in Body, saving them in user
	// Returning an error message if body is wrong
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Generating a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}

	user.Password = string(hashedPassword) // Saving hashed password in user model

	if result := db.DB.Db.Create(user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created", "User": user})
}

func Login(c *fiber.Ctx) error {

	loginUser := new(models.User) // Creating model of user

	// Getting the parameters sent in Body, saving them in user
	if err := c.BodyParser(&loginUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	dbUser := new(models.User) // Getting user from database

	// Getting user from database and Checking if user exists and theres no error
	if result := db.DB.Db.Where("email = ?", loginUser.Email).First(&dbUser); result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Comparing passwords
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginUser.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Generating token
	token, err := GenerateToken(loginUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate token"})
	}

	// Creating a cookie to set up the jwt authorization
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 2),  // Cookie expires in 24 hours
		HTTPOnly: true,                           // Important for security (set false for https only)
		Secure:   false,                          // Set true in production
		SameSite: fiber.CookieSameSiteStrictMode, // Important for security.
	}

	c.Cookie(&cookie) // Setting up cookie in fiber context

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful", "token": token})
}

func Logout(c *fiber.Ctx) error {

	// Resetting the cookie
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteStrictMode,
	}

	c.Cookie(&cookie) // Setting up cookie in fiber context

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful"})
}
