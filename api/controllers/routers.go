package controllers

import (
	"github.com/liubkkkko/firstAPI/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.GET("/", s.Home)  // working
	s.Router.GET("/test", s.TestRout) // working

	// Login Route
	s.Router.POST("/login", s.Login) //working
	s.Router.POST("/logout", s.Logout) //working

	//Author routes
	s.Router.POST("/authors", s.CreateAuthor) //working
	s.Router.GET("/authors", s.GetAuthors)    //working
	s.Router.GET("/author/:id", s.GetAuthor)  //working
	s.Router.PUT("/authors/:id", s.UpdateAuthor, middlewares.SetMiddlewareAuthentication)  //working
	s.Router.DELETE("/authors/:id", s.DeleteAuthor, middlewares.SetMiddlewareAuthentication)  //working

	//Workspace routes
	s.Router.POST("/workspces", s.CreateWorspace, middlewares.SetMiddlewareAuthentication) //working
	s.Router.GET("/workspaces", s.GetWorkspaces) //working
	s.Router.GET("/workspace/:id", s.GetWorkspace) //working
	s.Router.GET("/workspaces/authors/:id", s.GetWorkspacesByAuthorId) //working
	s.Router.PUT("/workspace/:id", s.AddOneMoreAuthorToWorkspace, middlewares.SetMiddlewareAuthentication) //working
	s.Router.PUT("/workspaces/:id", s.UpdateWorkspace, middlewares.SetMiddlewareAuthentication) //working
	s.Router.DELETE("/workspaces/:id", s.DeleteWorkspace, middlewares.SetMiddlewareAuthentication) //working (only if you try to delete own workspace)
	s.Router.GET("/workspace", s.CheckIfIAuthor) //working

	//Job routes
	s.Router.POST("/jobs", s.CreateJob, middlewares.SetMiddlewareAuthentication) //working
	s.Router.GET("/jobs", s.GetJobs) //working
	s.Router.GET("/job/:id", s.GetJob) //working
	s.Router.PUT("/job/:id", s.UpdateJob, middlewares.SetMiddlewareAuthentication) //working
	s.Router.DELETE("/job/:id", s.DeleteJob, middlewares.SetMiddlewareAuthentication) //working
}
