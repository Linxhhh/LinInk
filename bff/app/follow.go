package app

import (
	"strconv"

	"github.com/Linxhhh/LinInk/api/proto/follow"
	"github.com/Linxhhh/LinInk/pkg/jwts"
	"github.com/Linxhhh/LinInk/pkg/res"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FollowHandler struct {
	svc follow.FollowServiceClient
}

func NewFollowHandler(svc follow.FollowServiceClient) *FollowHandler {
	return &FollowHandler{
		svc: svc,
	}
}

func (hdl *FollowHandler) RegistryRouter(router *gin.Engine) {
	ur := router.Group("userRelation")
	ur.POST("follow", hdl.Follow)
	ur.GET("follow", hdl.FollowData)
	ur.GET("followees", hdl.FolloweeList)
	ur.GET("followers", hdl.FollowerList)
}

func (hdl *FollowHandler) Follow(ctx *gin.Context) {

	// 绑定参数
	type Req struct {
		Id     int64 `json:"id"`     // id 表示被关注的用户 id
		Follow bool  `json:"follow"` // true 表示关注，false 表示取消
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var err error
	if req.Follow {
		// err = hdl.svc.Follow(ctx, claims.UserId, req.Id)
		_, err = hdl.svc.Follow(ctx.Request.Context(), &follow.FollowRequest{
			FollowerId: claims.UserId,
			FolloweeId: req.Id,
		})
	} else {
		// err = hdl.svc.CancelFollow(ctx, claims.UserId, req.Id)
		_, err = hdl.svc.CancelFollow(ctx.Request.Context(), &follow.CancelFollowRequest{
			FollowerId: claims.UserId,
			FolloweeId: req.Id,
		})
	}
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("操作成功", ctx)
}

func (hdl *FollowHandler) FollowData(ctx *gin.Context) {

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	uid := claims.UserId

	// 绑定参数
	id := ctx.Query("id")
	if id != "" {
		uid, _ = strconv.ParseInt(id, 10, 64)
	}

	// 获取用户关系数据
	// data, err := hdl.svc.GetFollowData(ctx, uid)
	resp, err := hdl.svc.GetFollowData(ctx.Request.Context(), &follow.GetFollowDataRequest{Uid: uid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			res.OKWithData(resp.GetFollowData(), ctx)
			return
		}
		res.FailWithMsg("系统错误", ctx)
		return
	}

	if uid != claims.UserId {
		// data.IsFollowed, err = hdl.svc.GetFollowed(ctx, claims.UserId, uid)
		followedResp, err := hdl.svc.GetFollowed(ctx.Request.Context(), &follow.GetFollowedRequest{
			FollowerId: claims.UserId,
			FolloweeId: uid,
		})
		resp.FollowData.IsFollowed = followedResp.GetIsFollowed()
		if err != nil {
			res.FailWithMsg("系统错误", ctx)
			return
		}
	}

	// 返回响应
	type Resp struct {
		Followers  int64 `json:"followers"`  // 粉丝数量
		Followees  int64 `json:"followees"`  // 关注数量
		IsFollowed bool  `json:"isFollowed"` // 是否已关注
	}
	res.OKWithData(Resp{
		Followers:  resp.GetFollowData().Followers,
		Followees:  resp.GetFollowData().Followees,
		IsFollowed: resp.GetFollowData().IsFollowed,
	}, ctx)
}

// FolloweeList 获取关注列表
func (hdl *FollowHandler) FolloweeList(ctx *gin.Context) {

	// 绑定参数
	type ListReq struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}
	var req ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	// articleList, err := hdl.svc.GetFolloweeList(ctx, claims.UserId, req.Page, req.PageSize)
	resp, err := hdl.svc.GetFolloweeList(ctx.Request.Context(), &follow.GetFolloweeListRequest{
		FollowerId: claims.UserId,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 返回响应
	res.OKWithData(resp.GetFolloweeList(), ctx)
}

// FollowerList 获取粉丝列表
func (hdl *FollowHandler) FollowerList(ctx *gin.Context) {

	// 绑定参数
	type ListReq struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}
	var req ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	// articleList, err := hdl.svc.GetFollowerList(ctx, claims.UserId, req.Page, req.PageSize)
	resp, err := hdl.svc.GetFollowerList(ctx.Request.Context(), &follow.GetFollowerListRequest{
		FolloweeId: claims.UserId,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 返回响应
	res.OKWithData(resp.GetFollowerList(), ctx)
}
