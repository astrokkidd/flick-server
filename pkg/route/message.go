package route

import (
	"net/http"

	"github.com/astrokkidd/flick/pkg/database"
	"github.com/astrokkidd/flick/pkg/identity"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Message struct {
	queries      *database.Queries
	conn         *pgx.Conn
	tokenHandler *identity.TokenHandler
}

func NewMessageHandler(queries *database.Queries, conn *pgx.Conn, tokenHandler *identity.TokenHandler) Message {
	return Message{queries, conn, tokenHandler}
}

func (message *Message) GetMessages(c echo.Context) error {
	var body struct {
		ChatID int64 `param:"id"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input").SetInternal(err)
	}

	//-- Begin transaction --//
	ctx := c.Request().Context()
	tx, err := message.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := message.queries.WithTx(tx)

	messages, err := qtx.GetMessages(ctx, database.GetMessagesParams{
		ChatID: body.ChatID,
		Limit:  6,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "get messages failed")
	}

	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"messages": messages,
	})
}

func (message *Message) CreateMessage(c echo.Context) error {
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	senderID := claims.ID()

	var body struct {
		Content    string `json:"content"`
		CypherText string ``
		Nonce      string ``
		ChatID     int64  `param:"id"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json").SetInternal(err)
	}

	//-- Begin tx --//
	ctx := c.Request().Context()
	tx, err := message.conn.Begin(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := message.queries.WithTx(tx)

	createMessageParams := database.CreateMessageParams{
		ChatID:     body.ChatID,
		SenderID:   senderID,
		CypherText: []byte(body.Content),
		Nonce:      []byte(body.Nonce),
	}

	messageId, err := qtx.CreateMessage(ctx, createMessageParams)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "message creation failed").SetInternal(err)
	}

	err = qtx.UpdateChatLastMessage(ctx, database.UpdateChatLastMessageParams{ChatID: body.ChatID, LastMessageID: &messageId})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "update last message failed").SetInternal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.JSON(http.StatusCreated, nil)
}
