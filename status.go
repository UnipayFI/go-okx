package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
)

// SystemStatusState is the lifecycle state of a system maintenance window.
type SystemStatusState string

const (
	SystemStatusStateScheduled SystemStatusState = "scheduled"
	SystemStatusStateOngoing   SystemStatusState = "ongoing"
	SystemStatusStatePreOpen   SystemStatusState = "pre_open"
	SystemStatusStateCompleted SystemStatusState = "completed"
	SystemStatusStateCanceled  SystemStatusState = "canceled"
)

// GetSystemStatusService -- GET /api/v5/system/status (public)
//
// Returns scheduled and ongoing system maintenance windows. The data array is
// empty when no maintenance is planned or in progress.
type GetSystemStatusService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSystemStatusService() *GetSystemStatusService {
	return &GetSystemStatusService{c: c, params: map[string]string{}}
}

// SetState filters by maintenance state (scheduled / ongoing / pre_open /
// completed / canceled).
func (s *GetSystemStatusService) SetState(state SystemStatusState) *GetSystemStatusService {
	s.params["state"] = string(state)
	return s
}

func (s *GetSystemStatusService) Do(ctx context.Context) ([]SystemStatus, error) {
	req := request.Get(ctx, s.c, "/api/v5/system/status", s.params)
	return request.DoList[SystemStatus](req)
}

// SystemStatus is one OKX system maintenance window.
type SystemStatus struct {
	Title               string            `json:"title"`
	State               SystemStatusState `json:"state"`
	Begin               time.Time         `json:"begin"`
	End                 time.Time         `json:"end"`
	PreOpenBegin        time.Time         `json:"preOpenBegin"`
	Href                string            `json:"href"`
	ServiceType         string            `json:"serviceType"`
	System              string            `json:"system"`
	ScheduleDescription string            `json:"scheDesc"`
	MaintenanceType     string            `json:"maintType"`
	Env                 string            `json:"env"`
}

// GetAnnouncementsService -- GET /api/v5/support/announcements (signed)
//
// Returns paginated OKX support announcements. Each data element wraps a page of
// announcement details plus the total page count.
type GetAnnouncementsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAnnouncementsService() *GetAnnouncementsService {
	return &GetAnnouncementsService{c: c, params: map[string]string{}}
}

// SetAnnType filters by announcement type (see GetAnnouncementTypesService for
// the available values).
func (s *GetAnnouncementsService) SetAnnType(annType string) *GetAnnouncementsService {
	s.params["annType"] = annType
	return s
}

// SetPage selects the page of results (1-based).
func (s *GetAnnouncementsService) SetPage(page int) *GetAnnouncementsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

func (s *GetAnnouncementsService) Do(ctx context.Context) ([]Announcements, error) {
	req := request.Get(ctx, s.c, "/api/v5/support/announcements", s.params).WithSign()
	return request.DoList[Announcements](req)
}

// Announcements is one page of OKX support announcements.
type Announcements struct {
	TotalPage string             `json:"totalPage"`
	Details   []AnnouncementItem `json:"details"`
}

// AnnouncementItem is a single announcement entry.
type AnnouncementItem struct {
	AnnouncementType string    `json:"annType"`
	Title            string    `json:"title"`
	URL              string    `json:"url"`
	PushTime         time.Time `json:"pTime"`
	BusinessPTime    time.Time `json:"businessPTime"`
}

// GetAnnouncementTypesService -- GET /api/v5/support/announcement-types (signed)
//
// Returns the set of announcement types usable as the annType filter of
// GetAnnouncementsService.
type GetAnnouncementTypesService struct {
	c *Client
}

func (c *Client) NewGetAnnouncementTypesService() *GetAnnouncementTypesService {
	return &GetAnnouncementTypesService{c: c}
}

func (s *GetAnnouncementTypesService) Do(ctx context.Context) ([]AnnouncementType, error) {
	req := request.Get(ctx, s.c, "/api/v5/support/announcement-types").WithSign()
	return request.DoList[AnnouncementType](req)
}

// AnnouncementType is a selectable announcement category.
type AnnouncementType struct {
	AnnouncementType            string `json:"annType"`
	AnnouncementTypeDescription string `json:"annTypeDesc"`
}

// EconomicCalendarImportance is the market-impact rating of a calendar event.
type EconomicCalendarImportance string

const (
	EconomicCalendarImportanceLow    EconomicCalendarImportance = "1"
	EconomicCalendarImportanceMedium EconomicCalendarImportance = "2"
	EconomicCalendarImportanceHigh   EconomicCalendarImportance = "3"
)

// GetEconomicCalendarService -- GET /api/v5/public/economic-calendar (signed)
//
// Returns macro-economic calendar events (releases, forecasts and actuals)
// scoped by region and importance.
type GetEconomicCalendarService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetEconomicCalendarService() *GetEconomicCalendarService {
	return &GetEconomicCalendarService{c: c, params: map[string]string{}}
}

// SetRegion filters by country/region.
func (s *GetEconomicCalendarService) SetRegion(region string) *GetEconomicCalendarService {
	s.params["region"] = region
	return s
}

// SetImportance filters by market-impact rating (1 low / 2 medium / 3 high).
func (s *GetEconomicCalendarService) SetImportance(importance EconomicCalendarImportance) *GetEconomicCalendarService {
	s.params["importance"] = string(importance)
	return s
}

// SetBefore returns events occurring before the given time (pagination by
// event date, exclusive).
func (s *GetEconomicCalendarService) SetBefore(before time.Time) *GetEconomicCalendarService {
	s.params["before"] = strconv.FormatInt(before.UnixMilli(), 10)
	return s
}

// SetAfter returns events occurring after the given time (pagination by event
// date, exclusive).
func (s *GetEconomicCalendarService) SetAfter(after time.Time) *GetEconomicCalendarService {
	s.params["after"] = strconv.FormatInt(after.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of results (default 100, max 100).
func (s *GetEconomicCalendarService) SetLimit(limit int) *GetEconomicCalendarService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetEconomicCalendarService) Do(ctx context.Context) ([]EconomicCalendar, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/economic-calendar", s.params).WithSign()
	return request.DoList[EconomicCalendar](req)
}

// EconomicCalendar is one macro-economic calendar event. The numeric-looking
// fields (actual, previous, forecast, prevInitial) are kept as strings because
// OKX returns them with trailing units/symbols (e.g. "2.4%", "1.2K", "-").
type EconomicCalendar struct {
	CalendarID      string                     `json:"calendarId"`
	Date            time.Time                  `json:"date"`
	Region          string                     `json:"region"`
	Category        string                     `json:"category"`
	Event           string                     `json:"event"`
	ReferenceDate   time.Time                  `json:"refDate"`
	Actual          string                     `json:"actual"`
	Previous        string                     `json:"previous"`
	Forecast        string                     `json:"forecast"`
	DateSpan        string                     `json:"dateSpan"`
	Importance      EconomicCalendarImportance `json:"importance"`
	UpdateTime      time.Time                  `json:"uTime"`
	PreviousInitial string                     `json:"prevInitial"`
	Currency        string                     `json:"ccy"`
	Unit            string                     `json:"unit"`
}
