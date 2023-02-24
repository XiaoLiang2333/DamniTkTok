package main

import (
	"DamniTkTok/JsonStruct"
	"DamniTkTok/Service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"os"
	"path"
)

func main() {
	h := server.Default(server.WithHostPorts(Service.Port), server.WithMaxRequestBodySize(2000*1024*1024))

	douyin := h.Group("/douyin")
	douyin.GET("/feed/", Service.FeedAction)
	userGroup := douyin.Group("/user")
	userGroup.POST("/register/", Service.Register)
	userGroup.POST("/login/", Service.Login)
	userGroup.GET("/", Service.Getinfo)
	//User services
	commentGroup := douyin.Group("/comment")
	commentGroup.POST("/action/", Service.CommentAction)
	commentGroup.GET("/list/", Service.CommentList)

	publishGroup := douyin.Group("/publish")
	publishGroup.POST("/action/", Service.PublishAction)
	publishGroup.GET("/list/", Service.List)
	publishGroup.GET("/list/:filename", func(c context.Context, ctx *app.RequestContext) {
		filename := ctx.Param("filename")
		filePath := path.Join("cover", filename+".jpg")
		filestream, err := os.ReadFile(filePath)
		if err != nil {
			ctx.JSON(500, &JsonStruct.PlayVideoRsp{
				StatusCode: 1,
				StatusMsg:  "Failed to turn to stream",
			})
			return
		}
		ctx.Write(filestream)
	})
	publishGroup.GET("/action/:filename", func(c context.Context, ctx *app.RequestContext) {
		filename := ctx.Param("filename")
		filePath := path.Join("out", filename+".mp4")
		filestream, err := os.ReadFile(filePath)
		if err != nil {
			ctx.JSON(500, &JsonStruct.PlayVideoRsp{
				StatusCode: 1,
				StatusMsg:  "Failed to turn to stream",
			})
			return
		}
		ctx.Write(filestream)
	})
	// Favorite Services
	favoriteGroup := douyin.Group("/favorite")
	favoriteGroup.POST("/action/", Service.FavorAction)
	favoriteGroup.GET("/list/", Service.UserFavorList)

	h.Spin()
}
