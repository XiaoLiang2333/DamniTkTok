package Service

import (
	"DamniTkTok/Database"
	"DamniTkTok/JsonStruct"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/segmentio/ksuid"

	"os"
	"path"
	"strconv"
	"time"
)

func PublishAction(ctx context.Context, c *app.RequestContext) {
	form, _ := c.MultipartForm()
	token := form.Value["token"][0]
	title := form.Value["title"][0]
	file := form.File["data"]
	var userinfo JsonStruct.User
	result := Database.DB.Table("users").Where("token = ?", token).First(&userinfo)
	if result.Error != nil {
		Failmsg := "Wrong token"
		c.JSON(consts.StatusUnauthorized, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  &Failmsg,
		})
		return
	}
	opened, _ := file[0].Open()
	var data = make([]byte, file[0].Size)
	ReadSize, err := opened.Read(data)
	if err != nil {
		var msg *string
		Failmsg := "Open file failed"
		msg = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	if ReadSize != int(file[0].Size) {
		var msg *string
		Failmsg := "Size not match"
		msg = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	filename := ksuid.New().String()
	filePath := path.Join("video", filename+".mp4")
	dir := path.Dir(filePath)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		var msg *string
		Failmsg := "Failed to create directory"
		msg = &Failmsg
		c.JSON(500, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}

	err = os.WriteFile(filePath, data, os.FileMode(0755))
	if err != nil {
		Failmsg := "Failed to sava file"
		c.JSON(consts.StatusUnauthorized, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  &Failmsg,
		})
		return
	}
	cover := Cover(filePath)
	outpath := VideoConverter(filePath)
	Host := c.Host()
	coverurl := UrlHeader + string(Host) + CoverRouting + cover
	url := UrlHeader + string(Host) + Videorouting + outpath
	var Video JsonStruct.Video
	err = Database.DB.AutoMigrate(&JsonStruct.Video{})
	if err != nil {
		Failmsg := "Failed to create a table"
		c.JSON(500, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  &Failmsg,
		})
		return
	}
	Video = JsonStruct.Video{
		PlayURL:  url,
		Title:    title,
		UserID:   userinfo.ID,
		CoverURL: coverurl,
	}
	Database.DB.Save(&Video)
	result2 := Database.DB.Table("users").Where("token = ?", token).Update("work_count", userinfo.WorkCount+1)
	if result2.Error != nil {
		var msg *string
		Failmsg := "Failed to Update"
		msg = &Failmsg
		c.JSON(500, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	var msg *string
	Failmsg := "Success"
	msg = &Failmsg
	c.JSON(consts.StatusOK, &JsonStruct.PublishRsp{
		StatusCode: 0,
		StatusMsg:  msg,
	})
} //此为投稿接口
func List(ctx context.Context, c *app.RequestContext) {
	_, tokenBool := c.GetQuery("token")
	user_id, user_idBool := c.GetQuery("user_id")
	userid, _ := strconv.ParseInt(user_id, 10, 64)
	if !tokenBool || !user_idBool {
		Failmsg := "no passed data"
		c.JSON(500, &JsonStruct.ListRsp{
			StatusCode: 1,
			StatusMsg:  &Failmsg,
		})
		return
	}
	var userrsp = JsonStruct.User{}
	result := Database.DB.Table("users").Where("id = ?", userid).First(&userrsp)
	if result.Error != nil {
		var msg *string
		Failmsg := "no passed data"
		msg = &Failmsg
		c.JSON(500, &JsonStruct.ListRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	n, _ := strconv.Atoi(strconv.FormatInt(userrsp.WorkCount, 10))
	var video []JsonStruct.Video
	result2 := Database.DB.Table("videos").Limit(n).Where("user_id = ?", userrsp.ID).Find(&video)
	if result2.Error != nil {
		var msg *string
		Failmsg := "Fail to find the videolist"
		msg = &Failmsg
		c.JSON(500, &JsonStruct.ListRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	var videos []*JsonStruct.RspVideo
	for _, v := range video {
		User := ReadUser(v.UserID)
		videos = append(videos, &JsonStruct.RspVideo{
			Author:        User,
			CommentCount:  v.CommentCount,
			CoverURL:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			ID:            v.ID,
			IsFavorite:    v.IsFavorite,
			PlayURL:       v.PlayURL,
			Title:         v.Title,
		})
	}
	Failmsg := "Success"
	c.JSON(200, &JsonStruct.ListRsp{
		StatusCode: 200,
		StatusMsg:  &Failmsg,
		VideoList:  videos,
	})
}
func FeedAction(ctx context.Context, c *app.RequestContext) {
	latest_time := c.Query("latest_time")

	var last int64
	now := time.Now().UnixMilli()
	if latest_time == "0" {

		latest_time = strconv.FormatInt(now, 10)
		last, _ = strconv.ParseInt(latest_time, 10, 64)
	} else {
		last, _ = strconv.ParseInt(latest_time, 10, 64)
	}
	d5415 := time.UnixMilli(last)
	fmt.Println(d5415)
	//latest_time传入默认值时转换为当前时间
	var Limitvideo []JsonStruct.Video
	result := Database.DB.Table("videos").Where("created_at >= ?", time.UnixMilli(last)).Order("created_at desc").Find(&Limitvideo)
	if result.RowsAffected > 1 || last == now {
		var count []JsonStruct.Video
		Vcount := Database.DB.Table("videos").Find(&count)
		var videos []*JsonStruct.RspVideo
		video := make([]JsonStruct.Video, Vcount.RowsAffected)
		now := time.Now().UnixMilli()
		last, _ = strconv.ParseInt(strconv.FormatInt(now, 10), 10, 64)
		Database.DB.Table("videos").Where("created_at < ?", time.UnixMilli(last)).Limit(30).Order("created_at desc").Find(&video)
		//获取最新投稿列表（<=30)
		for _, v := range video {
			User := ReadUser(v.UserID)
			videos = append(videos, &JsonStruct.RspVideo{
				Author:        User,
				CommentCount:  v.CommentCount,
				CoverURL:      v.CoverURL,
				FavoriteCount: v.FavoriteCount,
				ID:            v.ID,
				IsFavorite:    v.IsFavorite,
				PlayURL:       v.PlayURL,
				Title:         v.Title,
			})
		} //对应数据填入RspvideoModel
		var nextTime int64
		if len(video) == 0 {

			msg := "Already the newest"
			c.JSON(consts.StatusOK, JsonStruct.FeedRsp{
				NextTime:   nil,
				StatusCode: 1,
				StatusMsg:  &msg,
				VideoList:  nil,
			})
			return
		} //防止反复刷新导致数组越界
		nextTime = video[0].CreatedAt.UnixNano() / int64(time.Millisecond)
		msg := "Success"
		c.JSON(consts.StatusOK, JsonStruct.FeedRsp{

			NextTime:   &nextTime,
			StatusCode: 0,
			StatusMsg:  &msg,
			VideoList:  videos,
		})
	} else {
		msg := "fail"
		c.JSON(500, JsonStruct.FeedRsp{
			NextTime:   nil,
			StatusCode: 1,
			StatusMsg:  &msg,
			VideoList:  nil,
		})
	}
}

func ReadUser(userid int64) (u *JsonStruct.RspUser) {
	var User JsonStruct.RspUser
	result := Database.DB.Table("users").Where("id = ?", userid).First(&User)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	return &User
} //根据用户ID获取对应用户资料
