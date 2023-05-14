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
	s.Router.PUT("/users/:id", s.UpdateUser, middlewares.SetMiddlewareAuthentication)    //WORKING
	s.Router.DELETE("/users/:id", s.DeleteUser, middlewares.SetMiddlewareAuthentication) //WORKING

	//Posts routes
	s.Router.POST("/posts", s.CreatePost, middlewares.SetMiddlewareAuthentication)       //work if you use token
	s.Router.GET("/posts", s.GetPosts)                                                   //working
	s.Router.GET("/posts/:id", s.GetPost)                                                //working
	s.Router.PUT("/posts/:id", s.UpdatePost, middlewares.SetMiddlewareAuthentication)    //working (only if you Update your own post)
	s.Router.DELETE("/posts/:id", s.DeletePost, middlewares.SetMiddlewareAuthentication) //working (only if you Delete your own post)

	//Task routes
	s.Router.POST("/tasks", s.CreateTask) //working
	s.Router.GET("/tasks", s.GetTasks)    //working
}
