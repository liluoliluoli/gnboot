package router

import (
	"github.com/piupuer/go-helper/router"
	v1 "gnboot/api/v1"
)

func InitRoleRouter(r *router.Router) {
	router1 := r.Casbin("/role")
	router2 := r.CasbinAndIdempotence("/role")
	router1.GET("/list", v1.FindRole)
	router1.GET("/list/:ids", v1.FindRoleByIds)
	router2.POST("/create", v1.CreateRole)
	router1.PATCH("/update/:id", v1.UpdateRoleById)
	router1.DELETE("/delete/batch", v1.BatchDeleteRoleByIds)
}
