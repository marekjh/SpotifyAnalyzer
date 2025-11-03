package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.TokenCache.Mutex.RLock()
		cache, ok := s.TokenCache.Data[c.ClientIP()]
		s.TokenCache.Mutex.RUnlock()

		if !ok {
			s.respondWithError(c, http.StatusUnauthorized, errors.New("please log in to access the service"))

			return
		}

		if cache.Token.Expiry.Before(time.Now()) {
			token, err := s.Authenticator.RefreshToken(c.Request.Context(), cache.Token)
			if err != nil {
				s.respondWithError(c, http.StatusInternalServerError, err)

				return
			}

			s.updateTokenCache(c, token)
		}

		setSpotifyClient(c, cache.SpotifyClient)
	}
}

type myRecentTracksRequest struct {
	Limit          int
	BeforeHoursAgo int64
	AfterHoursAgo  int64
}

func (s *Server) myRecentTracksMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req myRecentTracksRequest

		err := c.ShouldBindQuery(&req)
		if err != nil {
			s.respondWithError(c, http.StatusBadRequest, err)

			return
		}

		setMyRecentTracksRequest(c, &req)
	}
}
