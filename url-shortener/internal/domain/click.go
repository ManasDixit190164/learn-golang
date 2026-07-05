package domain

import (
	"time"

	"github.com/google/uuid"
)

type Click struct {
	ID        uuid.UUID `json:"id"`
	URLID     uuid.UUID `json:"url_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
	ClickedAt time.Time `json:"clicked_at"`
}

type DailyClickStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type AnalyticsResponse struct {
	URLID       uuid.UUID        `json:"url_id"`
	TotalClicks int64            `json:"total_clicks"`
	DailyClicks []DailyClickStat `json:"daily_clicks"`
}
