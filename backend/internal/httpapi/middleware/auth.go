package middleware

import (
	"strings"

	"github.com/critiq17/tripListik/internal/auth"
	"github.com/critiq17/tripListik/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ctxUserIDKey     = "userID"
	ctxTelegramIDKey = "telegramID"
)

func RequireAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header")
		}

		tokenStr := parts[1]
		claims, err := parseToken(cfg, tokenStr)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token subject")
		}

		c.Locals(ctxUserIDKey, userID)
		c.Locals(ctxTelegramIDKey, claims.TelegramID)
		return c.Next()
	}
}

// OptionalAuth resolves the user from a Bearer token when one is present,
// but lets the request through anonymously otherwise. Public endpoints need
// this so viewer-specific fields (viewer_is_member, private trip access)
// work for logged-in users.
func OptionalAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Next()
		}

		claims, err := parseToken(cfg, parts[1])
		if err != nil {
			return c.Next()
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return c.Next()
		}

		c.Locals(ctxUserIDKey, userID)
		c.Locals(ctxTelegramIDKey, claims.TelegramID)
		return c.Next()
	}
}

func parseToken(cfg *config.Config, tokenStr string) (*auth.Claims, error) {
	claims := &auth.Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	v := c.Locals(ctxUserIDKey)
	id, ok := v.(uuid.UUID)
	return id, ok
}

func RequireAuthQuery(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Query("token", "")
		if tokenStr == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		claims, err := parseToken(cfg, tokenStr)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token subject")
		}

		c.Locals(ctxUserIDKey, userID)
		c.Locals(ctxTelegramIDKey, claims.TelegramID)
		return c.Next()
	}
}
