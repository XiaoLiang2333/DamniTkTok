package Service

import (
	"DamniTkTok/JsonStruct"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/segmentio/ksuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Register(ctx context.Context, c *app.RequestContext) {
	username, usernameBool := c.GetQuery("username")
	password, passwordBool := c.GetQuery("password")
	if !usernameBool || !passwordBool {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.RegisterReply{
			StatusCode: 1,
			StatusMsg:  "no passed data",
		})
		return
	}
	if len(username) > 32 || len(password) > 32 {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.RegisterReply{
			StatusCode: 1,
			StatusMsg:  "Input is too long(>32)",
		})
		return
	}
	dsn := "root:XHL13458218248.@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	resp := &JsonStruct.RegisterReply{}
	var userregister JsonStruct.UserRegister
	UserInfo.AutoMigrate(&JsonStruct.UserRegister{})
	result := UserInfo.Where("user_name = ?", username).First(&userregister)
	if result.Error == nil {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.RegisterReply{
			StatusCode: 1,
			StatusMsg:  "Already existed",
		})
		return
	}
	token := ksuid.New().String()
	userregister = JsonStruct.UserRegister{UserName: username, UserPassWord: password, Token: token}
	UserInfo.Create(&userregister)
	resp = &JsonStruct.RegisterReply{
		StatusCode: 0,
		StatusMsg:  "Register Query Success",
		Token:      token,
		UserID:     userregister.ID,
	}
	c.JSON(consts.StatusOK, resp)
}

/**此为注册接口对应的基础服务实现*/
func Login(ctx context.Context, c *app.RequestContext) {
	username, usernameBool := c.GetQuery("username")
	password, passwordBool := c.GetQuery("password")
	if !usernameBool || !passwordBool {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.LoginReply{
			StatusCode: 1,
			StatusMsg:  "no passed data",
		})
		return
	}
	dsn := "root:XHL13458218248.@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	UserInfo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	resp := &JsonStruct.LoginReply{}
	var userregister JsonStruct.UserRegister
	result := UserInfo.Where("user_name = ?", username).First(&userregister)
	if result.Error != nil {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.LoginReply{
			StatusCode: 1,
			StatusMsg:  "Found no user",
		})
		return
	}
	result2 := UserInfo.Where("user_name = ?", username).Where("user_pass_word = ?", password).First(&userregister)
	if result2.Error != nil {
		c.JSON(consts.StatusUnauthorized, &JsonStruct.LoginReply{
			StatusCode: 1,
			StatusMsg:  "password incorrect",
		})
		return
	}

	resp = &JsonStruct.LoginReply{
		StatusCode: 0,
		StatusMsg:  "Login Query Success",
		Token:      userregister.Token,
		UserID:     userregister.ID,
	}
	c.JSON(consts.StatusOK, resp)
}

/**此为登陆对应的具体服务实现*/
