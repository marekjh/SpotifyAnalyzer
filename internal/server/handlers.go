package server

import (
	"crypto/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
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
		client, err := retrieveSpotifyClient(c)
		if err != nil {
			s.respondWithError(c, http.StatusInternalServerError, err)

			return
		}

		user, err := client.CurrentUser(c.Request.Context())
		if err != nil {
			s.respondWithError(c, http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (s *Server) handleMyRecentTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		client, err := retrieveSpotifyClient(c)
		if err != nil {
			s.respondWithError(c, http.StatusInternalServerError, err)

			return
		}

		req, err := retrieveMyRecentTracksRequest(c)
		if err != nil {
			s.respondWithError(c, http.StatusInternalServerError, err)

			return
		}

		before := req.beforeHoursAgo
		if req.beforeHoursAgo != 0 {
			before = time.Now().Add(-time.Duration(req.beforeHoursAgo) * time.Hour).UnixMilli()
		}

		after := req.afterHoursAgo
		if req.afterHoursAgo != 0 {
			after = time.Now().Add(-time.Duration(req.afterHoursAgo) * time.Hour).UnixMilli()
		}

		opts := &spotify.RecentlyPlayedOptions{
			Limit:         spotify.Numeric(req.limit),
			BeforeEpochMs: before,
			AfterEpochMs:  after,
		}

		recentTracks, err := client.PlayerRecentlyPlayedOpt(c.Request.Context(), opts)
		if err != nil {
			s.respondWithError(c, http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusOK, recentTracks)
	}
}
