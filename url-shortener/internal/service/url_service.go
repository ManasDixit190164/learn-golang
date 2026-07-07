package service // package declaration for the module

import ( // start import block
	"context" // import package
	"errors" // import package
	"strings" // import package
	"time" // import package

	"github.com/google/uuid" // import package
	"github.com/manasdixit/url-shortener/internal/domain" // import package
	"github.com/manasdixit/url-shortener/internal/repository" // import package
	"github.com/manasdixit/url-shortener/internal/utils" // import package
) // end import block or block scope

var ( // declare variable
	ErrInvalidURL       = errors.New("invalid original url") // assign value
	ErrAliasUnavailable = errors.New("custom alias is not available") // assign value
	ErrURLNotFound      = errors.New("url not found") // assign value
	ErrURLInactive      = errors.New("url is inactive") // assign value
	ErrURLExpired       = errors.New("url has expired") // assign value
) // end import block or block scope

type URLService struct { // declare struct type
	urlRepo         repository.URLRepository // execute statement
	clickRepo       repository.ClickRepository // execute statement
	baseURL         string // execute statement
	shortCodeLength int // execute statement
} // end block

func NewURLService( // declare function
	urlRepo repository.URLRepository, // execute statement
	clickRepo repository.ClickRepository, // execute statement
	baseURL string, // execute statement
	shortCodeLength int, // execute statement
) *URLService { // execute statement
	return &URLService{ // return statement
		urlRepo:         urlRepo, // execute statement
		clickRepo:       clickRepo, // execute statement
		baseURL:         strings.TrimRight(baseURL, "/"), // execute statement
		shortCodeLength: shortCodeLength, // execute statement
	} // end block
} // end block

func (s *URLService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateURLRequest) (domain.URLResponse, error) { // declare method
	req.OriginalURL = strings.TrimSpace(req.OriginalURL) // trim whitespace
	if !utils.IsValidHTTPURL(req.OriginalURL) { // validate that a URL is http or https
		return domain.URLResponse{}, ErrInvalidURL // return statement
	} // end block

	var shortCode string // execute statement
	var customAlias *string // execute statement

	if req.CustomAlias != nil && strings.TrimSpace(*req.CustomAlias) != "" { // trim whitespace
		alias := strings.TrimSpace(*req.CustomAlias) // trim whitespace
		if !utils.IsValidAlias(alias) { // if condition
			return domain.URLResponse{}, ErrInvalidInput // return statement
		} // end block

		exists, err := s.urlRepo.CustomAliasExists(ctx, alias) // check custom alias availability
		if err != nil { // if condition
			return domain.URLResponse{}, err // return statement
		} // end block
		if exists { // if condition
			return domain.URLResponse{}, ErrAliasUnavailable // return statement
		} // end block

		shortCode = alias // assign value
		customAlias = &alias // assign value
	} else { // close block
		code, err := s.generateUniqueShortCode(ctx) // declare and initialize variable
		if err != nil { // if condition
			return domain.URLResponse{}, err // return statement
		} // end block
		shortCode = code // assign value
	} // end block

	url := domain.ShortURL{ // declare and initialize variable
		UserID:      userID, // execute statement
		OriginalURL: req.OriginalURL, // execute statement
		ShortCode:   shortCode, // execute statement
		CustomAlias: customAlias, // execute statement
		Title:       req.Title, // execute statement
		IsActive:    true, // execute statement
		ExpiresAt:   req.ExpiresAt, // execute statement
	} // end block

	createdURL, err := s.urlRepo.Create(ctx, url) // save a new URL record
	if err != nil { // if condition
		if errors.Is(err, repository.ErrConflict) { // handle duplicate database entry
			return domain.URLResponse{}, ErrAliasUnavailable // return statement
		} // end block
		return domain.URLResponse{}, err // return statement
	} // end block

	return domain.NewURLResponse(createdURL, s.baseURL), nil // return statement
} // end block

func (s *URLService) List(ctx context.Context, userID uuid.UUID) ([]domain.URLResponse, error) { // declare method
	urls, err := s.urlRepo.ListByUserID(ctx, userID) // retrieve a user’s URLs
	if err != nil { // if condition
		return nil, err // return statement
	} // end block

	responses := make([]domain.URLResponse, 0, len(urls)) // declare and initialize variable
	for _, url := range urls { // for loop
		responses = append(responses, domain.NewURLResponse(url, s.baseURL)) // assign value
	} // end block
	return responses, nil // return statement
} // end block

func (s *URLService) GetByID(ctx context.Context, userID, id uuid.UUID) (domain.URLResponse, error) { // declare method
	url, err := s.urlRepo.FindByIDAndUserID(ctx, id, userID) // get URL by ID and owner
	if err != nil { // if condition
		if errors.Is(err, repository.ErrNotFound) { // handle missing database record
			return domain.URLResponse{}, ErrURLNotFound // return statement
		} // end block
		return domain.URLResponse{}, err // return statement
	} // end block
	return domain.NewURLResponse(url, s.baseURL), nil // return statement
} // end block

func (s *URLService) Update(ctx context.Context, userID, id uuid.UUID, req domain.UpdateURLRequest) (domain.URLResponse, error) { // declare method
	url, err := s.urlRepo.FindByIDAndUserID(ctx, id, userID) // get URL by ID and owner
	if err != nil { // if condition
		if errors.Is(err, repository.ErrNotFound) { // handle missing database record
			return domain.URLResponse{}, ErrURLNotFound // return statement
		} // end block
		return domain.URLResponse{}, err // return statement
	} // end block

	if req.OriginalURL != nil { // if condition
		originalURL := strings.TrimSpace(*req.OriginalURL) // trim whitespace
		if !utils.IsValidHTTPURL(originalURL) { // validate that a URL is http or https
			return domain.URLResponse{}, ErrInvalidURL // return statement
		} // end block
		url.OriginalURL = originalURL // assign value
	} // end block

	if req.Title != nil { // if condition
		url.Title = req.Title // assign value
	} // end block

	if req.IsActive != nil { // if condition
		url.IsActive = *req.IsActive // assign value
	} // end block

	if req.ExpiresAt != nil { // if condition
		url.ExpiresAt = req.ExpiresAt // assign value
	} // end block

	updatedURL, err := s.urlRepo.Update(ctx, url) // update URL metadata
	if err != nil { // if condition
		return domain.URLResponse{}, err // return statement
	} // end block

	return domain.NewURLResponse(updatedURL, s.baseURL), nil // return statement
} // end block

func (s *URLService) Delete(ctx context.Context, userID, id uuid.UUID) error { // declare method
	if err := s.urlRepo.Deactivate(ctx, id, userID); err != nil { // deactivate a URL
		if errors.Is(err, repository.ErrNotFound) { // handle missing database record
			return ErrURLNotFound // return statement
		} // end block
		return err // return statement
	} // end block
	return nil // return statement
} // end block

func (s *URLService) ResolveRedirect(ctx context.Context, shortCode string) (domain.ShortURL, error) { // declare method
	url, err := s.urlRepo.FindByShortCode(ctx, shortCode) // get URL by short code
	if err != nil { // if condition
		if errors.Is(err, repository.ErrNotFound) { // handle missing database record
			return domain.ShortURL{}, ErrURLNotFound // return statement
		} // end block
		return domain.ShortURL{}, err // return statement
	} // end block

	if !url.IsActive { // if condition
		return domain.ShortURL{}, ErrURLInactive // return statement
	} // end block

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) { // if condition
		return domain.ShortURL{}, ErrURLExpired // return statement
	} // end block

	return url, nil // return statement
} // end block

func (s *URLService) TrackClick(ctx context.Context, click domain.Click) error { // declare method
	return s.clickRepo.Create(ctx, click) // record click event
} // end block

func (s *URLService) Analytics(ctx context.Context, userID, urlID uuid.UUID) (domain.AnalyticsResponse, error) { // declare method
	url, err := s.urlRepo.FindByIDAndUserID(ctx, urlID, userID) // get URL by ID and owner
	if err != nil { // if condition
		if errors.Is(err, repository.ErrNotFound) { // handle missing database record
			return domain.AnalyticsResponse{}, ErrURLNotFound // return statement
		} // end block
		return domain.AnalyticsResponse{}, err // return statement
	} // end block

	total, err := s.clickRepo.CountByURLID(ctx, url.ID) // count clicks for a URL
	if err != nil { // if condition
		return domain.AnalyticsResponse{}, err // return statement
	} // end block

	daily, err := s.clickRepo.DailyStatsByURLID(ctx, url.ID, 30) // load daily click stats
	if err != nil { // if condition
		return domain.AnalyticsResponse{}, err // return statement
	} // end block

	return domain.AnalyticsResponse{ // return statement
		URLID:       url.ID, // execute statement
		TotalClicks: total, // execute statement
		DailyClicks: daily, // execute statement
	}, nil // close block
} // end block

func (s *URLService) generateUniqueShortCode(ctx context.Context) (string, error) { // declare method
	for i := 0; i < 5; i++ { // for loop
		code, err := utils.GenerateShortCode(s.shortCodeLength) // declare and initialize variable
		if err != nil { // if condition
			return "", err // return statement
		} // end block

		exists, err := s.urlRepo.ShortCodeExists(ctx, code) // check existing short code
		if err != nil { // if condition
			return "", err // return statement
		} // end block
		if !exists { // if condition
			return code, nil // return statement
		} // end block
	} // end block

	return "", errors.New("failed to generate unique short code") // return statement
} // end block
