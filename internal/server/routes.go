package server

func (s *Server) addRoutes() {
	s.Engine.GET("/login", s.handleLogin())

	s.Engine.GET("/auth-response", s.handleAuthResponse())
}
