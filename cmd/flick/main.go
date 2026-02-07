package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/astrokkidd/flick/pkg/database"
	"github.com/astrokkidd/flick/pkg/identity"
	"github.com/astrokkidd/flick/pkg/route"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/joho/godotenv/autoload"
)

var cfg Config

func init() {
	cfg.Load()
}

func main() {
	//url := "postgres://postgres:admin@localhost:5432/flick"
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.PostgresUrl)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	queries := database.New(conn)

	tokenHandler := identity.NewTokenHandler(cfg.JwtSecret)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/v1")

	//-- AUTH --//
	authHandler := route.NewAuthHandler(queries, &tokenHandler)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	//-- USER --//
	userHandler := route.NewUserHandler(queries, conn, &tokenHandler)
	users := api.Group("/users", identity.Authenticate(&tokenHandler))
	users.PUT("/pfp", userHandler.UpdateProfilePicture)
	users.PUT("/pfp/delete", userHandler.RemoveProfilePicture)
	users.PUT("/display-name", userHandler.UpdateDisplayName)
	users.PUT("/password", userHandler.UpdatePassword)
	users.GET("/profile", userHandler.GetProfile)

	//-- FRIENDS --//
	requestHandler := route.NewRequestHandler(queries, conn, &tokenHandler)
	friends := api.Group("/friends", identity.Authenticate(&tokenHandler))
	friends.GET("", requestHandler.GetFriends)
	friends.GET("/requests/received", requestHandler.GetReceivedRequests)
	friends.GET("/requests/sent", requestHandler.GetSentRequests)
	friends.POST("/requests/send", requestHandler.SendRequest)
	friends.POST("/requests/:id/accept", requestHandler.AcceptRequest)
	friends.POST("/requests/:id/delete", requestHandler.DeleteRequest)

	//-- CHATS --//
	chatHandler := route.NewChatHandler(queries, conn, &tokenHandler)
	chat := api.Group("/chats", identity.Authenticate(&tokenHandler))
	chat.POST("", chatHandler.CreateChat)
	chat.GET("", chatHandler.GetChats)
	chat.POST("/:id/read", chatHandler.SetLastReadMessage)
	chat.POST("/:id/typing/:status", chatHandler.SetTypingStatus)

	//-- MESSAGES --//
	messageHandler := route.NewMessageHandler(queries, conn, &tokenHandler)
	chat.POST("/:id/messages", messageHandler.CreateMessage)
	chat.GET("/:id/messages", messageHandler.GetMessages)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	log.Println("Flick API running on http://localhost:8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down Flick API...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	log.Println("Server stopped cleanly")
}
