package JsonStruct

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              int64 `gorm:"primaryKey" `
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Avatar          string         `json:"avatar"`                      // 用户头像
	BackgroundImage string         `json:"background_image"`            // 用户个人页顶部大图
	FavoriteCount   int64          `json:"favorite_count"`              // 喜欢数
	FollowCount     int64          `json:"follow_count"`                // 关注总数
	FollowerCount   int64          `json:"follower_count"`              // 粉丝总数
	IsFollow        bool           `json:"is_follow"`                   // true-已关注，false-未关注
	Name            string         `json:"name"`                        // 用户名称
	Signature       string         `json:"signature"`                   // 个人简介
	TotalFavorited  int64          `json:"total_favorited"`             // 获赞数量
	WorkCount       int64          `gorm:"default:0" json:"work_count"` // 作品数
	UserPassWord    string         `json:"user_pass_word"`
	Token           string         `json:"token"`
}

type Video struct {
	ID            int64 `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	UserID        int64
	CommentCount  int64  `json:"comment_count"`  // 视频的评论总数
	CoverURL      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"` // 视频的点赞总数
	IsFavorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string `json:"play_url"`       // 视频播放地址
	Title         string `json:"title"`          // 视频标题
}

type Comment struct {
	Content    string `json:"content"`              // 评论内容
	CreateDate string `json:"create_date"`          // 评论发布日期，格式 mm-dd
	ID         int64  `gorm:"primaryKey" json:"id"` //评论id
	UserID     int64  //用户id
	VideoID    int64
}
type FavoriteList struct {
	UserID    int64 `gorm:"primaryKey" json:"user_id"`
	VideoID   int64 `gorm:"primaryKey" json:"video_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
