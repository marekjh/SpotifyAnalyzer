package server

func (s *Server) addRoutes() {
	s.Engine.GET("/login", s.handleLogin())

	s.Engine.GET("/auth-response", s.handleAuthResponse())

	apiV1 := s.Engine.Group("/api/v1")
	apiV1.Use(s.refreshToken())
	
	apiV1.GET("/myaccount", )
}
