package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func New() *Server {
	s := &Server{
		router: gin.Default(),
	}

	healthzRouter := s.router.Group("/healthz/")
	healthzRouter.GET("/liveness", s.liveness)
	healthzRouter.GET("/readiness", s.readiness)

	s.router.GET("/", s.ShowAddItemForm)
	s.router.POST("/add-item", s.AddItem)
	s.router.GET("/remove-cart-item", s.DeleteCartItem)

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
