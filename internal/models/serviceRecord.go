package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/attendancerecord"
	_ "github.com/ogidi/church-media/ent/attendancerecord"
	"github.com/ogidi/church-media/ent/member"
	"github.com/ogidi/church-media/ent/service"
	"time"
)

type ServiceModel struct {
	Db *ent.Client
}

func (m *ServiceModel) CreateServiceRecord(ctx context.Context, dto *CreateServiceDto) error {
	// transaction
	tx, err := m.Db.Tx(ctx)
	if err != nil {
		return errors.New("failed to begin transaction")
	}

	defer func() {
		if err := recover(); err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			panic(err)
		} else if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		} else {
			err = tx.Commit()
		}
	}()

	serviceCreated, err := tx.Service.Create().SetDate(dto.Date).SetName(dto.Name).SetType(service.Type(dto.Types)).Save(ctx)
	if err != nil {
		return err
	}

	_, err = tx.AttendanceRecord.Create().
		SetMales(dto.Males).
		SetFemales(dto.Females).
		SetOffering(dto.Offering).
		SetTithe(dto.Tithe).
		SetService(serviceCreated).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *ServiceModel) GetAllServiceRecords(
	ctx context.Context,
	page int,
	pageSize int,
	sortField string,
	sortOrder string,
	serviceType string,
	dateFilter string,
) ([]*ent.AttendanceRecord, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Base query with relation
	query := m.Db.AttendanceRecord.Query().WithService().Offset((page - 1) * pageSize).Limit(pageSize)

	// Filter by service type if provided
	if serviceType != "" {
		query = query.Where(
			attendancerecord.HasServiceWith(
				service.TypeEQ(service.Type(serviceType)),
			),
		)
	}

	// Sorting logic
	switch sortField {
	case "males":
		if sortOrder == "desc" {
			query = query.Order(ent.Desc("males"))
		} else {
			query = query.Order(ent.Asc("males"))
		}
	case "females":
		if sortOrder == "desc" {
			query = query.Order(ent.Desc("females"))
		} else {
			query = query.Order(ent.Asc("females"))
		}
	case "offering":
		if sortOrder == "desc" {
			query = query.Order(ent.Desc("offering"))
		} else {
			query = query.Order(ent.Asc("offering"))
		}
	case "tithe":
		if sortOrder == "desc" {
			query = query.Order(ent.Desc("tithe"))
		} else {
			query = query.Order(ent.Asc("tithe"))
		}
	default:
		query = query.Order(ent.Desc("id"))
	}

	// Fetch filtered + paginated records
	records, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Count total number of filtered records
	countQuery := m.Db.AttendanceRecord.Query()
	if serviceType != "" {
		countQuery = countQuery.Where(
			attendancerecord.HasServiceWith(
				service.TypeEQ(service.Type(serviceType)),
			),
		)
	}

	if dateFilter != "" {
		filterDate, err := time.Parse("2006-01-02", dateFilter)
		if err == nil {
			startOfDay := time.Date(filterDate.Year(), filterDate.Month(), filterDate.Day(), 0, 0, 0, 0, filterDate.Location())
			endOfDay := startOfDay.Add(24 * time.Hour)
			query = query.Where(
				attendancerecord.HasServiceWith(
					service.DateGTE(startOfDay),
					service.DateLT(endOfDay),
				),
			)

			countQuery = countQuery.Where(
				attendancerecord.HasServiceWith(service.DateGTE(startOfDay), service.DateLT(endOfDay)))
		}
	}

	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (m *ServiceModel) DeleteRecord(ctx context.Context, id int) error {
	err := m.Db.AttendanceRecord.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *ServiceModel) GetServiceStatistics(ctx context.Context, serviceType string) (*ServiceStatistics, error) {
	stats := &ServiceStatistics{}

	// Base query for services
	serviceQuery := m.Db.Service.Query()
	recordQuery := m.Db.AttendanceRecord.Query()

	// Apply filter if specified
	if serviceType != "" {
		serviceQuery = serviceQuery.Where(service.TypeEQ(service.Type(serviceType)))
		recordQuery = recordQuery.Where(
			attendancerecord.HasServiceWith(
				service.TypeEQ(service.Type(serviceType)),
			),
		)
	}

	// Get total services count
	total, err := serviceQuery.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count services: %w", err)
	}
	stats.TotalServices = total

	// Get total attendance
	//totalMales, err := recordQuery.Aggregate(ent.Sum("males")).Int(ctx)
	//totalFemales, err := recordQuery.Aggregate(ent.Sum("females")).Int(ctx)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to count services: %w", err)
	//}
	//stats.TotalAttendance = totalMales + totalFemales

	totalAttendance, err := recordQuery.Clone().
		Aggregate(
			ent.Sum(attendancerecord.FieldMales),
			ent.Sum(attendancerecord.FieldFemales),
		).Ints(ctx)
	if len(totalAttendance) == 2 {
		stats.TotalServices = totalAttendance[0] + totalAttendance[1]
	}

	// Get total offering
	totalOffering, err := recordQuery.
		Clone(). // Need to clone because we can't reuse after aggregation
		Aggregate(ent.Sum("offering")).
		Float64(ctx)
	stats.TotalOffering = totalOffering

	// Get total tithe
	totalTithe, err := recordQuery.
		Clone().
		Aggregate(ent.Sum("tithe")).
		Float64(ctx)
	stats.TotalTithe = totalTithe

	if stats.TotalServices > 0 {
		stats.AverageAttendance = float64(stats.TotalAttendance) / float64(stats.TotalServices)
	}

	recentServices, err := serviceQuery.Clone().Where(service.DateGTE(time.Now().AddDate(0, 0, -30))).Count(ctx)
	stats.RecentServices = recentServices

	return stats, nil
}

func (m *MemberModel) GetBirthdaysThisMonth(ctx context.Context) (int, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0)

	return m.Db.Member.Query().
		Where(member.DobGTE(start), member.DobLT(end)).
		Count(ctx)
}

// GetBirthdayMembers returns members whose birthday is today
func (m *MemberModel) GetBirthdayMembers(ctx context.Context) ([]*ent.Member, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	return m.Db.Member.Query().
		Where(member.DobGTE(start), member.DobLT(end)).
		All(ctx)
}

func (m *MemberModel) GetMonthlyGrowth(ctx context.Context) ([]int, error) {
	// Get current year and calculate start date
	now := time.Now()
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	// Query to count members by month of creation
	_, err := m.Db.Member.Query().
		Where(member.CreatedAtGTE(startOfYear)).
		GroupBy(member.FieldCreatedAt).
		Aggregate(ent.Count()).
		Ints(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly growth: %w", err)
	}

	// For Ent-only solution (less efficient):
	counts := make([]int, 12)
	for month := 1; month <= 12; month++ {
		start := time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0)
		count, err := m.Db.Member.Query().
			Where(member.CreatedAtGTE(start), member.CreatedAtLT(end)).
			Count(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to count members for month %d: %w", month, err)
		}
		counts[month-1] = count
	}

	// Calculate cumulative growth
	growth := make([]int, 12)
	total := 0
	for i, count := range counts {
		total += count
		growth[i] = total
	}

	return growth, nil
}

func (m *MemberModel) GetAgeDistribution(ctx context.Context) ([]int, error) {
	// Age groups: 0-18, 19-35, 36-50, 51-65, 65+
	ageGroups := make([]int, 5)

	// Current date for age calculation
	now := time.Now()

	// For Ent-only solution (less efficient):
	members, err := m.Db.Member.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}

	for _, member := range members {
		age := now.Year() - member.Dob.Year()
		if now.YearDay() < member.Dob.YearDay() {
			age-- // Adjust for birthdays not yet occurred this year
		}

		switch {
		case age <= 18:
			ageGroups[0]++
		case age <= 35:
			ageGroups[1]++
		case age <= 50:
			ageGroups[2]++
		case age <= 65:
			ageGroups[3]++
		default:
			ageGroups[4]++
		}
	}

	return ageGroups, nil
}
