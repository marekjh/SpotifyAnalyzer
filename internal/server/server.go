package server

import (
	"context"
	"log"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/marekjh/spotifyanalyzer/internal/auth"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Server struct {
	Authenticator *spotifyauth.Authenticator
	Engine        *gin.Engine
	EnvVars       *Env
	Logger        *zap.SugaredLogger
	TokenCache    *TokenCache
	SpotifyClient *spotify.Client
}

type Env struct {
	ClientID        string `env:"SPOTIFY_ID"`
	ClientSecret    string `env:"SPOTIFY_SECRET"`
	AuthRedirectURL string `env:"SPOTIFY_REDIRECT_URL" envDefault:"https://127.0.0.1/auth-response"`
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

	sc := spotify.New(http.DefaultClient)

	tokenCache := &TokenCache{
		Data: make(map[string]*oauth2.Token),
	}

	s := &Server{
		Authenticator: authenticator,
		Engine:        engine,
		EnvVars:       &envVars,
		Logger:        logger.Sugar(),
		TokenCache:    tokenCache,
		SpotifyClient: sc,
	}

	s.addRoutes()

	return s
}
