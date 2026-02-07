package route

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/astrokkidd/flick/pkg/database"
	"github.com/astrokkidd/flick/pkg/identity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	queries      *database.Queries
	tokenHandler *identity.TokenHandler
}

func NewAuthHandler(queries *database.Queries, tokenHandler *identity.TokenHandler) Auth {
	return Auth{queries, tokenHandler}
}

func (auth *Auth) Login(c echo.Context) error {
	// Accept both form and JSON
	var form struct {
		DisplayName string `form:"display_name" json:"display_name"`
		Password    string `form:"password"     json:"password"`
	}

	if err := c.Bind(&form); err != nil {
		c.Logger().Warnf("login bind error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	form.DisplayName = strings.TrimSpace(form.DisplayName)
	if form.DisplayName == "" || form.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "display_name and password are required")
	}

	// Add a small timeout so a slow DB doesn’t hang the request forever
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	// Look up user by display name
	user, err := auth.queries.FindUserByDisplayName(ctx, database.FindUserByDisplayNameParams{
		DisplayName: form.DisplayName,
	})
	if err != nil {
		// sqlc usually wraps sql.ErrNoRows (or pgx.ErrNoRows); handle "not found" as invalid creds
		if errors.Is(err, sql.ErrNoRows) {
			// Do not reveal whether name or password was wrong
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		c.Logger().Errorf("login db error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// Constant-time password verification
	validated, err := identity.Password(form.Password).ValidatePassword(user.PasswordHash)
	if !validated {
		// Same generic message to avoid user enumeration
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	token, err := auth.tokenHandler.Sign(identity.UserClaims{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		ProfileImageURL: derefString(user.PfpUrl),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  fmt.Sprint(user.UserID),
			Issuer:   "api.getflick.chat",
			Audience: jwt.ClaimStrings{"api.getflick.chat"},
		},
	})
	if err != nil {
		c.Logger().Errorf("token sign error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"access_token": token,
		"user_id":      user.UserID,
		"display_name": form.DisplayName,
		"pfp_url":      user.PfpUrl,
	})
}

func (auth *Auth) Register(c echo.Context) error {
	var form struct {
		DisplayName string `json:"display_name"`
		Password    string `json:"password"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		UserKey     string `json:"user_key"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if strings.TrimSpace(form.DisplayName) == "" ||
		strings.TrimSpace(form.Password) == "" ||
		strings.TrimSpace(form.FirstName) == "" ||
		strings.TrimSpace(form.LastName) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "all required fields must be provided")
	}

	defaultPfp := fmt.Sprintf("https://api.dicebear.com/7.x/notionists-neutral/png?seed=%s", url.QueryEscape(form.DisplayName))

	password := identity.Password(form.Password)
	hash, err := password.GenerateHash()
	if err != nil {
		panic(err)
	}

	user, err := auth.queries.CreateUser(c.Request().Context(), database.CreateUserParams{
		DisplayName:  form.DisplayName,
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		PasswordHash: hash,
		UserKey:      []byte(form.UserKey),
		PfpUrl:       &defaultPfp,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return echo.NewHTTPError(http.StatusConflict, "username already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create user").SetInternal(err)
	}

	token, err := auth.tokenHandler.Sign(identity.UserClaims{
		FirstName:       form.FirstName,
		LastName:        form.LastName,
		ProfileImageURL: defaultPfp,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  strconv.FormatInt(user.UserID, 10),
			Issuer:   "api.getflick.chat",
			Audience: jwt.ClaimStrings{"api.getflick.chat"},
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not generate token").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"access_token": token,
		"user_id":      user.UserID,
		"display_name": user.DisplayName,
		"pfp_url":      user.PfpUrl,
	})
}

func (auth *Auth) Recover(c echo.Context) error {
	// Perform TOTP recovery flow
	return nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
