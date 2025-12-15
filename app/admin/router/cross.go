package router

import (
	"go-admin/app/admin/apis"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCrossRouter)
}

// registerCrossRouter
func registerCrossRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.CrossApi{}
	r := v1.Group("/cross") //.Use(authMiddleware.MiddlewareFunc()) //.Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
	}
}
