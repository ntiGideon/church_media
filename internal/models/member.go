package models

import (
	"context"
	"fmt"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/member"
	"strconv"
	"strings"
	"time"
)

type MemberModelInterface interface {
	CreateMember(ctx context.Context, dto *CreateMemberDto) error
}

type MemberModel struct {
	Db *ent.Client
}

func (m *MemberModel) GetMemberById(ctx context.Context, id int) (*ent.Member, error) {
	memberData, err := m.Db.Member.Query().Where(member.IDEQ(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return memberData, nil
}

func (m *MemberModel) ListMembers(ctx context.Context, page int, pageSize int, sortField string, sortOrder string) ([]*ent.Member, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := m.Db.Member.Query()

	// Default sort by createdAt if no sortField provided
	if sortField == "" {
		sortField = member.FieldCreatedAt // Assuming the field is named "created_at"
		sortOrder = "desc"                // default sort order, newest first
	}

	// Apply sorting
	order := ent.Asc(sortField)
	if sortOrder == "desc" {
		order = ent.Desc(sortField)
	}
	query = query.Order(order)

	// Count total records
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	members, err := query.
		Offset(offset).
		Limit(pageSize).
		Select(
			member.FieldID,
			member.FieldSurname,
			member.FieldOtherNames,
			member.FieldPhotoURL,
			member.FieldResidence,
			member.FieldGender,
			member.FieldDob,
		).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return members, total, nil
}

func (m *MemberModel) DeleteMember(ctx context.Context, id int) error {
	err := m.Db.Member.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *MemberModel) GetMonthlyGrowth2(ctx context.Context) ([]string, []int, error) {
	// Get current year
	currentYear := time.Now().Year()

	// Initialize results
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	counts := make([]int, 12)

	for i := 0; i < 12; i++ {
		month := i + 1
		startDate := time.Date(currentYear, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		var endDate time.Time

		if month == 12 {
			endDate = time.Date(currentYear+1, 1, 1, 0, 0, 0, 0, time.UTC)
		} else {
			endDate = time.Date(currentYear, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
		}

		count, err := m.Db.Member.Query().
			Where(
				member.CreatedAtGTE(startDate),
				member.CreatedAtLT(endDate),
			).
			Count(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to count members for month %d: %w", month, err)
		}

		counts[i] = count
	}

	return months, counts, nil
}

func (m *MemberModel) GetYearlyGrowth(ctx context.Context) ([]string, []int, error) {
	// Get current year
	currentYear := time.Now().Year()
	years := []int{currentYear - 3, currentYear - 2, currentYear - 1, currentYear}
	counts := make([]int, len(years))

	for i, year := range years {
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

		count, err := m.Db.Member.Query().
			Where(
				member.CreatedAtGTE(startDate),
				member.CreatedAtLT(endDate),
			).
			Count(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to count members for year %d: %w", year, err)
		}

		counts[i] = count
	}

	// Convert years to strings
	yearLabels := make([]string, len(years))
	for i, y := range years {
		yearLabels[i] = strconv.Itoa(y)
	}

	return yearLabels, counts, nil
}

func (m *MemberModel) GetMemberAgeDistribution(ctx context.Context) ([]string, []int, error) {
	// Define our age groups
	ageGroups := []struct {
		name   string
		minAge int
		maxAge int
	}{
		{"0-18", 0, 18},
		{"19-35", 19, 35},
		{"36-50", 36, 50},
		{"51-65", 51, 65},
		{"65+", 66, 150}, // Using 150 as upper limit for "65+"
	}

	// Get current year for age calculation
	currentYear := time.Now().Year()

	var labels []string
	var counts []int

	// Query count for each age group
	for _, group := range ageGroups {
		count, err := m.Db.Member.Query().
			Where(
				member.And(
					// Calculate age: currentYear - birthYear >= minAge
					member.DobGTE(time.Date(currentYear-group.maxAge, 1, 1, 0, 0, 0, 0, time.UTC)),
					// Calculate age: currentYear - birthYear <= maxAge
					member.DobLTE(time.Date(currentYear-group.minAge, 12, 31, 23, 59, 59, 0, time.UTC)),
				),
			).
			Count(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to count age group %s: %w", group.name, err)
		}

		labels = append(labels, group.name)
		counts = append(counts, count)
	}

	return labels, counts, nil
}

func (m *MemberModel) GetDistinctRegions(ctx context.Context) ([]string, []int, error) {
	regions, err := m.Db.Member.
		Query().
		GroupBy(member.FieldRegion).
		Strings(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get distinct regions: %w", err)
	}

	// Then get counts for each region
	var counts []int
	for _, region := range regions {
		count, err := m.Db.Member.
			Query().
			Where(member.RegionEQ(region)).
			Count(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to count members in region %s: %w", region, err)
		}
		counts = append(counts, count)
	}

	return regions, counts, nil
}

func (m *MemberModel) GetMemberStatistics(ctx context.Context) (*MemberStats, error) {
	stats := &MemberStats{}

	// Start a transaction
	tx, err := m.Db.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	stats.TotalMembers, err = tx.Member.Query().
		Count(ctx)
	if err != nil {
		return nil, err
	}

	stats.ActiveMembers, err = tx.Member.Query().
		Where(member.IsActive(true)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	stats.BaptizedMembers, err = tx.Member.Query().
		Where(member.IsActive(true), member.IsBaptized(true)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	if stats.TotalMembers > 0 {
		stats.BaptismPercentage = float64(stats.BaptizedMembers) / float64(stats.TotalMembers) * 100
	} else {
		stats.BaptismPercentage = 0
	}

	stats.MembersWithTitheCards, err = tx.Member.Query().
		Where(member.IsActive(true), member.HasTitleCard(true)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	stats.MaleMembers, err = tx.Member.Query().
		Where(member.IsActive(true), member.GenderEQ(member.GenderMale)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	stats.FemaleMembers, err = tx.Member.Query().
		Where(member.IsActive(true), member.GenderEQ(member.GenderFemale)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	stats.RecentBaptisms, err = tx.Member.Query().
		Where(member.IsActive(true), member.IsBaptized(true), member.BaptismDateGTE(time.Now().AddDate(0, 0, -30))).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (m *MemberModel) GetRecentMembers(ctx context.Context, limit int) ([]*RecentMembers, error) {
	members, err := m.Db.Member.Query().
		Where(member.IsActive(true)).
		Order(ent.Desc(member.FieldCreatedAt)).
		Limit(limit).All(ctx)
	if err != nil {
		return nil, err
	}
	var result []*RecentMembers
	for _, m := range members {
		result = append(result, &RecentMembers{
			ID:             m.ID,
			Surname:        m.Surname,
			OtherNames:     m.OtherNames,
			Mobile:         m.Mobile,
			Email:          m.Email,
			PhotoURL:       m.PhotoURL,
			MembershipYear: m.MembershipYear,
			IsBaptized:     m.IsBaptized,
			BaptismDate:    m.BaptismDate,
			CreatedAt:      m.CreatedAt,
		})
	}
	return result, nil
}

func (m *MemberModel) UpdateMember(ctx context.Context, id int, dto CreateMemberDto) error {
	// Get the existing member from the database
	existingMember, err := m.Db.Member.Get(ctx, id)
	if err != nil {
		return err
	}

	// Normalize and trim mobile and email inputs
	cleanMobile := strings.TrimSpace(dto.Mobile)
	cleanEmail := strings.TrimSpace(strings.ToLower(dto.Email))
	existingMobile := strings.TrimSpace(existingMember.Mobile)
	existingEmail := strings.TrimSpace(strings.ToLower(existingMember.Email))

	// Check for mobile number conflict
	if cleanMobile != "" && cleanMobile != existingMobile {
		mobileExists, err := m.Db.Member.Query().
			Where(member.MobileEQ(cleanMobile)).
			Where(member.IDNEQ(id)). // ⬅️ MUST exclude current member!
			Exist(ctx)
		if err != nil {
			return err
		}
		if mobileExists {
			return PhoneInUse
		}
	}

	if cleanEmail != "" && cleanEmail != existingEmail {
		emailExists, err := m.Db.Member.Query().
			Where(member.EmailEQ(cleanEmail)).
			Where(member.IDNEQ(id)). // ⬅️ MUST exclude current member!
			Exist(ctx)
		if err != nil {
			return err
		}
		if emailExists {
			return EmailConstraint
		}
	}

	// Start building the update statement
	update := m.Db.Member.UpdateOneID(id).
		SetIDNumber(dto.IdNumber).
		SetFormNumber(dto.FormNumber).
		SetSurname(dto.Surname).
		SetOtherNames(dto.OtherName).
		SetDob(dto.Dob).
		SetGender(member.Gender(dto.Gender)).
		SetHometown(dto.HomeTown).
		SetRegion(dto.Region).
		SetResidence(dto.Residence).
		SetAddress(dto.Address).
		SetSundaySchoolClass(dto.SundaySchoolClass).
		SetOccupation(dto.Occupation).
		SetHasTitleCard(dto.HasTitheCard).
		SetDayBorn(dto.DayBorn).
		SetHasSpouse(dto.HasSpouse).
		SetSpouseIDNumber(dto.SpouseIdNumber).
		SetSpouseName(dto.SpouseName).
		SetSpouseOccupation(dto.SpouseOccupation).
		SetSpouseContact(dto.SpouseContact).
		SetIsBaptized(dto.IsBaptised).
		SetBaptizedBy(dto.BaptisedBy).
		SetBaptismChurch(dto.BaptismChurch).
		SetBaptismCertNumber(dto.BaptismCertNumber).
		SetBaptismDate(dto.BaptismDate).
		SetMembershipYear(dto.MembershipYear)

	// Conditionally update mobile
	if cleanMobile != "" && cleanMobile != existingMobile {
		update = update.SetMobile(cleanMobile)
	}

	// Conditionally update email
	if cleanEmail != "" && cleanEmail != existingEmail {
		update = update.SetEmail(cleanEmail)
	}

	// Conditionally update photo URL
	if strings.TrimSpace(dto.PhotoURL) != "" {
		update = update.SetPhotoURL(strings.TrimSpace(dto.PhotoURL))
	}

	// Save the updated member
	_, err = update.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	return nil
}

func (m *MemberModel) CreateMember(ctx context.Context, dto CreateMemberDto) error {
	// Check for existing mobile (only if it's not empty)
	if dto.Mobile != "" {
		existing, err := m.Db.Member.Query().Where(member.MobileEQ(dto.Mobile)).Exist(ctx)
		if err != nil {
			return err
		}
		if existing {
			return PhoneInUse
		}
	}

	// Check for existing email (only if it's not empty)
	if dto.Email != "" {
		emailExist, err := m.Db.Member.Query().Where(member.EmailEQ(dto.Email)).Exist(ctx)
		if err != nil {
			return err
		}
		if emailExist {
			return EmailConstraint
		}
	}

	_, err := m.Db.Member.Create().
		SetIDNumber(dto.IdNumber).
		SetFormNumber(dto.FormNumber).
		SetSurname(dto.Surname).
		SetOtherNames(dto.OtherName).
		SetDob(dto.Dob).
		SetGender(member.Gender(dto.Gender)).
		SetHometown(dto.HomeTown).
		SetRegion(dto.Region).
		SetResidence(dto.Residence).
		SetAddress(dto.Address).
		SetMobile(dto.Mobile).
		SetEmail(dto.Email).
		SetSundaySchoolClass(dto.SundaySchoolClass).
		SetOccupation(dto.Occupation).
		SetHasTitleCard(dto.HasTitheCard).
		SetDayBorn(dto.DayBorn).
		SetHasSpouse(dto.HasSpouse).
		SetSpouseIDNumber(dto.SpouseIdNumber).
		SetSpouseName(dto.SpouseName).
		SetSpouseOccupation(dto.SpouseOccupation).
		SetSpouseContact(dto.SpouseContact).
		SetIsBaptized(dto.IsBaptised).
		SetBaptizedBy(dto.BaptisedBy).
		SetBaptismChurch(dto.BaptismChurch).
		SetBaptismCertNumber(dto.BaptismCertNumber).
		SetBaptismDate(dto.BaptismDate).
		SetMembershipYear(dto.MembershipYear).
		SetIsActive(true).
		SetPhotoURL(dto.PhotoURL).
		Save(ctx)
	return err
}
