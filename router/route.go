package router

import (
	"github.com/devesh44/gin-poc/handler"
	"github.com/gin-gonic/gin"
)

type router struct {
	r *gin.Engine
}

func New(r *gin.Engine) *router {
	return &router{r}
}

func (rt *router) RegisterRoutes() {
	rg := rt.r.Group("/v1/api")
	rt.auth(rg)
	// rt.order(rg)
}
func (rt *router) auth(rg *gin.RouterGroup) {
	ah := handler.NewAuthHandler()
	r := rg.Group("/auth")
	{
		r.POST("/register", ah.Create())
	}

}

// func (rt *router) order(rg *gin.RouterGroup) {
// 	handler.NewOrderHandler()
// 	rt.r.POST("/order")
// }
