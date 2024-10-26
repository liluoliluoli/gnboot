package router

import (
	"github.com/piupuer/go-helper/router"
	v1 "gnboot/api/v1"
)

func InitLeaveRouter(r *router.Router) {
	router1 := r.Casbin("/leave")
	router2 := r.CasbinAndIdempotence("/leave")
	router1.GET("/list", v1.FindLeave)
	router2.POST("/create", v1.CreateLeave)
	router1.PATCH("/update/:id", v1.UpdateLeaveById)
	router1.DELETE("/delete/batch", v1.BatchDeleteLeaveByIds)
}
