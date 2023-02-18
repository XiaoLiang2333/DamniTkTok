package main

import (
	"DamniTkTok/JsonStruct"
	"DamniTkTok/Service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"os"
	"path"
)

func main() {
	h := server.Default(server.WithHostPorts(Service.Port))

	douyin := h.Group("/douyin")
	userGroup := douyin.Group("/user")
	userGroup.POST("/register/", Service.Register)
	userGroup.POST("/login/", Service.Login)
	userGroup.GET("/", Service.Getinfo)
	//User services
	publishGroup := douyin.Group("/publish")
	publishGroup.POST("/action/", Service.PublishAction)
	publishGroup.GET("/list/", Service.List)
	publishGroup.POST("/action/:filename", func(c context.Context, ctx *app.RequestContext) {
		filename := ctx.Param("filename")
		filePath := path.Join("video", filename+".mp4")
		filestream, err := os.ReadFile(filePath)
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, &JsonStruct.PlayVideoRsp{
				StatusCode: 1,
				StatusMsg:  "Failed to turn to stream",
			})
			return
		}
		ctx.Write(filestream)
	})
	h.Spin()
}