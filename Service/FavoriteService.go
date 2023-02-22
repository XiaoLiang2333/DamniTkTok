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
	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	result := UserInfo.Table("users").Where("token = ?", token).First(&userinfo)
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
	err = UserInfo.AutoMigrate(&JsonStruct.FavoriteList{})
	if err != nil {
		panic("failed to create table")
		return
	}

	// 通过 token 查询对应的用户ID —— result
	var userregister JsonStruct.User
	result = UserInfo.Where("token = ?", token).Find(&userregister)
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
		result := UserInfo.Create(&userfavorite)
		// 对应视频点赞数 +1
		var uservideo JsonStruct.Video
		UserInfo.Where("id = ?", videoId).Find(&uservideo)
		UserInfo.Table("videos").Where("id = ?", videoId).Update("favorite_count", uservideo.FavoriteCount+1)
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
		result := UserInfo.Where(map[string]interface{}{"user_id": userregister.ID, "video_id": int64VideoId}).Find(&userfavorite)
		// 对应视频点赞数 -1
		var uservideo JsonStruct.Video
		UserInfo.Where("id = ?", videoId).Find(&uservideo)
		UserInfo.Table("videos").Where("id = ?", videoId).Update("favorite_count", uservideo.FavoriteCount-1)
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
		UserInfo.Unscoped().Delete(&userfavorite)

		resp = &JsonStruct.FavoriteActionRsp{
			StatusCode: 0,
			StatusMsg:  "Operate Successfully",
		}
		c.JSON(consts.StatusOK, resp)
	}

}

// UserFavorList List 喜欢列表接口实现
func UserFavorList(ctx context.Context, c *app.RequestContext) {
	// 获取客户端参数 user_id token
	user_id, user_idBool := c.GetQuery("user_id")
	token, tokenBool := c.GetQuery("token")
	// 检查客户端参数 user_id token
	if !user_idBool || !tokenBool {
		Failmsg := "no passed data"
		c.JSON(500, &JsonStruct.FavoriteListRsp{
			StatusCode: 1,
			StatusMsg:  &Failmsg,
			VideoList:  nil,
		})
		return
	}
	fmt.Println(user_id)
	// 验证 token
	var userinfo JsonStruct.User
	TikTok, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	result := TikTok.Table("users").Where("token = ?", token).First(&userinfo)
	if result.Error != nil {
		Failmsg := "Wrong token"
		c.JSON(500, &JsonStruct.FavoriteListRsp{
			StatusCode: 1,
			StatusMsg:  &Failmsg,
			VideoList:  nil,
		})
		return
	}

	/* 获取用户的所有喜欢列表（点赞视频）*/
	var userfavorite []JsonStruct.FavoriteList
	result = TikTok.Table("favorite_lists").Where("user_id = ?", user_id).Find(&userfavorite)
	// 处理查询异常
	if result.Error != nil {
		msg := "Query Error"
		c.JSON(consts.StatusUnauthorized, &JsonStruct.FavoriteListRsp{
			StatusCode: 1,
			StatusMsg:  &msg,
			VideoList:  nil,
		})
		return
	}

	var videos []*JsonStruct.RspVideo
	TikTok.Table("favorite_lists").Where("user_id = ?", user_id).Find(&userfavorite)
	for _, v := range userfavorite {
		video := GetFavoriteVideo(v.VideoID)
		User := ReadUser(v.UserID)
		videos = append(videos, &JsonStruct.RspVideo{
			Author:        User,
			CommentCount:  video.CommentCount,
			CoverURL:      video.CoverURL,
			FavoriteCount: video.FavoriteCount,
			ID:            video.ID,
			IsFavorite:    video.IsFavorite,
			PlayURL:       video.PlayURL,
			Title:         video.Title,
		})
	} //提取用户喜欢列表并进行填装
	msg := "Success"
	c.JSON(consts.StatusOK, JsonStruct.FavoriteListRsp{
		StatusCode: 0,
		StatusMsg:  &msg,
		VideoList:  videos,
	})
}
func GetFavoriteVideo(videoid int64) (u *JsonStruct.Video) {
	TikTok, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var video JsonStruct.Video
	TikTok.Table("videos").Where("id = ?", videoid).First(&video)
	return &video
} //根据videoid提取视频
