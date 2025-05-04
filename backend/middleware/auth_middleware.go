package middleware

import (
  "strings"

  "backend/handlers"

  "github.com/gofiber/fiber/v3"
  "github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c fiber.Ctx) error {
  authHeader := c.Get("Authorization")
  if authHeader == "" {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
  }

  parts := strings.Split(authHeader, " ")
  if len(parts) != 2 || parts[0] != "Bearer" {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
  }

  tokenString := parts[1]

  token, err := jwt.ParseWithClaims(tokenString, &handlers.Claims{}, func(token *jwt.Token) (interface{}, error) {
    return handlers.JwtSecret, nil
  })

  if err != nil {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token", "details": err.Error()})
  }

  claims, ok := token.Claims.(*handlers.Claims)
  if !ok || !token.Valid {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
  }

  c.Locals("userID", claims.ID)
  c.Locals("userRole", claims.Role)

  return c.Next()
}

func FacultyOnly(c fiber.Ctx) error {
  role, ok := c.Locals("userRole").(string)
  if !ok || role != "faculty" {
    return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Faculty access required"})
  }
  return c.Next()
}

func StudentOnly(c fiber.Ctx) error {
  role, ok := c.Locals("userRole").(string)
  if !ok || role != "student" {
    return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Student access required"})
  }
  return c.Next()
}

