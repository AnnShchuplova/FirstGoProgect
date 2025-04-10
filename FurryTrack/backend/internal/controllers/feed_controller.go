package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"FurryTrack/internal/services"
)

type FeedController struct {
	feedService    *services.FeedService
	relationService *services.UserRelationService
}

func NewFeedController(
	feedService *services.FeedService,
	relationService *services.UserRelationService,
) *FeedController {
	return &FeedController{
		feedService:    feedService,
		relationService: relationService,
	}
}

// GetMainFeed - основная лента (все посты)
func (c *FeedController) GetMainFeed(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uuid.UUID)

	posts, err := c.feedService.GetMainFeed(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": posts,
		"meta": gin.H{
			"count": len(posts),
			"type":  "regular",
		},
	})
}

// GetMarketFeed - лента продаж (посты с типом market)
func (c *FeedController) GetMarketFeed(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uuid.UUID)

	posts, err := c.feedService.GetMarketFeed(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": posts,
		"meta": gin.H{
			"count": len(posts),
			"type":  "market",
		},
	})
}

// GetFollowingFeed - лента подписок (посты тех, на кого подписан пользователь)
func (c *FeedController) GetFollowingFeed(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uuid.UUID)

	posts, err := c.feedService.GetFollowingFeed(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": posts,
		"meta": gin.H{
			"count": len(posts),
			"type":  "following",
		},
	})
}