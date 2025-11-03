package server

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// Config for cookies
const (
	BasePath = "/"
	Domain   = "127.0.0.1"
	Hour     = 3600
)

type TokenCache struct {
	Data  map[string]*oauth2.Token
	Mutex sync.RWMutex
}

func (s *Server) updateTokenCache(c *gin.Context, token *oauth2.Token) {
	s.TokenCache.Mutex.Lock()
	s.TokenCache.Data[c.ClientIP()] = token
	s.TokenCache.Mutex.Unlock()
}

func (s *Server) respondWithError(c *gin.Context, code int, err error) {
	s.Logger.Error(err)

	c.AbortWithStatusJSON(code, gin.H{"Error": err.Error()})
}

func setCookie(c *gin.Context, name, value string) {
	c.SetCookie(name, value, Hour, BasePath, Domain, true, false)
}
