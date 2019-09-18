package psa

import "time"

// Site ...
type Site struct {
	Name     string
	SiteCode string
}

// Excludes detail tickets that should be exluded from scorecard
type Excludes struct {
	Summary []string `json:"summary"`
}

// Ticket ...
type Ticket struct {
	ID          int       `json:"id"`
	DateEntered time.Time `json:"dateEntered"`
	Company     Company   `json:"company"`
	Board       Board     `json:"board"`
	Summary     string    `josn:"summary"`
	Info        Info      `json:"_info"`
}

// Info ..
type Info struct {
	LastUpdated         time.Time `json:"lastUpdated"`
	UpdatedBy           string    `json:"updatedBy"`
	DateEntered         time.Time `json:"dateEntered"`
	EnteredBy           string    `json:"enteredBy"`
	ActivitiesHRef      string    `json:"activities_href"`
	ScheduleentriesHRef string    `json:"scheduleentries_href"`
	DocumentsHRef       string    `json:"documents_href"`
	ConfigurationsHRef  string    `json:"configurations_href"`
	TasksHRef           string    `json:"tasks_href"`
	NotesHRef           string    `json:"notes_href"`
	ProductsHRef        string    `json:"products_href"`
	TimeentriesHRef     string    `json:"timeentries_href"`
	ExpenseEntriesHRef  string    `json:"expenseEntries_href"`
}

// Company ...
type Company struct {
	ID       int    `json:"id"`
	SiteCode string `json:"identifier"`
	Name     string `json:"name"`
}

// Board ...
type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Query ...
type Query struct {
	OrderBy    string `json:"orderBy"`
	Conditions string `json:"conditions"`
}

// TicketNote ..
type TicketNote struct {
	ID                    int       `json:"id"`
	TicketID              int       `json:"ticketId"`
	Text                  string    `json:"text"`
	DetailDescriptionFlag string    `json:"DetailDescriptionFlag"`
	InternalAnalysisFlag  string    `json:"InternalAnalysisFlag"`
	ResolutionFlag        string    `json:"ResolutionFlag"`
	Member                Member    `json:"Member"`
	Contact               Contact   `json:"Contact"`
	CustomerUpdatedFlag   string    `json:"CustomerUpdatedFlag"`
	ProcessNotifications  string    `json:"ProcessNotifications"`
	DateCreated           time.Time `json:"DateCreated"`
	CreatedBy             string    `json:"CreatedBy"`
	InternalFlag          string    `json:"InternalFlag"`
	ExternalFlag          string    `json:"ExternalFlag"`
}

// Member ..
type Member struct {
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// Contact ..
type Contact struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
