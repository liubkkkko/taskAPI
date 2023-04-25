package controllers

import (
	"github.com/liubkkkko/firstAPI/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route // working
	s.Router.GET("/", s.Home)

	// working
	s.Router.GET("/test", s.TestRout)

	// Login Route
	s.Router.POST("/login", s.Login)

	//Users routes
	s.Router.POST("/users", s.CreateUser)                                                //working
	s.Router.GET("/users", s.GetUsers)                                                   //WORKING
	s.Router.GET("/users/:id", s.GetUser)                                                //WORKING
	s.Router.PUT("/users/:id", s.UpdateUser, middlewares.SetMiddlewareAuthentication)    //NOT WORKING
	s.Router.DELETE("/users/:id", s.DeleteUser, middlewares.SetMiddlewareAuthentication) //NOT WORKING

	//Posts routes
	s.Router.POST("/posts", s.CreatePost)                                                //seems to work
	s.Router.GET("/posts", s.GetPosts)                                                   //working
	s.Router.GET("/posts/:id", s.GetPost)                                                //NOT WORKING
	s.Router.PUT("/posts/:id", s.UpdatePost, middlewares.SetMiddlewareAuthentication)    //no info
	s.Router.DELETE("/posts/:id", s.DeletePost, middlewares.SetMiddlewareAuthentication) //no info
}
