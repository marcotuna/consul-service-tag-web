package app

import (
	"consul-service-tag-web/model"
	"consul-service-tag-web/service"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// App ...
type App struct {
	ConsulClient *service.ConsulClient
	HTTPServer   *http.Server
	Config       *model.Config
}

// NewApp initializes application
func NewApp(config *model.Config) (*App, error) {

	app := App{}

	// Initialize consul client connection
	log.Info("Initializing Consul Client...")

	consulClient, err := service.NewConsulClient()

	if err != nil {
		return nil, err
	}

	app.ConsulClient = consulClient

	// Config
	if app.Config == nil {
		app.Config = config
	}

	// Gin HTTP Server
	log.Info("Initializing HTTP Server...")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLGlob("template/*.tmpl")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.tmpl", gin.H{
			"title": "Main website",
		})
	})

	apiV1 := r.Group("/v1")
	{
		apiV1.GET("/services", app.getServices)
	}
	r.Use(gin.Recovery())

	srv := &http.Server{
		Addr:    config.HTTP.Address,
		Handler: r,
	}

	app.HTTPServer = srv

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	return &app, nil
}

// Shutdown closes running services
func (a *App) Shutdown() error {
	log.Info("Shutting HTTP Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Info("Exitting...")

	return nil
}

func (a *App) getServices(c *gin.Context) {
	// Parse Service Tags
	configServiceTags := strings.Split(a.Config.Tag.Filter, ",")

	serviceTags, err := a.ConsulClient.GetConsulServices(configServiceTags)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status":  true,
		"message": serviceTags,
	})
}
