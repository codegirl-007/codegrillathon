package main

import (
	"codegrillathon/handlers"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/twitch"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("HELP ME LOAD SOME ENV FILE")
	}
	tmpl, err := template.ParseGlob("templates/*.html")

	if err != nil {
		panic(fmt.Errorf("failed to parse templates: %w\n", err))
	}

	goth.UseProviders(
		twitch.New(os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_SECRET"), "https://localhost:8080/auth/twitch/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "https://localhost:8080/auth/github/callback"),
	)

	h := handlers.Handler{Template: *tmpl}
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// auth routes
	mux.HandleFunc("/auth/twitch", h.Auth)
	mux.HandleFunc("/auth/twitch/callback", h.Callback)
	mux.HandleFunc("/logout/", h.Logout)

	// pages route
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/welcome", h.Welcome)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("Server started...")
		if err := server.ListenAndServeTLS("server.crt", "server.key"); err != http.ErrServerClosed {
			fmt.Printf("Serve error: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server force to shutdown: %v\n", err)
	} else {
		fmt.Println("Server exit gracefully")
	}
}
