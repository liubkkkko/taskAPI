package controllers

import (
	"github.com/liubkkkko/firstAPI/api/middlewares"
)

func (server *Server) initializeRoutes() {

	// Home Route
	server.Router.GET("/", server.Home)         // working
	server.Router.GET("/test", server.TestRout) // working
	server.Router.GET("/author", server.IdIfYouHaveToken)

	// Login Route
	server.Router.POST("/login", server.Login)                                            //working
	server.Router.POST("/logout", server.Logout, middlewares.SetMiddlewareAuthentication) //working

	//Author routes
	server.Router.POST("/authors", server.CreateAuthor)                                                //working
	server.Router.GET("/authors", server.GetAuthors, middlewares.SetMiddlewareAuthentication)                                                   //working
	server.Router.GET("/author/:id", server.GetAuthor, middlewares.SetMiddlewareAuthentication)                                                 //working
	server.Router.PUT("/authors/:id", server.UpdateAuthor, middlewares.SetMiddlewareAuthentication)    //working
	server.Router.DELETE("/authors/:id", server.DeleteAuthor, middlewares.SetMiddlewareAuthentication) //working

	//Workspace routes
	server.Router.POST("/workspces", server.CreateWorspace, middlewares.SetMiddlewareAuthentication)                 //working
	server.Router.GET("/workspaces", server.GetWorkspaces, middlewares.SetMiddlewareAuthentication)                                                           //working
	server.Router.GET("/workspace/:id", server.GetWorkspace, middlewares.SetMiddlewareAuthentication)                                                         //working
	server.Router.GET("/workspaces/authors/:id", server.GetWorkspacesByAuthorId, middlewares.SetMiddlewareAuthentication)                                     //working
	server.Router.PUT("/workspace/:id", server.AddOneMoreAuthorToWorkspace, middlewares.SetMiddlewareAuthentication) //working
	server.Router.PUT("/workspaces/:id", server.UpdateWorkspace, middlewares.SetMiddlewareAuthentication)            //working
	server.Router.DELETE("/workspaces/:id", server.DeleteWorkspace, middlewares.SetMiddlewareAuthentication)         //working (only if you try to delete own workspace)
	server.Router.GET("/workspace", server.CheckIfIAuthor, middlewares.SetMiddlewareAuthentication)                                                           //working

	//Job routes
	server.Router.POST("/jobs", server.CreateJob, middlewares.SetMiddlewareAuthentication)      //working
	server.Router.GET("/jobs", server.GetJobs, middlewares.SetMiddlewareAuthentication)                                                  //working
	server.Router.GET("/job/:id", server.GetJob, middlewares.SetMiddlewareAuthentication)                                                //working
	server.Router.GET("/jobs/:id", server.GetJobsByWorkspaceId, middlewares.SetMiddlewareAuthentication)
	server.Router.PUT("/job/:id", server.UpdateJob, middlewares.SetMiddlewareAuthentication)    //working
	server.Router.DELETE("/job/:id", server.DeleteJob, middlewares.SetMiddlewareAuthentication) //working
}
