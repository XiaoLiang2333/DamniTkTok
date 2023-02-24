package Service

import (
	"DamniTkTok/Database"
	"DamniTkTok/JsonStruct"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func CommentAction(ctx context.Context, c *app.RequestContext) {
	token, tokenBool := c.GetQuery("token")
	video_id, video_idBool := c.GetQuery("video_id")
	action_type, action_typeBool := c.GetQuery("action_type")
	var ms1 *string
	Failmsg := "Wrong token"
	ms1 = &Failmsg
	uncom := JsonStruct.CommentRes{}
	if !tokenBool || !video_idBool || !action_typeBool {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.CommentRsp{
			StatusCode: 1,
			StatusMsg:  ms1,
			Comment:    &uncom,
		})
		return
	}
	var userinfo JsonStruct.User
	result := Database.DB.Table("users").Where("token = ?", token).First(&userinfo)
	if result.Error != nil {
		var ms2 *string
		Failmsg := "Wrong token"
		ms2 = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.CommentRsp{
			StatusCode: 1,
			StatusMsg:  ms2,
			Comment:    &uncom,
		})
		return
	}

	err := Database.DB.AutoMigrate(&JsonStruct.Comment{}, &JsonStruct.User{})
	if err != nil {
		var ms3 *string
		Failmsg := "Failed to create a table"
		ms3 = &Failmsg
		c.JSON(consts.StatusOK, &JsonStruct.CommentRsp{
			StatusCode: 1,
			StatusMsg:  ms3,
			Comment:    &uncom,
		})
		return
	}
	actiontype, _ := strconv.Atoi(action_type)

	//1
	if actiontype == 1 {
		comment_text, comment_textBool := c.GetQuery("comment_text")
		var ms4 *string
		Failmsg := "comment_text err"
		ms4 = &Failmsg
		if !comment_textBool {
			c.JSON(consts.StatusUnauthorized, &JsonStruct.CommentRsp{
				StatusCode: 1,
				StatusMsg:  ms4,
				Comment:    &uncom,
			})
			return
		}
		t := time.Now().Format("2006-01-02")
		m, _ := strconv.Atoi(video_id)
		var Comment JsonStruct.Comment
		Comment = JsonStruct.Comment{
			Content:    comment_text,
			CreateDate: t,
			UserID:     userinfo.ID,
			VideoID:    int64(m),
		}

		result := Database.DB.Table("comments").Create(&Comment)
		if result.Error != nil {
			var ms5 *string
			Failmsg := "add err"
			ms5 = &Failmsg
			c.JSON(consts.StatusOK, &JsonStruct.CommentRsp{
				StatusCode: 1,
				StatusMsg:  ms5,
				Comment:    &uncom,
			})
			return
		}

		var com JsonStruct.Video
		result = Database.DB.Table("videos").Where("id = ?", video_id).First(&com)
		if result.Error != nil {
			var ms5 *string
			Failmsg := "video_id err"
			ms5 = &Failmsg
			c.JSON(consts.StatusOK, &JsonStruct.CommentRsp{
				StatusCode: 1,
				StatusMsg:  ms5,
				Comment:    &uncom,
			})
			return
		}
		Video := JsonStruct.Video{
			CommentCount: com.CommentCount,
		}
		result4 := Database.DB.Table("videos").Where("id = ?", video_id).Update("favorite_count", Video.CommentCount+1)
		if result4.Error != nil {
			var ms6 *string
			Failmsg := "Failed to Update"
			ms6 = &Failmsg
			c.JSON(500, &JsonStruct.CommentRsp{
				StatusCode: 1,
				StatusMsg:  ms6,
				Comment:    &uncom,
			})
			return
		}
		var ms7 *string
		commentuserres := JsonStruct.RspUser{
			userinfo.ID,
			userinfo.Name,
			userinfo.FollowCount,
			userinfo.FollowerCount,
			userinfo.IsFollow,
			userinfo.Avatar,
			userinfo.BackgroundImage,
			userinfo.Signature,
			userinfo.TotalFavorited,
			userinfo.WorkCount,
			userinfo.FavoriteCount,
		}

		CommentRes := JsonStruct.CommentRes{
			Content:    comment_text,
			CreateDate: t,
			ID:         Comment.ID,
			User:       commentuserres,
		}

		Failmsg = "Success"
		ms7 = &Failmsg
		c.JSON(consts.StatusOK, &JsonStruct.CommentRsp{
			StatusCode: 0,
			StatusMsg:  ms7,
			Comment:    &CommentRes,
		})

	} else if actiontype == 2 {
		comment_id, comment_idBool := c.GetQuery("comment_id")
		if !comment_idBool {
			var ms8 *string
			Failmsg := "comment_id err"
			ms8 = &Failmsg
			c.JSON(consts.StatusUnauthorized, &JsonStruct.CommentRsp{

				StatusCode: 1,
				StatusMsg:  ms8,
				Comment:    &uncom,
			})
			return
		}
		Database.DB.Delete(&JsonStruct.Comment{}, comment_id)
		var ms9 *string
		Failmsg := "Success"
		ms9 = &Failmsg
		c.JSON(consts.StatusOK, &JsonStruct.CommentRsp{
			StatusCode: 0,
			StatusMsg:  ms9,
			Comment:    &uncom,
		})
	} else {
		var ms10 *string
		Failmsg := "actiontype err"
		ms10 = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.CommentRsp{
			StatusCode: 1,
			StatusMsg:  ms10,
			Comment:    &uncom,
		})
	}

}

// 此为评论视频接口，已测试，无问题
func CommentList(ctx context.Context, c *app.RequestContext) {
	_, tokenBool := c.GetQuery("token")
	video_id, video_idBool := c.GetQuery("video_id")

	comment := []JsonStruct.Comment{}

	if !tokenBool || !video_idBool {
		var msg *string
		Failmsg := "no passed data"
		msg = &Failmsg
		c.JSON(500, &JsonStruct.CommentListRsp{
			StatusCode:  1,
			StatusMsg:   msg,
			CommentList: comment,
		})
		return
	}
	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	comments := make([]JsonStruct.Comment, 8)
	UserInfo.Table("comments").Where("video_id = ?", video_id).Find(&comments)
	successInfo := "get commentlist successfully"
	success := &successInfo
	c.JSON(consts.StatusOK, &JsonStruct.CommentListRsp{
		CommentList: comments,
		StatusCode:  0,
		StatusMsg:   success,
	})

}

//此为拉取评论列表接口，已测试，无问题
