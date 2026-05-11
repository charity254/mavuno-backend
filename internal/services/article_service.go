package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

type ArticleService struct {
	repo *storage.ArticleRepository
}

func NewArticleService(repo *storage.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) CreateArticle(authorID uuid.UUID, title, content, category string) (*models.Article, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}

	a := &models.Article{
		ID:        uuid.New(),
		AuthorID:  authorID,
		Title:     title,
		Content:   content,
		Category:  category,
		Deleted:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateArticle(a); err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return a, nil
}

func (s *ArticleService) GetArticles(category string) ([]*models.Article, error) {
	articles, err := s.repo.GetAllArticles(category)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}
	if articles == nil {
		articles = []*models.Article{}
	}
	return articles, nil
}

func (s *ArticleService) GetArticleByID(id uuid.UUID) (*models.Article, error) {
	article, err := s.repo.GetArticleByID(id)
	if err != nil {
		return nil, fmt.Errorf("article not found")
	}
	return article, nil
}

func (s *ArticleService) DeleteArticle(id, authorID uuid.UUID) error {
	article, err := s.repo.GetArticleByID(id)
	if err != nil {
		return fmt.Errorf("article not found")
	}

	if article.AuthorID != authorID {
		return fmt.Errorf("you do not have permission to delete this article")
	}

	return s.repo.SoftDeleteArticle(id, authorID)
}