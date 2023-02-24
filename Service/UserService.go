package Service

import (
	"DamniTkTok/Database"
	"DamniTkTok/JsonStruct"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/segmentio/ksuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

func Register(ctx context.Context, c *app.RequestContext) {
	username, usernameBool := c.GetQuery("username")
	password, passwordBool := c.GetQuery("password")
	if !usernameBool || !passwordBool {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.RegisterResponse{
			StatusCode: 1,
			StatusMsg:  "no passed data",
		})
		return
	}
	if len(username) > 32 || len(password) > 32 {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.RegisterResponse{
			StatusCode: 1,
			StatusMsg:  "Input is too long(>32)",
		})
		return
	}
	resp := &JsonStruct.RegisterResponse{}
	var userregister JsonStruct.User
	Database.DB.Table("users").AutoMigrate(&userregister)
	result := Database.DB.Table("users").Where("name = ?", username).First(&userregister)
	if result.Error == nil {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.RegisterResponse{
			StatusCode: 1,
			StatusMsg:  "Already existed",
		})
		return
	}
	token := ksuid.New().String()
	userregister = JsonStruct.User{Name: username, UserPassWord: password, Token: token}
	Database.DB.Table("users").Create(&userregister)
	resp = &JsonStruct.RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "Register Query Success",
		Token:      token,
		UserID:     userregister.ID,
	}
	c.JSON(consts.StatusOK, resp)
} //此为注册接口对应的基础服务实现

func Login(ctx context.Context, c *app.RequestContext) {
	username, usernameBool := c.GetQuery("username")
	password, passwordBool := c.GetQuery("password")
	var userregister JsonStruct.User
	var usertoken *string
	var msg *string
	if !usernameBool || !passwordBool {
		failmsg := "no passed data"
		c.JSON(consts.StatusUnauthorized, &JsonStruct.LoginResponse{
			StatusCode: 1,
			StatusMsg:  &failmsg,
		})
		return
	}
	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	resp := &JsonStruct.LoginResponse{}

	result := UserInfo.Table("users").Where("name = ?", username).First(&userregister)
	if result.Error != nil {
		failmsg := "Can't find user"
		msg = &failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.LoginResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	result2 := UserInfo.Table("users").Where("name = ?", username).Where("user_pass_word = ?", password).First(&userregister)
	if result2.Error != nil {
		failmsg := "Incorrect password"
		msg = &failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.LoginResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	token := ksuid.New().String()
	usertoken = &token
	failmsg := "Login Query Success"
	msg = &failmsg
	resp = &JsonStruct.LoginResponse{
		StatusCode: 0,
		StatusMsg:  &failmsg,
		Token:      usertoken,
		UserID:     &userregister.ID,
	}
	UserInfo.Table("users").Where("name = ?", username).Update("token", token)
	c.JSON(consts.StatusOK, resp)
} //此为登陆对应的具体服务实现

func Getinfo(ctx context.Context, c *app.RequestContext) {
	user_id, user_idBool := c.GetQuery("user_id")
	_, tokenBool := c.GetQuery("token")
	userid, _ := strconv.ParseInt(user_id, 10, 64)
	if !user_idBool || !tokenBool {
		var msg *string
		Failmsg := "no passed data"
		msg = &Failmsg
		c.JSON(consts.StatusUnauthorized, &JsonStruct.GetInfoResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		})
		return
	}
	user := ReadUser(userid)
	successInfo := "Get Imformation successfully"
	UserResp := &JsonStruct.RspUser{
		ID:              user.ID,
		Name:            user.Name,
		FollowCount:     user.FollowerCount,
		FollowerCount:   user.FollowCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
	InfoResp := &JsonStruct.GetInfoResponse{
		StatusCode: 0,
		StatusMsg:  &successInfo,
		User:       UserResp,
	}
	c.JSON(consts.StatusOK, InfoResp)
} //此为获取用户信息接口对应的服务实现
