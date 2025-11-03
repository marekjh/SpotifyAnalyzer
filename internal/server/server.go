package server

import (
	"context"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/marekjh/spotifyanalyzer/internal/auth"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
)

type Server struct {
	Authenticator *spotifyauth.Authenticator
	Engine        *gin.Engine
	EnvVars       *Env
	Logger        *zap.SugaredLogger
	TokenCache    *TokenCache
}

type Env struct {
	ClientID        string `env:"SPOTIFY_ID"`
	ClientSecret    string `env:"SPOTIFY_SECRET"`
	AuthRedirectURL string `env:"SPOTIFY_REDIRECT_URL" envDefault:"http://127.0.0.1:8080/callback"`
}

func NewServer(ctx context.Context) *Server {
	engine := gin.New()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err.Error())
	}

	var envVars Env
	err = env.Parse(&envVars)
	if err != nil {
		log.Fatalf("failed to parse environment variables: %s", err.Error())
	}

	authenticator := spotifyauth.New(spotifyauth.WithRedirectURL(envVars.AuthRedirectURL), spotifyauth.WithScopes(auth.Scopes...))

	clientCache := &TokenCache{
		Data: make(map[string]Subcache),
	}

	s := &Server{
		Authenticator: authenticator,
		Engine:        engine,
		EnvVars:       &envVars,
		Logger:        logger.Sugar(),
		TokenCache:    clientCache,
	}

	s.addRoutes()

	return s
}
