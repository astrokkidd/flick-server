package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astrokkidd/flick/pkg/database"
	"github.com/astrokkidd/flick/pkg/identity"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Request struct {
	queries      *database.Queries
	conn         *pgx.Conn
	tokenHandler *identity.TokenHandler
}

type FriendResponse struct {
	UserID       int64  `json:"user_id"`
	PfpURL       string `json:"pfp_url"`
	DisplayName  string `json:"display_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FriendshipTs string `json:"friendship_ts"` // string for React
}

func NewRequestHandler(queries *database.Queries, conn *pgx.Conn, tokenHandler *identity.TokenHandler) Request {
	return Request{queries, conn, tokenHandler}
}

func (r *Request) SendRequest(c echo.Context) error {
	//-- Verify user --//
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	//-- Get display name from payload --//
	var payload struct {
		DisplayName string `json:"display_name"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json").SetInternal(err)
	}
	dn := strings.TrimSpace(payload.DisplayName)
	if dn == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "display_name is required")
	}

	//-- Begin tx --//
	ctx := c.Request().Context()
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	//-- Get the user id from payload display name --//
	receiver, err := qtx.FindUserByDisplayName(ctx, database.FindUserByDisplayNameParams{DisplayName: dn})
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found").SetInternal(err)
	}
	if receiver.UserID == uid {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot send request to yourself")
	}

	// -- See if friend request already exists --//
	alreadyExists, err := qtx.DoesFriendRequestExist(ctx, database.DoesFriendRequestExistParams{
		SenderID:   uid,
		ReceiverID: receiver.UserID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "friend request existence check failed").SetInternal(err)
	}
	if alreadyExists {
		return echo.NewHTTPError(http.StatusConflict, "friend request already exists")
	}

	//-- See if users are already friends --//
	alreadyFriends, err := qtx.AreUsersFriends(ctx, database.AreUsersFriendsParams{UserID: uid, UserID_2: receiver.UserID})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "friend existence check failed").SetInternal(err)
	}
	if alreadyFriends {
		return echo.NewHTTPError(http.StatusConflict, "users already friends")
	}

	//-- Create friend request --//
	fr, err := qtx.CreateFriendRequest(ctx, database.CreateFriendRequestParams{
		SenderID:   uid,
		ReceiverID: receiver.UserID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "insert failed").SetInternal(err)
	}

	//-- Commit all queries --//
	if err := tx.Commit(ctx); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, fr)
}

func (request *Request) GetReceivedRequests(c echo.Context) error {
	//-- Verify user --//
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	//-- List all friend requests to user --//
	result, err := request.queries.ListReceivedFriendRequestsWithUser(c.Request().Context(), database.ListReceivedFriendRequestsWithUserParams{ReceiverID: uid})
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}
	fmt.Println("DEBUG result:", result)

	return c.JSON(http.StatusOK, map[string]any{
		"requests": result,
	})
}

func (request *Request) GetSentRequests(c echo.Context) error {
	//-- Verify user --//
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	//-- List all friend requests to user --//
	result, err := request.queries.ListSentFriendRequestsWithUser(c.Request().Context(), database.ListSentFriendRequestsWithUserParams{SenderID: uid})
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}
	fmt.Println("DEBUG result:", result)

	return c.JSON(http.StatusOK, map[string]any{
		"requests": result,
	})
}

func (request *Request) GetFriends(c echo.Context) error {
	//-- Verify user --//
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	//-- List all friends of user --//
	result, err := request.queries.ListAllFriends(c.Request().Context(), database.ListAllFriendsParams{UserID: uid})
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	friends := make([]FriendResponse, len(result))
	for i, r := range result {
		friends[i] = FriendResponse{
			UserID:       r.UserID,
			PfpURL:       *r.PfpUrl,
			DisplayName:  r.DisplayName,
			FirstName:    r.FirstName,
			LastName:     r.LastName,
			FriendshipTs: r.FriendshipTs.Format(time.RFC3339), // ISO string
		}
	}
	b, _ := json.Marshal(friends)
	fmt.Println("DEBUG JSON:", string(b))

	return c.JSON(http.StatusOK, friends)
}

func (request *Request) AcceptRequest(c echo.Context) error {
	//-- Verify user --//
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	//-- Get request id from params --//
	rid_str := c.Param("id")
	rid, err := strconv.ParseInt(rid_str, 10, 64)
	if err != nil || rid <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request id")
	}

	//-- Begin tx --//
	ctx := c.Request().Context()
	tx, err := request.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := request.queries.WithTx(tx)

	//-- Get friend request sender --//
	fid, err := qtx.GetUserByRequestID(ctx, database.GetUserByRequestIDParams{RequestID: rid})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "query failed")
	}

	//-- Create friendship --//
	err = qtx.CreateFriendship(ctx, database.CreateFriendshipParams{UserID: uid, FriendID: fid})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "friendship insert failed")
	}

	//-- Delete deprecated friend request --//
	res, err := qtx.DeleteFriendRequest(ctx, database.DeleteFriendRequestParams{SenderID: fid, RequestID: rid})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "request delete failed")
	}

	//-- Commit queries --//
	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.JSON(http.StatusOK, res)
}

func (request *Request) DeleteRequest(c echo.Context) error {
	//-- Verify user --//
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}

	_, err = strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user id in token")
	}

	//-- Get request id from params --//
	rid_str := c.Param("id")
	rid, err := strconv.ParseInt(rid_str, 10, 64)
	if err != nil || rid <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request id")
	}

	//-- Begin tx --//
	ctx := c.Request().Context()
	tx, err := request.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := request.queries.WithTx(tx)

	//-- Get friend request sender --//
	fid, err := qtx.GetUserByRequestID(ctx, database.GetUserByRequestIDParams{
		RequestID: rid,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "query failed")
	}

	//-- Delete deprecated friend request --//
	res, err := qtx.DeleteFriendRequest(ctx, database.DeleteFriendRequestParams{
		SenderID:  fid,
		RequestID: rid,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "request delete failed")
	}

	//-- Commit queries --//
	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.JSON(http.StatusOK, res)
}
