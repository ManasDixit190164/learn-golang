package domain // package declaration for the module

import ( // start import block
	"time" // import package

	"github.com/google/uuid" // import package
) // end import block or block scope

type Click struct { // declare struct type
	ID        uuid.UUID `json:"id"` // execute statement
	URLID     uuid.UUID `json:"url_id"` // execute statement
	IPAddress string    `json:"ip_address"` // execute statement
	UserAgent string    `json:"user_agent"` // execute statement
	Referrer  string    `json:"referrer"` // execute statement
	ClickedAt time.Time `json:"clicked_at"` // execute statement
} // end block

type DailyClickStat struct { // declare struct type
	Date  string `json:"date"` // execute statement
	Count int64  `json:"count"` // execute statement
} // end block

type AnalyticsResponse struct { // declare struct type
	URLID       uuid.UUID        `json:"url_id"` // execute statement
	TotalClicks int64            `json:"total_clicks"` // execute statement
	DailyClicks []DailyClickStat `json:"daily_clicks"` // execute statement
} // end block
