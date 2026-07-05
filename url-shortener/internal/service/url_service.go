package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/manasdixit/url-shortener/internal/domain"
	"github.com/manasdixit/url-shortener/internal/repository"
	"github.com/manasdixit/url-shortener/internal/utils"
)

var (
	ErrInvalidURL       = errors.New("invalid original url")
	ErrAliasUnavailable = errors.New("custom alias is not available")
	ErrURLNotFound      = errors.New("url not found")
	ErrURLInactive      = errors.New("url is inactive")
	ErrURLExpired       = errors.New("url has expired")
)

type URLService struct {
	urlRepo         repository.URLRepository
	clickRepo       repository.ClickRepository
	baseURL         string
	shortCodeLength int
}

func NewURLService(
	urlRepo repository.URLRepository,
	clickRepo repository.ClickRepository,
	baseURL string,
	shortCodeLength int,
) *URLService {
	return &URLService{
		urlRepo:         urlRepo,
		clickRepo:       clickRepo,
		baseURL:         strings.TrimRight(baseURL, "/"),
		shortCodeLength: shortCodeLength,
	}
}

func (s *URLService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateURLRequest) (domain.URLResponse, error) {
	req.OriginalURL = strings.TrimSpace(req.OriginalURL)
	if !utils.IsValidHTTPURL(req.OriginalURL) {
		return domain.URLResponse{}, ErrInvalidURL
	}

	var shortCode string
	var customAlias *string

	if req.CustomAlias != nil && strings.TrimSpace(*req.CustomAlias) != "" {
		alias := strings.TrimSpace(*req.CustomAlias)
		if !utils.IsValidAlias(alias) {
			return domain.URLResponse{}, ErrInvalidInput
		}

		exists, err := s.urlRepo.CustomAliasExists(ctx, alias)
		if err != nil {
			return domain.URLResponse{}, err
		}
		if exists {
			return domain.URLResponse{}, ErrAliasUnavailable
		}

		shortCode = alias
		customAlias = &alias
	} else {
		code, err := s.generateUniqueShortCode(ctx)
		if err != nil {
			return domain.URLResponse{}, err
		}
		shortCode = code
	}

	url := domain.ShortURL{
		UserID:      userID,
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
		CustomAlias: customAlias,
		Title:       req.Title,
		IsActive:    true,
		ExpiresAt:   req.ExpiresAt,
	}

	createdURL, err := s.urlRepo.Create(ctx, url)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.URLResponse{}, ErrAliasUnavailable
		}
		return domain.URLResponse{}, err
	}

	return domain.NewURLResponse(createdURL, s.baseURL), nil
}

func (s *URLService) List(ctx context.Context, userID uuid.UUID) ([]domain.URLResponse, error) {
	urls, err := s.urlRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.URLResponse, 0, len(urls))
	for _, url := range urls {
		responses = append(responses, domain.NewURLResponse(url, s.baseURL))
	}
	return responses, nil
}

func (s *URLService) GetByID(ctx context.Context, userID, id uuid.UUID) (domain.URLResponse, error) {
	url, err := s.urlRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.URLResponse{}, ErrURLNotFound
		}
		return domain.URLResponse{}, err
	}
	return domain.NewURLResponse(url, s.baseURL), nil
}

func (s *URLService) Update(ctx context.Context, userID, id uuid.UUID, req domain.UpdateURLRequest) (domain.URLResponse, error) {
	url, err := s.urlRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.URLResponse{}, ErrURLNotFound
		}
		return domain.URLResponse{}, err
	}

	if req.OriginalURL != nil {
		originalURL := strings.TrimSpace(*req.OriginalURL)
		if !utils.IsValidHTTPURL(originalURL) {
			return domain.URLResponse{}, ErrInvalidURL
		}
		url.OriginalURL = originalURL
	}

	if req.Title != nil {
		url.Title = req.Title
	}

	if req.IsActive != nil {
		url.IsActive = *req.IsActive
	}

	if req.ExpiresAt != nil {
		url.ExpiresAt = req.ExpiresAt
	}

	updatedURL, err := s.urlRepo.Update(ctx, url)
	if err != nil {
		return domain.URLResponse{}, err
	}

	return domain.NewURLResponse(updatedURL, s.baseURL), nil
}

func (s *URLService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	if err := s.urlRepo.Deactivate(ctx, id, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrURLNotFound
		}
		return err
	}
	return nil
}

func (s *URLService) ResolveRedirect(ctx context.Context, shortCode string) (domain.ShortURL, error) {
	url, err := s.urlRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.ShortURL{}, ErrURLNotFound
		}
		return domain.ShortURL{}, err
	}

	if !url.IsActive {
		return domain.ShortURL{}, ErrURLInactive
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return domain.ShortURL{}, ErrURLExpired
	}

	return url, nil
}

func (s *URLService) TrackClick(ctx context.Context, click domain.Click) error {
	return s.clickRepo.Create(ctx, click)
}

func (s *URLService) Analytics(ctx context.Context, userID, urlID uuid.UUID) (domain.AnalyticsResponse, error) {
	url, err := s.urlRepo.FindByIDAndUserID(ctx, urlID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.AnalyticsResponse{}, ErrURLNotFound
		}
		return domain.AnalyticsResponse{}, err
	}

	total, err := s.clickRepo.CountByURLID(ctx, url.ID)
	if err != nil {
		return domain.AnalyticsResponse{}, err
	}

	daily, err := s.clickRepo.DailyStatsByURLID(ctx, url.ID, 30)
	if err != nil {
		return domain.AnalyticsResponse{}, err
	}

	return domain.AnalyticsResponse{
		URLID:       url.ID,
		TotalClicks: total,
		DailyClicks: daily,
	}, nil
}

func (s *URLService) generateUniqueShortCode(ctx context.Context) (string, error) {
	for i := 0; i < 5; i++ {
		code, err := utils.GenerateShortCode(s.shortCodeLength)
		if err != nil {
			return "", err
		}

		exists, err := s.urlRepo.ShortCodeExists(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique short code")
}
