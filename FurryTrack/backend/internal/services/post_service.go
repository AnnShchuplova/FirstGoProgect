package services

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
	//"time"
	"errors"
	"github.com/google/uuid"
)

type PostService struct {
	postRepo *repositories.PostRepository
	petRepo  *repositories.PetRepository
}

func NewPostService(postRepo *repositories.PostRepository, petRepo *repositories.PetRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		petRepo:  petRepo,
	}
}

// CreatePost создает новый пост
func (s *PostService) CreatePost(post *models.Post) error {
	// Проверяем, что питомец принадлежит пользователю
	if post.PetID != uuid.Nil {
		pet, err := s.petRepo.FindByID(post.PetID)
		if err != nil || pet.OwnerID != post.AuthorID {
			return errors.New("invalid pet")
		}
	}

	return s.postRepo.Create(post)
}

// GetFeed возвращает ленту постов
func (s *PostService) GetFeed(userID uuid.UUID, page, limit int) ([]models.Post, error) {
	offset := (page - 1) * limit
	return s.postRepo.GetByUserIDs([]uuid.UUID{userID}, offset, limit)
}



func (s *PostService) LikePost(userID, postID uuid.UUID) error {
	if _, err := s.postRepo.GetByID(postID); err != nil {
		return err
	}

	return s.postRepo.AddLike(postID, userID)
}

// AddComment - добавление комментария
func (s *PostService) AddComment(userID, postID uuid.UUID, content string) (*models.Comment, error) {
	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

	err := s.postRepo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetComments - получение комментариев
func (s *PostService) GetComments(postID uuid.UUID) ([]models.Comment, error) {
	return s.postRepo.GetCommentsByPost(postID)
}
