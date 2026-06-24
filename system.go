package okx

import (
	"context"
	"time"

	"github.com/UnipayFI/go-okx/request"
)

// GetSystemTimeService -- GET /api/v5/public/time
//
// Returns the current OKX server time. Used by Client.SyncServerTime to align
// the request-signing clock.
type GetSystemTimeService struct {
	c *Client
}

func (c *Client) NewGetSystemTimeService() *GetSystemTimeService {
	return &GetSystemTimeService{c: c}
}

func (s *GetSystemTimeService) Do(ctx context.Context) (*SystemTime, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/time")
	return request.DoOne[SystemTime](req)
}

// SystemTime is the OKX server time.
type SystemTime struct {
	Timestamp time.Time `json:"ts"`
}
