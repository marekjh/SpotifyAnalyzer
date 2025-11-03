package server

import (
	"crypto/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

const AuthCookie = "auth-cookie"

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := rand.Text()
		authURL := s.Authenticator.AuthURL(state)

		setCookie(c, AuthCookie, state)

		c.String(http.StatusOK, authURL)
	}
}

func (s *Server) handleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := c.Cookie(AuthCookie)
		if err != nil {
			s.respondWithError(c, http.StatusInternalServerError, err)

			return
		}

		token, err := s.Authenticator.Token(c.Request.Context(), state, c.Request)
		if err != nil {
			s.respondWithError(c, http.StatusInternalServerError, err)

			return
		}

		s.updateTokenCache(c, token)
	}
}

func (s *Server) handleMyAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.SpotifyClient.CurrentUser(c.Request.Context())
		if err != nil {
			s.respondWithError(c, http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusOK, user)
	}
}
