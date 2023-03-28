package controllers

import "github.com/liubkkkko/firstAPI/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.GET("/", s.Home)
	// s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.POST("/login", s.Login)
	// s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.POST("/users", s.CreateUser)
	s.Router.GET("/users", s.GetUser)
	s.Router.GET("/users/{id}", s.GetUser)
	s.Router.PUT("/users/{id}", s.UpdateUser, middlewares.SetMiddlewareAuthentication)
	s.Router.DELETE("/users/{id}", s.DeleteUser, middlewares.SetMiddlewareAuthentication)
	// s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	// s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	// s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	// s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	// s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Posts routes
	s.Router.POST("/posts", s.CreatePost)
	s.Router.GET("/posts", s.GetPosts)
	s.Router.GET("/posts/{id}", s.GetPost)
	s.Router.PUT("/posts/{id}", s.UpdatePost, middlewares.SetMiddlewareAuthentication)
	s.Router.DELETE("/posts/{id}", s.DeletePost, middlewares.SetMiddlewareAuthentication)
	// s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	// s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	// s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	// s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	// s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")

	//New routes

	// s.Router.HandleFunc("/")
}
