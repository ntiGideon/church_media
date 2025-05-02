package main

import (
	"fmt"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/message"
	"github.com/ogidi/church-media/ent/user"
	"github.com/ogidi/church-media/internal/models"
	"github.com/ogidi/church-media/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

type templateData struct {
	Form            any
	Flash           string
	Toast           map[string]interface{}
	Message         models.Message
	Messages        []*ent.Message
	IsAuthenticated bool
	CurrentYear     int
	CSRFToken       string
	UpcomingEvents  []*ent.Event
	Events          []*ent.Event
	Event           *ent.Event
	Pagination      struct {
		CurrentPage int
		TotalPages  int
		Pages       []int
	}
}

type templateDataAdmin struct {
	Form                any
	Flash               string
	Toast               map[string]interface{}
	Message             models.Message
	Messages            []*ent.Message
	Member              ent.Member
	Members             []*ent.Member
	IsAuthenticated     bool
	CurrentYear         int
	UnreadMessagesCount int
	CSRFToken           string
	NewMemberID         int
	FormNumber          string
	IsEdit              bool
	Stats               models.MemberStats
	ServiceStatistics   models.ServiceStatistics
	RecentMembers       []*models.RecentMembers
	ChartData           any
	Records             []*ent.AttendanceRecord
	Pagination          Pagination
	Sort                models.Sort
	Filters             models.Filters
	CurrentServiceType  string
	CurrentDateFilter   string
	Event               ent.Event
	Events              ent.Events
	User                ent.User
	Users               []*ent.User
}

type Pagination struct {
	CurrentPage int
	PageSize    int
	TotalItems  int
	TotalPages  int
	OffsetStart int
	OffsetEnd   int
	SortField   string
	SortOrder   string
	Pages       []int
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filesystem which match the pattern 'html/pages/*.tmpl'. This essentially
	// gives us a slice of all the 'page' templates for the application, just
	// like before.
	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/base.gohtml",
			"html/partials/*.gohtml",
			page,
		}

		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func newAdminTemplateCache() (map[string]*template.Template, error) {
	adminCache := map[string]*template.Template{}

	// Use fs.Glob() to get all 'pagesAdmin' templates in the application
	adminPages, err := fs.Glob(ui.Files, "html/pagesAdmin/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, adminPage := range adminPages {
		name := filepath.Base(adminPage)

		// Patterns for 'pagesAdmin' templates
		adminPatterns := []string{
			"html/baseAdmin.gohtml",
			"html/partialsAdmin/*.gohtml",
			adminPage,
		}

		// Parse 'pagesAdmin' templates
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, adminPatterns...)
		if err != nil {
			return nil, err
		}

		adminCache[name] = ts
	}

	return adminCache, nil
}

// Calculates the number of hours or day a specific message was created
func daysAgo(t time.Time) string {
	now := time.Now()

	duration := now.Sub(t)

	if duration.Hours() < 24 {
		hours := int(duration.Hours())
		if hours == 1 {
			return fmt.Sprintf("%d hour ago", hours)
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := int(duration.Hours() / 24)
	if days == 1 {
		return fmt.Sprintf("%d day ago", days)
	}
	return fmt.Sprintf("%d days ago", days)
}

func specifyMessageSubject(key message.Subject) string {
	if key == message.SubjectEVENT_INFORMATION {
		return "Event information"
	} else if key == message.SubjectGENERAL_ENQUIRY {
		return "General Enquiry"
	} else if key == message.SubjectMINISTRY_QUESTION {
		return "Ministry Question"
	} else if key == message.SubjectPRAYER_REQUEST {
		return "Prayer Request"
	} else {
		return "Other"
	}
}

func reverseDateFormat(date time.Time) string {
	return strings.Split(date.String(), " ")[0]
}

func calculateAge(dateOfBirth time.Time) int {
	now := time.Now().Year()
	return now - dateOfBirth.Year()
}

func iterate(start, end int) []int {
	n := end - start + 1
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = start + i
	}
	return result
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006")
}

func readableRoles(role user.Role) string {
	if role == user.RoleAdmin {
		return "Admin"
	} else if role == user.RoleSuperadmin {
		return "SuperAdmin"
	} else if role == user.RoleDeacon {
		return "Deacon"
	} else if role == user.RoleMember {
		return "Member"
	} else if role == user.RoleContentAdmin {
		return "Content Admin"
	} else if role == user.RolePastor {
		return "Pastor"
	} else if role == user.RoleSecretary {
		return "Secretary"
	} else {
		return "Other"
	}
}

var functions = template.FuncMap{
	// put configured functions to be passed to the template here!
	"daysAgo":               daysAgo,
	"specifyMessageSubject": specifyMessageSubject,
	"reverseDateFormat":     reverseDateFormat,
	"add":                   func(a, b int) int { return a + b },
	"sub":                   func(a, b int) int { return a - b },
	"calculateAge":          calculateAge,
	"iterate":               iterate,
	"humanDate":             humanDate,
	"readableRoles":         readableRoles,
}
