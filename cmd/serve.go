package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amirhossein-karimi/todo/internal/cache"
	"github.com/amirhossein-karimi/todo/internal/config"
	"github.com/amirhossein-karimi/todo/internal/db"
	"github.com/amirhossein-karimi/todo/internal/handler"
	"github.com/amirhossein-karimi/todo/internal/logger"
	"github.com/amirhossein-karimi/todo/internal/metrics"
	"github.com/amirhossein-karimi/todo/internal/middleware"
	"github.com/amirhossein-karimi/todo/internal/repository"
	"github.com/amirhossein-karimi/todo/internal/service"
	"github.com/amirhossein-karimi/todo/internal/validation"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",

	RunE: func(cmd *cobra.Command, args []string) error {

		log, err := logger.New()
		if err != nil {
			panic(err)
		}

		config, err := config.Load("config.yaml")
		if err != nil {
			log.Fatal("cannot load config", zap.Error(err))
		}

		gormDB, err := db.NewConnection(config)

		if err != nil {
			log.Fatal("cannot connect to database", zap.Error(err))
		}

		r := gin.Default()

		pprof.Register(r)

		err = validation.Register()
		if err != nil {
			log.Fatal("cannot register validation", zap.Error(err))
		}

		group := r.Group("/api/v1/todo")
		cache := cache.NewRedisCache(config)
		todoRepository := repository.NewTodoRepository(gormDB)

		todoSvc := service.NewTodoService(todoRepository, cache)
		h := handler.NewHandler(todoSvc, log)
		metrics.Register()

		group.Use(middleware.Prometheus())

		r.GET("/metrics", gin.WrapH(promhttp.Handler()))

		group.POST("/create", h.Create)
		group.GET("/list", h.List)
		group.DELETE("/delete/:uuid", h.Delete)
		group.GET("/:uuid", h.GetInfo)
		group.PUT("/update/:uuid", h.Update)

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		server := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}

		go func() {
			fmt.Println("Server started on :8080")

			if err := server.ListenAndServe(); err != nil &&
				!errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("server error: %v\n", err)
			}
		}()

		sigChan := make(chan os.Signal, 1)

		signal.Notify(
			sigChan,
			syscall.SIGINT,
			syscall.SIGTERM,
		)

		<-sigChan

		fmt.Println("Shutdown signal received")

		// Give active requests 10 seconds to finish
		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}

		fmt.Println("Server stopped gracefully")

		return nil
	},
}

func init() {

	rootCmd.AddCommand(serveCmd)
}
