package services

import (
	"github.com/google/uuid"
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
)

type FeedService struct {
	postRepo     *repositories.PostRepository
	relationRepo *repositories.UserRelationRepository
}

func NewFeedService(
	postRepo *repositories.PostRepository,
	relationRepo *repositories.UserRelationRepository,
) *FeedService {
	return &FeedService{
		postRepo:     postRepo,
		relationRepo: relationRepo,
	}
}

// Обращения к репозиторию

func (s *FeedService) GetMainFeed(userID uuid.UUID) ([]models.Post, error) {
	return s.postRepo.GetPosts(
		s.postRepo.GetBaseQuery().
			Order("created_at DESC").
			Limit(100),
	)
}

func (s *FeedService) GetMarketFeed(userID uuid.UUID) ([]models.Post, error) {
	return s.postRepo.GetPosts(
		s.postRepo.GetBaseQuery().
			Where("post_type = ?", "market").
			Order("created_at DESC").
			Limit(100),
	)
}

func (s *FeedService) GetFollowingFeed(userID uuid.UUID) ([]models.Post, error) {
	following, err := s.relationRepo.GetFollowing(userID)
	if err != nil {
		return nil, err
	}

	if len(following) == 0 {
		return []models.Post{}, nil
	}

	followingIDs := make([]uuid.UUID, len(following))
	for i, rel := range following {
		followingIDs[i] = rel.FollowingID
	}

	return s.postRepo.GetPosts(
		s.postRepo.GetBaseQuery().
			Where("author_id IN ?", followingIDs).
			Order("created_at DESC").
			Limit(100),
	)
}