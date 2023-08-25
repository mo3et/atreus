package biz

import (
	"context"

	"github.com/toomanysource/atreus/pkg/common"

	"github.com/toomanysource/atreus/app/feed/service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type Video struct {
	Id            uint32 `json:"id"`
	Author        User   `json:"author"`
	CommentCount  uint32 `json:"comment_count"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount uint32 `json:"favorite_count"`
	IsFavorite    bool   `json:"is_favorite"`
	PlayUrl       string `json:"play_url"`
	Title         string `json:"title"`
}

type User struct {
	Id              uint32 `json:"id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	FavoriteCount   uint32 `json:"favorite_count"`
	FollowCount     uint32 `json:"follow_count"`
	FollowerCount   uint32 `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Signature       string `json:"signature"`
	TotalFavorited  uint32 `json:"total_favorited"`
	WorkCount       uint32 `json:"work_count"`
}

type FeedRepo interface {
	GetFeedList(context.Context, string, uint32) (int64, []*Video, error)
}

type PublishRepo interface {
	GetVideoList(ctx context.Context, latestTime string, userId uint32, number uint32) (int64, []*Video, error)
}

type FeedUsecase struct {
	repo   FeedRepo
	config *conf.JWT
	log    *log.Helper
}

func NewFeedUsecase(repo FeedRepo, conf *conf.JWT, logger log.Logger) *FeedUsecase {
	return &FeedUsecase{
		repo: repo, config: conf,
		log: log.NewHelper(log.With(logger, "model", "usecase/feed")),
	}
}

func (uc *FeedUsecase) FeedList(
	ctx context.Context, latestTime string, tokenString string,
) (int64, []*Video, error) {
	if tokenString == "" {
		return uc.repo.GetFeedList(ctx, latestTime, 0)
	}
	token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return 0, nil, err
	}
	data, err := common.GetTokenData(token)
	if err != nil {
		return 0, nil, err
	}
	userId := uint32(data["user_id"].(float64))
	return uc.repo.GetFeedList(ctx, latestTime, userId)
}
