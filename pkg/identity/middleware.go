package identity

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	userClaimsKey string = "userClaims"
)

func GetUserClaims(c echo.Context) (*UserClaims, error) {
	claims, ok := c.Get(userClaimsKey).(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("failed to resolve claims")
	}

	return claims, nil
}

func Authenticate(handler *TokenHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			header := r.Header.Get("Authorization")
			slog.Info("Authenticate called",
				"time", time.Now(),
				"path", r.URL.Path,
				"method", r.Method,
				"authHeader", header,
			)
			token := strings.TrimPrefix(header, "Bearer ")
			if token == header { // Bearer prefix not found
				slog.Warn("Authorization header missing or not Bearer format")
			}
			claims, err := handler.Verify(token)
			if err != nil {
				slog.Error("token verification failed", "error", err)
				return echo.ErrUnauthorized.WithInternal(err)
			}

			slog.Info("token verified", "subject", claims.Subject, "firstName", claims.FirstName)
			c.Set(userClaimsKey, claims)
			return next(c)
		}
	}
}
