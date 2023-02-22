package Service

import (
	"DamniTkTok/JsonStruct"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/segmentio/ksuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		var msg *string
		Failmsg := "Failed to sava file"
		msg = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	outpath := VideoConverter(filePath)
	Host := c.Host()
	url := UrlHeader + string(Host) + Videorouting + outpath
	var Video JsonStruct.Video
	var User JsonStruct.User
	err = UserInfo.AutoMigrate(&JsonStruct.Video{}, &JsonStruct.User{})
	if err != nil {
		var msg *string
		Failmsg := "Failed to create a table"
		msg = &Failmsg
		c.JSON(consts.StatusOK, &JsonStruct.PublishRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	User = JsonStruct.User{
		WorkCount: userinfo.WorkCount,
	}
	Video = JsonStruct.Video{
		PlayURL: url,
		Title:   title,
		UserID:  userinfo.ID,
	}
	UserInfo.Create(&Video)
	result2 := UserInfo.Table("users").Where("token = ?", token).Update("work_count", User.WorkCount+1)
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
	token, tokenBool := c.GetQuery("token")
	_, user_idBool := c.GetQuery("user_id")

	if !tokenBool || !user_idBool {
		var msg *string
		Failmsg := "no passed data"
		msg = &Failmsg
		c.JSON(500, &JsonStruct.ListRsp{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	var userrsp = &JsonStruct.RspUser{}
	result := UserInfo.Table("users").Where("token = ?", token).First(userrsp)
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
	result2 := UserInfo.Table("videos").Limit(n).Where("user_id = ?", userrsp.ID).Find(&video)
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
	now := time.Now().UnixMilli()
	var last int64
	if latest_time == "0" {
		latest_time = strconv.FormatInt(now, 10)
		last, _ = strconv.ParseInt(latest_time, 10, 64)
	} //latest_time传入默认值时转换为当前时间

	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} //连接数据库
	var videos []*JsonStruct.RspVideo
	var video []JsonStruct.Video
	last, _ = strconv.ParseInt(latest_time, 10, 64)
	UserInfo.Table("videos").Where("created_at < ?", time.UnixMilli(last)).Limit(30).Order("created_at desc").Find(&video) //获取最新投稿列表（<=30)
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
	nextTime := video[(len(video) - 1)].CreatedAt.Add(time.Duration(-1)).UnixMilli() //获取最新投稿时间
	msg := "Success"
	c.JSON(consts.StatusOK, JsonStruct.FeedRsp{
		NextTime:   &nextTime,
		StatusCode: 0,
		StatusMsg:  &msg,
		VideoList:  videos,
	})
}
func ReadUser(userid int64) (u *JsonStruct.RspUser) {
	var User JsonStruct.RspUser
	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	result := UserInfo.Table("users").Where("id = ?", userid).First(&User)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	return &User
} //根据用户ID获取对应用户资料
