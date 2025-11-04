package server

func (s *Server) addRoutes() {
	s.Engine.GET("/login", s.handleLogin())

	s.Engine.GET("/callback", s.handleCallback())

	apiV1 := s.Engine.Group("/api/v1")
	apiV1.Use(s.authMiddleware())
	apiV1.GET("/myaccount", s.handleMyAccount())
	apiV1.GET("/myrecenttracks", s.myRecentTracksMiddleware(), s.handleMyRecentTracks())
	apiV1.GET("/mydevices", s.handleMyDevices())
	apiV1.PUT("/play", s.playMiddleware(), s.handlePlay())
}
