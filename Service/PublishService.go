package Service

import (
	"DamniTkTok/JsonStruct"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/segmentio/ksuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path"
)

func PublishAction(ctx context.Context, c *app.RequestContext) {
	form, _ := c.MultipartForm()
	token := form.Value["token"][0]
	title := form.Value["title"][0]
	file := form.File["data"]
	var userinfo JsonStruct.User
	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
	url := "http://localhost:8080/douyin/publish/action/" + filename
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
		StatusCode: 1,
		StatusMsg:  msg,
	})
} //此为投稿接口
func List(ctx context.Context, c *app.RequestContext) {
	_, tokenBool := c.GetQuery("token")
	user_id, user_idBool := c.GetQuery("user_id")
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
	var userinfo = JsonStruct.User{}
	result := UserInfo.Table("users").Where("id = ?", user_id).First(&userinfo)
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
	var msg *string
	Failmsg := "Success"
	msg = &Failmsg
	c.JSON(500, &JsonStruct.ListRsp{
		StatusCode: 200,
		StatusMsg:  msg,
	})
} //此为发布列表接口，但存有BUG，18日24.前会调试完成，暂且不要调用
