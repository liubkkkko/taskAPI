package controllers

import (
	"github.com/liubkkkko/firstAPI/api/middlewares"
)

func (server *Server) initializeRoutes() {

    // Healthcheck (no auth)
    server.Router.GET("/health", server.Health)

	// Home Route
	server.Router.GET("/", server.Home)
	server.Router.GET("/test", server.TestRout)
	server.Router.GET("/author", server.IdIfYouHaveToken, middlewares.SetMiddlewareAuthentication)

	// Login/Refresh/Logout
	server.Router.POST("/login", server.Login)
	server.Router.POST("/refresh", server.Refresh) // тут credentials: include використовує cookie
	server.Router.POST("/logout", server.Logout, middlewares.SetMiddlewareAuthentication)

	// Author routes
	server.Router.POST("/authors", server.CreateAuthor)
	server.Router.GET("/authors", server.GetAuthors, middlewares.SetMiddlewareAuthentication)
	server.Router.GET("/author/:id", server.GetAuthor, middlewares.SetMiddlewareAuthentication)
	server.Router.PUT("/authors/:id", server.UpdateAuthor, middlewares.SetMiddlewareAuthentication)
	server.Router.DELETE("/authors/:id", server.DeleteAuthor, middlewares.SetMiddlewareAuthentication)

	// Workspace routes (більшість доступні тільки з авторизацією)
	server.Router.POST("/workspaces", server.CreateWorkspace, middlewares.SetMiddlewareAuthentication)
	server.Router.GET("/workspaces", server.GetWorkspaces, middlewares.SetMiddlewareAuthentication)
	server.Router.GET("/workspace/:id", server.GetWorkspace, middlewares.SetMiddlewareAuthentication)
	server.Router.GET("/workspaces/authors/:id", server.GetWorkspacesByAuthorId, middlewares.SetMiddlewareAuthentication)
	server.Router.PUT("/workspace/:id", server.AddOneMoreAuthorToWorkspace, middlewares.SetMiddlewareAuthentication)
	server.Router.PUT("/workspaces/:id", server.UpdateWorkspace, middlewares.SetMiddlewareAuthentication)
	server.Router.DELETE("/workspaces/:id", server.DeleteWorkspace, middlewares.SetMiddlewareAuthentication)
	server.Router.GET("/workspace", server.CheckIfIAuthor, middlewares.SetMiddlewareAuthentication)

	// Job routes
	server.Router.POST("/jobs", server.CreateJob, middlewares.SetMiddlewareAuthentication)
	server.Router.GET("/jobs", server.GetJobs, middlewares.SetMiddlewareAuthentication)
	server.Router.GET("/job/:id", server.GetJob, middlewares.SetMiddlewareAuthentication)
	server.Router.GET("/jobs/:id", server.GetJobsByWorkspaceId, middlewares.SetMiddlewareAuthentication)
	server.Router.PUT("/job/:id", server.UpdateJob, middlewares.SetMiddlewareAuthentication)
	server.Router.DELETE("/job/:id", server.DeleteJob, middlewares.SetMiddlewareAuthentication)

	// ...existing code...
    server.Router.GET("/auth/google", server.GoogleAuthRedirect)
    server.Router.GET("/auth/google/callback", server.GoogleAuthCallback)

}
