package route

import (
	"net/http"

	"github.com/astrokkidd/flick/pkg/database"
	"github.com/astrokkidd/flick/pkg/identity"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Chat struct {
	queries      *database.Queries
	conn         *pgx.Conn
	tokenHandler *identity.TokenHandler
}

func NewChatHandler(queries *database.Queries, conn *pgx.Conn, tokenHandler *identity.TokenHandler) Chat {
	return Chat{queries, conn, tokenHandler}
}

func (chat *Chat) CreateChat(c echo.Context) error {
	//-- Authenticated user --//
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	//-- Read participant_id from JSON --//
	var body struct {
		ParticipantID int64 `json:"participant_id"`
	}
	if err := c.Bind(&body); err != nil || body.ParticipantID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing or invalid participant_id")
	}
	if body.ParticipantID == uid {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot start chat with self")
	}

	//-- Begin transaction --//
	ctx := c.Request().Context()
	tx, err := chat.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := chat.queries.WithTx(tx)

	//-- Check if chat already exists between these two --//
	found, err := qtx.FindDirectChatBetween(ctx, database.FindDirectChatBetweenParams{
		UserID:   uid,
		UserID_2: body.ParticipantID,
	})
	if err == nil {
		// Chat already exists
		return c.JSON(http.StatusOK, echo.Map{"chat_id": found})
	}

	areFriends, err := qtx.AreUsersFriends(ctx, database.AreUsersFriendsParams{
		UserID:   uid,
		UserID_2: body.ParticipantID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user doesn't exist")
	}
	if !areFriends {
		return echo.NewHTTPError(http.StatusForbidden, "you must be friends to start a direct chat")
	}

	//-- Create new chat --//
	created, err := qtx.CreateEmptyChat(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create chat")
	}
	cid := created.ChatID

	//-- Add both participants --//
	for _, pid := range []int64{uid, body.ParticipantID} {
		err := qtx.AddParticipant(ctx, database.AddParticipantParams{
			ChatID: cid,
			UserID: pid,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not add participant")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.JSON(http.StatusCreated, echo.Map{"chat_id": cid})
}

func (chat *Chat) GetChats(c echo.Context) error {
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	//-- Begin tx --//
	ctx := c.Request().Context()
	tx, err := chat.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := chat.queries.WithTx(tx)

	listChatsWithParticipantParams := database.ListChatsWithParticipantParams{
		UserID: uid,
	}

	chats, err := qtx.ListChatsWithParticipant(ctx, listChatsWithParticipantParams)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "chats query failed")
	}

	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"chats": chats,
	})
}

func (chat *Chat) SetTypingStatus(c echo.Context) error {
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	var body struct {
		ChatID   int64 `param:"id"`
		IsTyping bool  `param:"status"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	//-- Begin tx --//
	ctx := c.Request().Context()
	tx, err := chat.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := chat.queries.WithTx(tx)

	isParticipant, err := qtx.IsUserInChat(ctx, database.IsUserInChatParams{ChatID: body.ChatID, UserID: uid})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not verify participant")
	}
	if !isParticipant {
		return echo.NewHTTPError(http.StatusForbidden, "not a participant in this chat")
	}

	_, err = qtx.SetTypingStatus(ctx, database.SetTypingStatusParams{IsTyping: body.IsTyping, ChatID: body.ChatID, UserID: uid})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not set typing status")
	}

	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.NoContent(http.StatusNoContent)
}

func (chat *Chat) SetLastReadMessage(c echo.Context) error {
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}
	uid := claims.ID()

	var body struct {
		ChatID    int64 `param:"chat_id"`
		MessageID int64 `json:"message_id"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	ctx := c.Request().Context()
	_, err = chat.queries.SetLastReadMessage(ctx, database.SetLastReadMessageParams{LastReadMessageID: &body.MessageID, ChatID: body.ChatID, UserID: uid})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not set last read message")
	}

	return c.NoContent(http.StatusNoContent)
}
