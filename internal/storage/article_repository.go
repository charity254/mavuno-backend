package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) CreateArticle(a *models.Article) error {
	query := `
		INSERT INTO articles (id, author_id, title, content, category, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query,
		a.ID, a.AuthorID, a.Title,
		a.Content, a.Category, a.Deleted,
		a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}

func (r *ArticleRepository) GetArticleByID(id uuid.UUID) (*models.Article, error) {
	query := `
		SELECT id, author_id, title, content, category, deleted, created_at, updated_at
		FROM articles
		WHERE id = $1 AND deleted = false
	`
	a := &models.Article{}
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.AuthorID, &a.Title,
		&a.Content, &a.Category, &a.Deleted,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("article not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return a, nil
}

func (r *ArticleRepository) GetAllArticles(category string) ([]*models.Article, error) {
	query := `
		SELECT id, author_id, title, content, category, deleted, created_at, updated_at
		FROM articles
		WHERE deleted = false
	`
	args := []interface{}{}

	// Add category filter if provided
	if category != "" {
		query += " AND category = $1"
		args = append(args, category)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		a := &models.Article{}
		err := rows.Scan(
			&a.ID, &a.AuthorID, &a.Title,
			&a.Content, &a.Category, &a.Deleted,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, a)
	}
	return articles, nil
}

func (r *ArticleRepository) SoftDeleteArticle(id, authorID uuid.UUID) error {
	query := `
		UPDATE articles
		SET deleted = true, updated_at = NOW()
		WHERE id = $1 AND author_id = $2 AND deleted = false
	`
	result, err := r.db.Exec(query, id, authorID)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("article not found or already deleted")
	}
	return nil
}