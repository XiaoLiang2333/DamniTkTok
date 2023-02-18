package JsonStruct

import "gorm.io/gorm"

type UserRegister struct {
	UserName     string
	UserPassWord string
	Token        string
	ID           int64 `gorm:"primarykey"`
}
type Videolist struct {
	gorm.Model
	Author        User   `gorm:"foreignKey:VideoID;references:ID;" json:"author"` // 视频作者信息
	CommentCount  int64  `json:"comment_count"`                                   // 视频的评论总数
	CoverURL      string `json:"cover_url"`                                       // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"`                                  // 视频的点赞总数
	VideoID       int64  `gorm:"primarykey" json:"id"`                            // 视频唯一标识
	IsFavorite    bool   `json:"is_favorite"`                                     // true-已点赞，false-未点赞
	PlayURL       string `json:"play_url"`                                        // 视频播放地址
	Title         string `json:"title"`                                           // 视频标题
}
type User struct {
	Avatar          string `json:"avatar"`               // 用户头像
	BackgroundImage string `json:"background_image"`     // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`       // 喜欢数
	FollowCount     int64  `json:"follow_count"`         // 关注总数
	FollowerCount   int64  `json:"follower_count"`       // 粉丝总数
	ID              int64  `gorm:"primarykey" json:"id"` // 用户id
	IsFollow        bool   `json:"is_follow"`            // true-已关注，false-未关注
	Name            string `json:"name"`                 // 用户名称
	Signature       string `json:"signature"`            // 个人简介
	TotalFavorited  string `json:"total_favorited"`      // 获赞数量
	WorkCount       int64  `json:"work_count"`           // 作品数
}
