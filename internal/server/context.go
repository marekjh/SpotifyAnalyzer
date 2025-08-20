package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// Keys for context
const (
	AuthKey = "auth"
)

var Error = "Something went wrong"

type TokenCache struct {
	Data  map[string]*oauth2.Token
	Mutex *sync.RWMutex
}

func (s *Server) updateTokenCache(c *gin.Context, token *oauth2.Token) {
	s.TokenCache.Mutex.RLock()
	cache := s.TokenCache.Data
	s.TokenCache.Mutex.RUnlock()

	cache[c.ClientIP()] = token

	s.TokenCache.Mutex.Lock()
	s.TokenCache.Data = cache
	s.TokenCache.Mutex.Unlock()
}

func (s *Server) respondWithError(c *gin.Context, code int, err error) {
	s.Logger.Error(err)
	c.AbortWithStatusJSON(code, gin.H{Error: err})
}

func setCookie(c *gin.Context, name, value string) {
	c.SetCookie(name, value, Hour, BasePath, Domain, true, false)
}

func performGetRequest(c *gin.Context, url string, ans any) error {
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status [%d] at [%s] received", resp.StatusCode, url)
	}

	return json.Unmarshal(content, ans)
}
