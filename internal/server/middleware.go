package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) refreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.TokenCache.Mutex.RLock()
		token, ok := s.TokenCache.Data[c.ClientIP()]
		s.TokenCache.Mutex.RUnlock()

		if !ok {
			s.respondWithError(c, http.StatusUnauthorized, errors.New("please login to access the service"))

			return
		}

		var err error
		if token.Expiry.Before(time.Now()) {
			token, err = s.Authenticator.RefreshToken(c.Request.Context(), token)
			if err != nil {
				s.respondWithError(c, http.StatusInternalServerError, err)

				return
			}
		}

		s.updateTokenCache(c, token)
	}
}
