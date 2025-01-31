package http

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammadne/ice-global/internal/services"
)

type Server struct {
	router *gin.Engine

	// services
	cartsService services.Carts
	itemsService services.Items
	usersService services.Users
}

//go:embed templates/*
var templates embed.FS

func New(cartsService services.Carts, itemsService services.Items, usersService services.Users) *Server {
	s := &Server{
		router:       gin.Default(),
		cartsService: cartsService,
		itemsService: itemsService,
		usersService: usersService,
	}

	// Parse the embedded templates
	tmplates, err := template.ParseFS(templates, "templates/*")
	if err != nil {
		slog.Error("error parsing templates", "Err", err)
		os.Exit(1)
	}
	// Set the embedded template to Gin's renderer
	s.router.SetHTMLTemplate(tmplates)

	healthzRouter := s.router.Group("/healthz/")
	healthzRouter.GET("/liveness", s.liveness)
	healthzRouter.GET("/readiness", s.readiness)

	s.router.GET("/", s.OptionalCookie, s.showAddItemForm)
	// s.router.POST("/add-item", s.RequiredCookie, s.addItem)
	// s.router.GET("/remove-cart-item", s.RequiredCookie, s.deleteCartItem)

	return s
}

func (s *Server) Serve(ctx context.Context, wg *sync.WaitGroup, port int) {
	defer wg.Done()
	var server *http.Server

	go func() {
		server = &http.Server{
			Handler:      s.router,
			Addr:         fmt.Sprintf("0.0.0.0:%d", port),
			WriteTimeout: 3 * time.Second,
			ReadTimeout:  3 * time.Second,
		}

		slog.Info("starting server", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil {
			slog.Error("error resolving server", slog.String("address", server.Addr), "Err", err)
			os.Exit(1)
		}

	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("error shutdown http server")
	}

	slog.Warn("gracefully shutdown the https servers")
}
