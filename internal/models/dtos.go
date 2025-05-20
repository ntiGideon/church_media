package models

import (
	"github.com/ogidi/church-media/ent/story"
	"github.com/ogidi/church-media/ent/user"
	"github.com/ogidi/church-media/internal/validator"
	"time"
)

type Member struct {
	ID                int
	IdNumber          string
	Surname           string
	OtherName         string
	Dob               time.Time
	Gender            string
	HomeTown          string
	Region            string
	Residency         string
	Address           string
	Mobile            string
	Email             string
	SundaySchoolClass string
	Occupation        string
	HasTitheCard      bool
	TitheCardNumber   string
	DayBorn           string
	HasSpouse         bool
	SpouseIdNumber    string
	SpouseName        string
	SpouseOccupation  string
	SpouseContact     string
	IsBaptised        bool
	BaptisedBy        string
	BaptismChurch     string
	BaptismCertNumber string
	BaptismDate       time.Time
	MembershipYear    int
	PhotoURL          string
	PhotoData         string
	IsActive          bool
	CreatedAt         time.Time
}

type CreateMessageDto struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Phone               string `form:"phone"`
	Subject             string `form:"subject"`
	Description         string `form:"description"`
	validator.Validator `form:"-"`
}

type CreateServiceDto struct {
	Name                string
	Types               string    `form:"service_type"`
	Males               int       `form:"males"`
	Females             int       `form:"females"`
	Offering            float64   `form:"offering"`
	Tithe               float64   `form:"tithe"`
	Date                time.Time `form:"date"`
	validator.Validator `form:"-"`
}

type ServiceStatistics struct {
	TotalServices      int
	TotalAttendance    int
	TotalOffering      float64
	TotalTithe         float64
	AverageAttendance  float64
	RecentServices     int
	HighestAttendance  int
	HighestOffering    float64
	LastServiceDate    time.Time
	MemberDemographics map[string]int
}

type CreateMemberDto struct {
	IdNumber  string    `form:"id_number"`
	Surname   string    `form:"surname"`
	OtherName string    `form:"other_name"`
	Dob       time.Time `form:"dob" time_format:"2006-01-02"`
	Gender    string    `form:"gender"`

	HomeTown          string    `form:"home_town"`
	Region            string    `form:"region"`
	Residence         string    `form:"residency"`
	Address           string    `form:"address"`
	Mobile            string    `form:"mobile"`
	Email             string    `form:"email"`
	SundaySchoolClass string    `form:"sunday_school_class"`
	Occupation        string    `form:"occupation"`
	HasTitheCard      bool      `form:"has_tithe_card"`
	TitheCardNumber   string    `form:"tithe_card_number"`
	DayBorn           string    `form:"day_born"`
	HasSpouse         bool      `form:"has_spouse"`
	SpouseIdNumber    string    `form:"spouse_id_number"`
	SpouseName        string    `form:"spouse_name"`
	SpouseOccupation  string    `form:"spouse_occupation"`
	SpouseContact     string    `form:"spouse_contact"`
	IsBaptised        bool      `form:"is_baptised"`
	BaptisedBy        string    `form:"baptised_by"`
	BaptismChurch     string    `form:"baptism_church"`
	BaptismCertNumber string    `form:"baptism_cert_number"`
	BaptismDate       time.Time `form:"baptism_date" time_format:"2006-01-02"`
	MembershipYear    int       `form:"membership_year"`
	PhotoURL          string    `form:"photo_url"`
	PhotoData         string    `form:"photo_data"`
	FormNumber        string    `form:"form_number"`

	validator.Validator `form:"-"`
}

type MemberStats struct {
	TotalMembers          int
	ActiveMembers         int
	BaptizedMembers       int
	MembersWithTitheCards int
	MaleMembers           int
	FemaleMembers         int
	RecentBaptisms        int
	BaptismPercentage     float64
	FirstService          ServiceStats
	SecondService         ServiceStats
	WeeklyTrend           []ServiceStats
	TotalOfferings        float64
	TotalTithe            float64

	PreviousYearMembers int
	GrowthRate          float64
	AttendanceRate      float64
	BirthdaysThisMonth  int
	OtherGenderMembers  int
	MalePercentage      float64
	FemalePercentage    float64
}

type RecentMembers struct {
	ID             int
	Surname        string
	OtherNames     string
	Mobile         string
	Email          string
	PhotoURL       string
	MembershipYear int
	IsBaptized     bool
	BaptismDate    time.Time
	CreatedAt      time.Time
}

type ServiceStats struct {
	Date        time.Time
	ServiceType string
	Males       int
	Females     int
	Total       int
	Offering    float64
	Tithe       float64
}

type ServiceData struct {
	Males    int
	Females  int
	Offering float64
	Tithe    float64
}

type ServiceRecordForm struct {
	Date          time.Time
	FirstService  ServiceData
	SecondService ServiceData
	validator.Validator
}

type CreateEventDto struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	StartTime   string `form:"start_time"`
	EndTime     string `form:"end_time"`
	Location    string `form:"location"`
	ImageUrl    string `form:"image_url"`
	Featured    bool   `form:"featured"`

	validator.Validator `form:"-"`
}

type RegisterDto struct {
	RegistrationToken string `form:"registration_token"`
	Username          string `form:"username"`
	Password          string `form:"password"`
	ConfirmPassword   string `form:"confirm_password"`
	Email             string `form:"email"`

	validator.Validator `form:"-"`
}

type LoginDto struct {
	EmailUsername string `form:"email_username"`
	Password      string `form:"password"`
	RememberMe    bool   `form:"remember_me"`

	validator.Validator `form:"-"`
}

type InviteDto struct {
	Firstname string    `form:"firstname"`
	Lastname  string    `form:"lastname"`
	Email     string    `form:"email"`
	Role      user.Role `form:"role"`
	ExpiresAt int       `form:"expires_at"`

	validator.Validator `form:"-"`
}

type StoriesDto struct {
	Title    string       `form:"title"`
	Body     string       `form:"body"`
	Image    string       `form:"image"`
	Excerpt  string       `form:"excerpt"`
	Status   story.Status `form:"status"`
	AuthorId int          `form:"author_id"`

	validator.Validator `form:"-"`
}

type ToastDto interface{}

type UpdateProfileRequest struct {
	FirstName       string    `form:"first_name"`
	Surname         string    `form:"surname"`
	PhoneNumber     string    `form:"phone_number"`
	DateOfBirth     time.Time `form:"date_of_birth"`
	Gender          string    `form:"gender"`
	Department      string    `form:"department"`
	Occupation      string    `form:"occupation"`
	MaritalStatus   string    `form:"marital_status"`
	Address         string    `form:"address"`
	CurrentPassword string    `form:"current_password"`
	NewPassword     string    `form:"new_password"`
	ConfirmPassword string    `form:"confirm_password"`
	ProfilePicture  string    `form:"profile_picture"`

	validator.Validator `form:"-"`
}

type Filters struct {
	Role   string
	Status string
	Search string
}

type Sort struct {
	Field string
	Order string
}

type Pagination struct {
	CurrentPage int
	PageSize    int
	TotalPages  int
	TotalItems  int
}

type EventFilter struct {
	Search   string
	Featured *bool
	FromDate time.Time
	ToDate   time.Time
	SortBy   string
	Order    string
	Page     int
	PageSize int
}

type AuditLogEntry struct {
	Action     string
	EntityType string
	EntityID   int
	EntityData map[string]interface{}
	CreatedBy  *int
	IPAddress  string
	UserAgent  string
	RequestID  string
	Metadata   map[string]interface{}
}
