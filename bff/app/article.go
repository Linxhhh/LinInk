package app

import (
	"log"
	"strconv"
	"time"

	"github.com/Linxhhh/LinInk/api/proto/article"
	"github.com/Linxhhh/LinInk/api/proto/interaction"
	"github.com/Linxhhh/LinInk/pkg/jwts"
	"github.com/Linxhhh/LinInk/pkg/res"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ArticleHandler struct {
	svc      article.ArticleServiceClient
	interSvc interaction.InteractionServiceClient
	biz      string
}

func NewArticleHandler(svc article.ArticleServiceClient, interSvc interaction.InteractionServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:      svc,
		interSvc: interSvc,
		biz:      "article",
	}
}

func (hdl *ArticleHandler) RegistryRouter(router *gin.Engine) {
	// 作者接口
	ag := router.Group("article")
	ag.POST("edit", hdl.Edit)
	ag.POST("publish", hdl.Publish)
	ag.DELETE("withdraw", hdl.Withdraw)
	ag.GET("count", hdl.Count)
	ag.GET("list", hdl.List)
	ag.GET("detail", hdl.Detail)

	// 读者接口
	pg := router.Group("pub")
	pg.GET("list", hdl.PubList)
	pg.GET("detail", hdl.PubDetail, hdl.Read)
	pg.POST("like", hdl.Like)
	pg.POST("collect", hdl.Collect)
	pg.GET("interaction", hdl.Interaction)
}

// Edit 新建帖子，或编辑旧帖子
func (hdl *ArticleHandler) Edit(ctx *gin.Context) {

	// 绑定参数
	var req ArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	/*
		aid, err := hdl.svc.Save(ctx, domain.Article{
			Id:       req.Id,
			Title:    req.Title,
			Content:  req.Content,
			AuthorId: claims.UserId,
		})
	*/
	resp, err := hdl.svc.Save(ctx, &article.SaveRequest{
		Article: &article.Article{
			Id:       req.Id,
			Title:    req.Title,
			Content:  req.Content,
			AuthorId: claims.UserId,
		},
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithData(gin.H{"article_id": resp.GetAid()}, ctx)
}

// Pubish 帖子发表
func (hdl *ArticleHandler) Publish(ctx *gin.Context) {

	// 绑定参数
	var req ArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	/*
		aid, err := hdl.svc.Publish(ctx, domain.Article{
			Id:       req.Id,
			Title:    req.Title,
			Content:  req.Content,
			AuthorId: claims.UserId,
		})
	*/
	resp, err := hdl.svc.Publish(ctx, &article.PublishRequest{
		Article: &article.Article{
			Id:       req.Id,
			Title:    req.Title,
			Content:  req.Content,
			AuthorId: claims.UserId,
		},
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	/*
		// 异步事件 —— feed 流推送
		if err = hdl.producer.ProduceEvent(events.ArticleEvent{
			Uid:   claims.UserId,
			Aid:   resp.GetAid(),
			Title: req.Title,
		}); err != nil {
			res.FailWithMsg("异步事件生成错误", ctx)
		}
	*/
	res.OKWithData(gin.H{"article_id": resp.GetAid()}, ctx)
}

// Withdraw 撤销发表
func (hdl *ArticleHandler) Withdraw(ctx *gin.Context) {

	// 绑定参数
	type Req struct{ Id int64 }
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	/*
		err := hdl.svc.Withdraw(ctx, claims.UserId, req.Id)
	*/
	_, err := hdl.svc.Withdraw(ctx, &article.WithdrawRequest{
		Uid: claims.UserId,
		Aid: req.Id,
	})
	if err != nil {
		res.FailWithMsg("撤销失败", ctx)
		return
	}
	res.OKWithMsg("撤销成功", ctx)
}

type ArticleRequest struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Count 获取用户制作库的帖子总数
func (hdl *ArticleHandler) Count(ctx *gin.Context) {

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	/*
		count, err := hdl.svc.Count(ctx, claims.UserId)
	*/
	resp, err := hdl.svc.Count(ctx, &article.CountRequest{
		Uid: claims.UserId,
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 返回响应
	res.OKWithData(gin.H{"total": resp.GetCount()}, ctx)
}

// List 获取用户制作库的帖子列表
func (hdl *ArticleHandler) List(ctx *gin.Context) {

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
	/*
		articleList, err := hdl.svc.List(ctx, claims.UserId, req.Page, req.PageSize)
	*/
	resp, err := hdl.svc.List(ctx, &article.ListRequest{
		Uid:      claims.UserId,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 返回响应
	res.OKWithData(resp.GetList(), ctx)
}

// Detail 获取制作库的帖子详情
func (hdl *ArticleHandler) Detail(ctx *gin.Context) {

	// 绑定参数
	aid, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if aid == 0 || err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	/*
		art, err := hdl.svc.Detail(ctx, claims.UserId, aid)
	*/
	resp, err := hdl.svc.Detail(ctx, &article.DetailRequest{
		Uid: claims.UserId,
		Aid: aid,
	})
	if err != nil {
		res.FailWithMsg("查询帖子失败", ctx)
		return
	}

	// 返回响应
	res.OKWithData(resp.GetArticle(), ctx)
}

// PubDetail 获取线上库的帖子详情
func (hdl *ArticleHandler) PubDetail(ctx *gin.Context) {

	// 绑定参数
	aid, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if aid == 0 || err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 调用下层服务
	/*
		art, err := hdl.svc.PubDetail(ctx, aid)
	*/
	resp, err := hdl.svc.PubDetail(ctx, &article.PubDetailRequest{Aid: aid})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 设置上下文，然后返回响应
	ctx.Set("article_id", aid)
	res.OKWithData(resp.GetArticle(), ctx)
}

/*
阅读、点赞、收藏功能待优化
*/

func (hdl *ArticleHandler) Read(ctx *gin.Context) {

	aid, exists := ctx.Get("article_id")
	if !exists {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 调用下层服务
	/*
		err := hdl.interSvc.IncrReadCnt(ctx, hdl.biz, aid.(int64))
	*/
	_, err := hdl.interSvc.IncrReadCnt(ctx, &interaction.IncrReadCntRequest{
		Biz:   hdl.biz,
		BizId: aid.(int64),
	})
	if err != nil {
		log.Panicln("IncrReadCnt 报错：err : ", err.Error())
	}
}

func (hdl *ArticleHandler) Like(ctx *gin.Context) {

	// 绑定参数
	type Req struct {
		Id   int64 `json:"id"`
		Like bool  `json:"like"` // true 表示点赞，false 表示取消
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
	if req.Like {
		// err = hdl.interSvc.Like(ctx, hdl.biz, req.Id, claims.UserId)
		_, err = hdl.interSvc.Like(ctx, &interaction.LikeRequest{
			Biz:   hdl.biz,
			BizId: req.Id,
			Uid:   claims.UserId,
		})
	} else {
		// err = hdl.interSvc.CancelLike(ctx, hdl.biz, req.Id, claims.UserId)
		_, err = hdl.interSvc.CancelLike(ctx, &interaction.CancelLikeRequest{
			Biz:   hdl.biz,
			BizId: req.Id,
			Uid:   claims.UserId,
		})
	}
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("操作成功", ctx)
}

func (hdl *ArticleHandler) Collect(ctx *gin.Context) {
	// 绑定参数
	type Req struct {
		Id      int64 `json:"id"`
		Collect bool  `json:"collect"` // true 表示点赞，false 表示取消
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
	if req.Collect {
		// err = hdl.interSvc.Collect(ctx, hdl.biz, req.Id, claims.UserId)
		_, err = hdl.interSvc.Collect(ctx, &interaction.CollectRequest{
			Biz:   hdl.biz,
			BizId: req.Id,
			Uid:   claims.UserId,
		})
	} else {
		// err = hdl.interSvc.CancelCollect(ctx, hdl.biz, req.Id, claims.UserId)
		_, err = hdl.interSvc.CancelCollect(ctx, &interaction.CancelCollectRequest{
			Biz:   hdl.biz,
			BizId: req.Id,
			Uid:   claims.UserId,
		})
	}
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("操作成功", ctx)
}

func (hdl *ArticleHandler) CollectionList(ctx *gin.Context) {

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	// arts, err := hdl.interSvc.CollectionList(ctx, hdl.biz, claims.UserId)
	resp, err := hdl.svc.CollectionList(ctx, &article.CollectionListRequest{
		Uid: claims.UserId,
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithData(resp.GetList(), ctx)
}

func (hdl *ArticleHandler) Interaction(ctx *gin.Context) {

	// 绑定参数
	aid, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if aid == 0 || err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	// i, err := hdl.interSvc.Get(ctx, hdl.biz, aid, claims.UserId)
	resp, err := hdl.interSvc.Get(ctx, &interaction.GetRequest{
		Biz:   hdl.biz,
		BizId: aid,
		Uid:   claims.UserId,
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 返回响应
	type Resp struct {
		ReadCnt     int64 `json:"readCnt"`
		LikeCnt     int64 `json:"likeCnt"`
		CollectCnt  int64 `json:"collectCnt"`
		IsLiked     bool  `json:"isLiked"`
		IsCollected bool  `json:"isCollected"`
	}
	res.OKWithData(Resp{
		ReadCnt:     resp.GetInteraction().ReadCnt,
		LikeCnt:     resp.GetInteraction().LikeCnt,
		CollectCnt:  resp.GetInteraction().CollectCnt,
		IsLiked:     resp.GetInteraction().IsLiked,
		IsCollected: resp.GetInteraction().IsCollected,
	}, ctx)
}

func (hdl *ArticleHandler) PubList(ctx *gin.Context) {

	// limit and offset
	type Req struct {
		Time   time.Time `json:"time"`
		Limit  int       `json:"limit"`
		Offset int       `json:"offset"`
	}
	var req Req
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 调用下层服务
	/*
		list, err := hdl.svc.PubList(ctx, req.Time, req.Limit, req.Offset)
	*/
	resp, err := hdl.svc.PubList(ctx, &article.PubListRequest{
		Timestamp: timestamppb.New(req.Time),
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
	})
	if err != nil {
		res.FailWithMsg("获取帖子失败", ctx)
		return
	}

	list := resp.GetList()
	if len(list) == 0 {
		res.OKWithMsg("目前没有新帖子", ctx)
		return
	}
	res.OKWithData(list, ctx)
}
