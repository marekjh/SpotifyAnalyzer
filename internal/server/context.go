package server

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

// Config for cookies
const (
	BasePath = "/"
	Domain   = "127.0.0.1"
	Hour     = 3600
)

// Keys for objects living in gin.Context.Keys
const (
	SpotifyClient         = "SpotifyClient"
	MyRecentTracksRequest = "MyRecentTracksRequest"
)

type TokenCache struct {
	Data  map[string]Subcache
	Mutex sync.RWMutex
}

type Subcache struct {
	Token         *oauth2.Token
	SpotifyClient *spotify.Client
}

func (s *Server) updateTokenCache(c *gin.Context, token *oauth2.Token) {
	s.TokenCache.Mutex.Lock()
	s.TokenCache.Data[c.ClientIP()] = Subcache{
		Token:         token,
		SpotifyClient: spotify.New(s.Authenticator.Client(c.Request.Context(), token)),
	}
	s.TokenCache.Mutex.Unlock()
}

func (s *Server) respondWithError(c *gin.Context, code int, err error) {
	s.Logger.Error(err)

	c.AbortWithStatusJSON(code, gin.H{"Error": err.Error()})
}

func setCookie(c *gin.Context, name, value string) {
	c.SetCookie(name, value, Hour, BasePath, Domain, true, false)
}

func retrieveSpotifyClient(c *gin.Context) (*spotify.Client, error) {
	value := c.Value(SpotifyClient)

	client, ok := value.(*spotify.Client)
	if !ok {
		return nil, errors.New("failed to type assert spotify client")
	}

	return client, nil
}

func setSpotifyClient(c *gin.Context, client *spotify.Client) {
	c.Set(SpotifyClient, client)
}

func retrieveMyRecentTracksRequest(c *gin.Context) (*myRecentTracksRequest, error) {
	value := c.Value(MyRecentTracksRequest)

	req, ok := value.(*myRecentTracksRequest)
	if !ok {
		return nil, errors.New("failed to type assert my recent tracks request")
	}

	return req, nil
}

func setMyRecentTracksRequest(c *gin.Context, req *myRecentTracksRequest) {
	c.Set(MyRecentTracksRequest, req)
}
