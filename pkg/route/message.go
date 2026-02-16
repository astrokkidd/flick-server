package route

import (
	"net/http"
	"time"

	"github.com/astrokkidd/flick/pkg/crypto"
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

type MessageResponse struct {
	MessageID int64  `json:"message_id"`
	SenderID  int64  `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func NewMessageHandler(queries *database.Queries, conn *pgx.Conn, tokenHandler *identity.TokenHandler) Message {
	return Message{queries, conn, tokenHandler}
}

func (message *Message) GetMessages(c echo.Context) error {
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	senderID := claims.ID()

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

	isInChat, err := qtx.IsUserInChat(ctx, database.IsUserInChatParams{
		UserID: senderID,
		ChatID: body.ChatID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user doesn't exist")
	}
	if !isInChat {
		return echo.NewHTTPError(http.StatusForbidden, "User not in chat")
	}

	messages, err := qtx.GetMessages(ctx, database.GetMessagesParams{
		ChatID: body.ChatID,
		Limit:  6,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "get messages failed")
	}

	var response []MessageResponse

	for _, m := range messages {
		plaintext, err := crypto.Decrypt(m.CypherText)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "decryption failed")
		}

		response = append(response, MessageResponse{
			MessageID: m.MessageID,
			SenderID:  m.SenderID,
			Content:   string(plaintext),
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	}

	return c.JSON(http.StatusOK, response)
}

func (message *Message) CreateMessage(c echo.Context) error {
	claims, err := identity.GetUserClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	senderID := claims.ID()

	var body struct {
		Content string `json:"content"`
		ChatID  int64  `param:"id"`
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

	isInChat, err := qtx.IsUserInChat(ctx, database.IsUserInChatParams{
		UserID: senderID,
		ChatID: body.ChatID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user doesn't exist")
	}
	if !isInChat {
		return echo.NewHTTPError(http.StatusForbidden, "User not in chat")
	}

	encrypted, err := crypto.Encrypt([]byte(body.Content))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "encryption failed")
	}

	createMessageParams := database.CreateMessageParams{
		ChatID:     body.ChatID,
		SenderID:   senderID,
		CypherText: encrypted,
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

	return c.JSON(http.StatusCreated, messageId)
}
