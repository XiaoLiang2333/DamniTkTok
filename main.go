package main

import (
	"DamniTkTok/Service"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default(server.WithHostPorts(Service.Port))

	douyin := h.Group("/douyin")
	userGroup := douyin.Group("/user")
	userGroup.POST("/register/", Service.Register)
	userGroup.POST("/login/", Service.Login)
	h.Spin()
}
