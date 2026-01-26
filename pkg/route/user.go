package route

import (
	"fmt"
	"net/http"

	"github.com/astrokkidd/flick/pkg/database"
	"github.com/astrokkidd/flick/pkg/identity"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type User struct {
	queries      *database.Queries
	conn         *pgx.Conn
	tokenHandler *identity.TokenHandler
}

func NewUserHandler(queries *database.Queries, conn *pgx.Conn, tokenHandler *identity.TokenHandler) User {
	return User{queries, conn, tokenHandler}
}

func (user *User) UpdateProfilePicture(c echo.Context) error {
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	var body struct {
		ProfilePictureURL string `json:"profile_picture_url"`
	}
	if err := c.Bind(&body); err != nil || body.ProfilePictureURL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing or invalid profile_picture_url")
	}

	err = user.queries.UpdateUserPfp(c.Request().Context(), database.UpdateUserPfpParams{
		UserID: uid,
		PfpUrl: &body.ProfilePictureURL,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update profile picture")
	}

	return echo.NewHTTPError(http.StatusOK, nil)
}

func (user *User) RemoveProfilePicture(c echo.Context) error {
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	//-- Begin tx --//
	ctx := c.Request().Context()
	tx, err := user.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := user.queries.WithTx(tx)

	u, err := qtx.FindUserByID(ctx, database.FindUserByIDParams{UserID: uid})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch user")
	}

	defaultPfp := fmt.Sprintf("https://api.dicebear.com/7.x/notionists-neutral/png?seed=%s", u.DisplayName)

	err = qtx.UpdateUserPfp(c.Request().Context(), database.UpdateUserPfpParams{
		UserID: uid,
		PfpUrl: &defaultPfp,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update profile picture")
	}

	//-- Commit queries --//
	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.JSON(http.StatusOK, nil)
}

func (user *User) UpdateDisplayName(c echo.Context) error {
	return echo.NewHTTPError(501, "not implemented")
}

func (user *User) UpdatePassword(c echo.Context) error {
	return echo.NewHTTPError(501, "not implemented")
}

func (user *User) GetProfile(c echo.Context) error {
	return echo.NewHTTPError(501, "not implemented")
}
