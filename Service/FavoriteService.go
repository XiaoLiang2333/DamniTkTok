package Service

import (
	"DamniTkTok/JsonStruct"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

// FavorAction Action 赞操作接口实现
func FavorAction(ctx context.Context, c *app.RequestContext) {
	// 获取客户端参数 token video_id action_type
	actionType, _ := c.GetQuery("action_type")

	token, _ := c.GetQuery("token")
	videoId, _ := c.GetQuery("video_id")
	int64VideoId, err := strconv.ParseInt(videoId, 10, 64)
	if len(actionType) == 0 || len(token) == 0 || len(videoId) == 0 {
		var msg *string
		Failmsg := "no passed data"
		msg = &Failmsg
		c.JSON(500, &JsonStruct.FavoriteActionRsp{
			StatusCode: 1,
			StatusMsg:  *msg,
		})
		return
	}

	// 验证 token
	var userinfo JsonStruct.User
	TikTok, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	result := TikTok.Table("users").Where("token = ?", token).First(&userinfo)
	if result.Error != nil {
		var msg *string
		Failmsg := "Wrong token"
		msg = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}

	// AutoMigrate 自动创建数据库表
	resp := &JsonStruct.FavoriteActionRsp{}
	var userfavorite JsonStruct.FavoriteList
	err = TikTok.AutoMigrate(&JsonStruct.FavoriteList{})
	if err != nil {
		panic("failed to create table")
		return
	}

	// 通过 token 查询对应的用户ID —— result
	var userregister JsonStruct.User
	result = TikTok.Where("token = ?", token).Find(&userregister)
	if result.Error != nil {
		var msg *string
		Failmsg := "User Not Found"
		msg = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.FavoriteActionRsp{
			StatusCode: 1,
			StatusMsg:  *msg,
		})
		return
	}

	switch actionType {
	// 将点赞操作插入数据库 action_type 1-点赞
	case strconv.Itoa(1):
		userfavorite = JsonStruct.FavoriteList{UserID: userregister.ID, VideoID: int64VideoId}
		result := TikTok.Create(&userfavorite)
		// 对应视频点赞数 +1
		var uservideo JsonStruct.Video
		TikTok.Where("id = ?", videoId).Find(&uservideo)
		TikTok.Table("videos").Where("id = ?", videoId).Update("favorite_count", uservideo.FavoriteCount+1)
		// 处理插入异常
		if result.Error != nil {
			var msg *string
			Failmsg := "Operation Failure"
			msg = &Failmsg
			c.JSON(consts.StatusUnauthorized, &JsonStruct.FavoriteActionRsp{
				StatusCode: 1,
				StatusMsg:  *msg,
			})
			return
		}
		// 正常执行
		resp = &JsonStruct.FavoriteActionRsp{
			StatusCode: 0,
			StatusMsg:  "Operate Successfully",
		}
		c.JSON(consts.StatusOK, resp)

	// 将取消点赞记录从数据库中删除 action_type 2-取消点赞
	case strconv.Itoa(2):
		result := TikTok.Where(map[string]interface{}{"user_id": userregister.ID, "video_id": int64VideoId}).Find(&userfavorite)
		// 对应视频点赞数 -1
		var uservideo JsonStruct.Video
		TikTok.Where("id = ?", videoId).Find(&uservideo)
		TikTok.Table("videos").Where("id = ?", videoId).Update("favorite_count", uservideo.FavoriteCount-1)
		// 处理查询异常
		if result.Error != nil {
			var msg *string
			Failmsg := "Query Error"
			msg = &Failmsg
			c.JSON(consts.StatusUnauthorized, &JsonStruct.FavoriteActionRsp{
				StatusCode: 1,
				StatusMsg:  *msg,
			})
			return
		}

		// 正常执行
		TikTok.Unscoped().Delete(&userfavorite)

		resp = &JsonStruct.FavoriteActionRsp{
			StatusCode: 0,
			StatusMsg:  "Operate Successfully",
		}
		c.JSON(consts.StatusOK, resp)
	}

}

// UserFavorList List 喜欢列表接口实现    还未完成开发
func UserFavorList(ctx context.Context, c *app.RequestContext) {
	// 获取客户端参数 user_id token
	userId, _ := c.GetQuery("user_id")
	token, _ := c.GetQuery("token")
	// 检查客户端参数 user_id token
	if len(userId) == 0 || len(token) == 0 {
		var msg *string
		Failmsg := "no passed data"
		msg = &Failmsg
		c.JSON(500, &JsonStruct.FavoriteActionRsp{
			StatusCode: 1,
			StatusMsg:  *msg,
		})
		return
	}

	// 验证 token
	var userinfo JsonStruct.User
	TikTok, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	result := TikTok.Table("users").Where("token = ?", token).First(&userinfo)
	if result.Error != nil {
		var msg *string
		Failmsg := "Wrong token"
		msg = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}

	/* 获取用户的所有喜欢列表（点赞视频）*/
	var userfavorite JsonStruct.FavoriteList
	result = TikTok.Where("user_id = ?", userId).Find(&userfavorite)
	// 处理查询异常
	if result.Error != nil {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.FavoriteActionRsp{
			StatusCode: 1,
			StatusMsg:  "Query Error",
		})
		return
	}
	// 正常执行
	fmt.Println(result.RowsAffected)
	/*
		// 返回响应
		resp := &JsonStruct.FavoriteListRsp{}
		var msg *string
		Failmsg := "Query Success"
		msg = &Failmsg
		resp = &JsonStruct.FavoriteListRsp{
			StatusCode: 0,
			StatusMsg: msg,
			VideoList:
		}
			c.JSON(consts.StatusOK, resp)
	*/
}
