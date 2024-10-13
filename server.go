package sfw

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/StevenACoffman/errgroup"
	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"

	"github.com/rs/zerolog/log"
)

type Server struct {
	apiGroup       hureg.APIGen
	adapter        *echo.Echo
	stopSignalChan chan os.Signal
}

type (
	APIGen = hureg.APIGen
)

// NewServer will create struct contains huma.API and *echo.Echo
func NewServer(title, version string, isDevelopmentMode bool) Server {
	openAPIPath := "/openapi"
	docsPath := "/docs"

	if !isDevelopmentMode {
		openAPIPath = ""
		docsPath = ""
	}

	config := huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   title,
				Version: version,
			},
			Components: &huma.Components{
				SecuritySchemes: map[string]*huma.SecurityScheme{
					"bearer": {
						Type:         "http",
						Scheme:       "bearer",
						BearerFormat: "token",
					},
				},
			},
		},
		OpenAPIPath:   openAPIPath,
		DocsPath:      docsPath,
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
		CreateHooks:   nil,
	}

	e := echo.New()
	e.HideBanner = true

	api := humaecho.New(e, config)
	apiGroup := hureg.NewAPIGen(api)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return Server{apiGroup: apiGroup, adapter: e, stopSignalChan: sigChan}
}

// API represents a Huma API wrapping a specific router.
func (s *Server) API() huma.API {
	return s.apiGroup.GetHumaAPI()
}

// Use adds middleware to the chain which is run after router.
func (s *Server) Use(middlewares ...echo.MiddlewareFunc) {
	s.adapter.Use(middlewares...)
}

// Start starts an HTTP server
func (s *Server) Start(port string) {
	if err := s.adapter.Start(port); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().AnErr("cannot start server", err)
	}
}

type CleanupJob struct {
	Name string
	Func func(context.Context) error
}

// Stop will stop the server is the signal channel is received. If not this will block forever
func (s *Server) Stop(ctx context.Context, jobs ...CleanupJob) {
	<-s.stopSignalChan

	log.Info().Msg("stop signal received")

	eg := errgroup.Group{}

	for _, job := range jobs {
		eg.Go(func() error {
			log.Info().Msgf("working on: %s", job.Name)
			return job.Func(ctx)
		})
	}

	err := eg.Wait()
	if err != nil {
		log.Err(err).Msg("cannot shutdown job")
	}

	err = s.adapter.Shutdown(ctx)
	if err != nil {
		log.Fatal().AnErr("cannot shutdown server", err)
	}
}

// Group create a new route group.
// It is recommended that you dont call APIGen.AddBasePath again, because that is already handled for you.
func (s *Server) Group(basePath string, groupFunc func(APIGen)) {
	group := s.apiGroup.AddBasePath(basePath)
	groupFunc(group)
}
