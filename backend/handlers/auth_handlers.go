package handlers

import (
  "context"
  "time"

  "backend/common"
  "backend/database"
  "backend/models"

  "github.com/gofiber/fiber/v3"
  "github.com/golang-jwt/jwt/v5"
  "github.com/jackc/pgx/v5"
)

var JwtSecret = []byte(common.MustGetEnv("JWT_SECRET"))

type Claims struct {
  ID   int    `json:"id"`
  Role string `json:"role"`
  jwt.RegisteredClaims
}

func generateJWT(id int, role string) (string, error) {
  claims := Claims{
    ID:   id,
    Role: role,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
      IssuedAt:  jwt.NewNumericDate(time.Now()),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(JwtSecret)
}

func Login(c fiber.Ctx) error {
  loginReq := new(models.LoginRequest)

  if err := c.Bind().Body(loginReq); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
  }

  if loginReq.ID == 0 || loginReq.Password == "" || (loginReq.Role != "student" && loginReq.Role != "faculty") {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID, password, and valid role are required"})
  }

  var storedPassword string
  var dateOfBirth *time.Time
  var query string

  if loginReq.Role == "student" {
    query = `SELECT password, date_of_birth FROM students WHERE id = $1`
  } else {
    query = `SELECT password, date_of_birth FROM faculty WHERE id = $1`
  }

  err := database.DB.QueryRow(context.Background(), query, loginReq.ID).Scan(&storedPassword, &dateOfBirth)

  if err == pgx.ErrNoRows {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
  } else if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error during login", "details": err.Error()})
  }

  if storedPassword == "" {
    if dateOfBirth == nil {
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
    }
    dobPassword := dateOfBirth.Format("2006-01-02")
    if loginReq.Password != dobPassword {
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
    }
  } else {
    if loginReq.Password != storedPassword {
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
    }
  }

  token, err := generateJWT(loginReq.ID, loginReq.Role)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token", "details": err.Error()})
  }

  return c.JSON(models.AuthResponse{Token: token, Role: loginReq.Role, ID: loginReq.ID})
}

