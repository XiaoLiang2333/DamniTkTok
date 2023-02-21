package JsonStruct

type RegisterResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}
type LoginResponse struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	Token      *string `json:"token"`       // 用户鉴权token
	UserID     *int64  `json:"user_id"`     // 用户id
}
type GetInfoResponse struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	User       *User   `json:"user"`        // 用户信息
}
type PublishRsp struct {
	StatusCode int64   `json:"status_code"`
	StatusMsg  *string `json:"status_msg"`
}
type ListRsp struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户发布的视频列表
}
type PlayVideoRsp struct {
	StatusCode int64
	StatusMsg  string
}
type CommentRsp struct {
	Comment    *CommentRes `json:"comment"`     // 评论成功返回评论内容，不需要重新拉取整个列表
	StatusCode int64       `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string     `json:"status_msg"`  // 返回状态描述
}
type CommentListRsp struct {
	CommentList []Comment `json:"comment_list"` // 评论列表
	StatusCode  int64     `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   *string   `json:"status_msg"`   // 返回状态描述
}

type CommentUserRes struct {
	ID              int64  `json:"id"`                          //用户id
	Name            string `json:"name"`                        // 用户名称
	FollowCount     int64  `json:"follow_count"`                // 关注总数
	FollowerCount   int64  `json:"follower_count"`              // 粉丝总数
	IsFollow        bool   `json:"is_follow"`                   // true-已关注，false-未关注
	Avatar          string `json:"avatar"`                      // 用户头像
	BackgroundImage string `json:"background_image"`            // 用户个人页顶部大图
	Signature       string `json:"signature"`                   // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`             // 获赞数量
	WorkCount       int64  `gorm:"default:0" json:"work_count"` // 作品数
	FavoriteCount   int64  `json:"favorite_count"`              // 喜欢数

}

type CommentRes struct {
	Content    string         `json:"content"`     // 评论内容
	CreateDate string         `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64          `json:"id"`          //评论id
	User       CommentUserRes `json:"user"`        // 评论用户信息
}
