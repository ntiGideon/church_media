package models

import (
	"context"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/member"
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
	BaptismDate       string
	MembershipYear    int
	PhotoUrl          string
	IsActive          bool
	CreatedAt         time.Time
}

type MemberModelInterface interface {
	CreateMember(ctx context.Context, dto *CreateMemberDto) error
}

type MemberModel struct {
	Db *ent.Client
}

func (m *MemberModel) CreateMember(ctx context.Context, dto CreateMemberDto) error {
	existing, err := m.Db.Member.Query().Where(member.Mobile(dto.Mobile)).Exist(ctx)
	if err != nil {
		return err
	}
	if existing {
		return PhoneInUse
	}

	_, err = m.Db.Member.Create().
		SetIDNumber(dto.IdNumber).
		SetSurname(dto.Surname).
		SetOtherNames(dto.OtherName).
		SetDob(dto.Dob).
		SetGender(member.Gender(dto.Gender)).
		SetHometown(dto.HomeTown).
		SetRegion(dto.Region).
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
		SetPhotoURL(dto.PhotoUrl).
		SetIsActive(dto.IsActive).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
