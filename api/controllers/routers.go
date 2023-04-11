package controllers

import (
	"github.com/liubkkkko/firstAPI/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.GET("/", s.Home)

	// Login Route
	s.Router.POST("/login", s.Login)

	//Users routes
	s.Router.POST("/users", s.CreateUser)
	s.Router.GET("/users", s.GetUser)
	s.Router.GET("/users/{id}", s.GetUser)
	s.Router.PUT("/users/{id}", s.UpdateUser, middlewares.SetMiddlewareAuthentication)
	s.Router.DELETE("/users/{id}", s.DeleteUser, middlewares.SetMiddlewareAuthentication)

	//Posts routes
	s.Router.POST("/posts", s.CreatePost)
	s.Router.GET("/posts", s.GetPosts)
	s.Router.GET("/posts/{id}", s.GetPost)
	s.Router.PUT("/posts/{id}", s.UpdatePost, middlewares.SetMiddlewareAuthentication)
	s.Router.DELETE("/posts/{id}", s.DeletePost, middlewares.SetMiddlewareAuthentication)
}
