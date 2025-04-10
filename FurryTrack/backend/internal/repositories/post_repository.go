package repositories

import (
	"FurryTrack/internal/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	*BaseRepository
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create  - создает новый пост
func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// GetFeed - возвращает ленту постов
func (r *PostRepository) GetFeed(userID uuid.UUID, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.
		Preload("User").
		Preload("Pet").
		Order("created_at desc").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// GetByUserIDs - возвращает посты пользователей с пагинацией
func (r *PostRepository) GetByUserIDs(userIDs []uuid.UUID, offset, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Where("user_id IN ?", userIDs).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// AddLike - добавляет лайк к посту
func (r *PostRepository) AddLike(postID uuid.UUID, userID uuid.UUID) error {
	var exists bool
	err := r.db.Model(&models.PostLike{}).
		Select("count(*) > 0").
		Where("post_id = ? AND user_id = ?", postID, userID).
		Find(&exists).Error

	if err != nil {
		return err
	}
	if exists {
		return errors.New("already liked")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Добавляем запись о лайке
		if err := tx.Create(&models.PostLike{
			PostID: postID,
			UserID: userID,
		}).Error; err != nil {
			return err
		}

		// Обновляем счетчик лайков в посте
		return tx.Model(&models.Post{}).
			Where("id = ?", postID).
			Update("likes_count", gorm.Expr("likes_count + 1")).
			Error
	})
}

// GetBaseQuery - возвращает базовый запрос для постов с предзагрузкой связанных данных
func (r *PostRepository) GetBaseQuery() *gorm.DB {
	return r.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Pet")
}

// GetPosts выполняет запрос и возвращает список постов
func (r *PostRepository) GetPosts(query *gorm.DB) ([]models.Post, error) {
	var posts []models.Post
	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}


// CreateComment - создание комментария
func (r *PostRepository) CreateComment(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

// GetCommentsByPost - получение комментариев к посту
func (r *PostRepository) GetCommentsByPost(postID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.
		Preload("User").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}


// Обновляет счетчик лайков, возвращает количсетво лайков
func (r *PostRepository) IncrementLikes(postID uuid.UUID) (int, error) {
    err := r.db.Model(&models.Post{}).
        Where("id = ?", postID).
        Update("likes_count", gorm.Expr("likes_count + 1")).
        Error
    if err != nil {
        return 0, err
    }

    var post models.Post
    err = r.db.Select("likes_count").
        First(&post, "id = ?", postID).
        Error

    return post.LikesCount, err
}

// Получение поста по ID
func (r *PostRepository) GetByID(id uuid.UUID) (*models.Post, error) {
    var post models.Post
    err := r.db.
        Preload("Author"). 
        Preload("Pet").     
        First(&post, "id = ?", id).
        Error
    
    return &post, err
}
