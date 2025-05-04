package handlers

import (
  "context"
  "errors"
  "time"

  "backend/common"
  "backend/database"
  "backend/models"

  "github.com/gofiber/fiber/v3"
  "github.com/golang-jwt/jwt/v5"
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

  var authenticatedUserID int
  var authenticatedUserRole string

  query := `SELECT user_id, user_role FROM authenticate_user($1, $2, $3)`
  err := database.DB.QueryRow(context.Background(), query,
    loginReq.ID,
    loginReq.Password,
    loginReq.Role,
  ).Scan(&authenticatedUserID, &authenticatedUserRole)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  token, err := generateJWT(authenticatedUserID, authenticatedUserRole)
  if err != nil {
    return sendInternalServerError(c, errors.New("Failed to generate token"))
  }

  return c.JSON(models.AuthResponse{Token: token, Role: authenticatedUserRole, ID: authenticatedUserID})
}

