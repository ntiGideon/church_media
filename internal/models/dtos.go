package models

import (
	"github.com/ogidi/church-media/internal/validator"
	"time"
)

type CreateMessageDto struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Phone               string `form:"phone"`
	Subject             string `form:"subject"`
	Description         string `form:"description"`
	validator.Validator `form:"-"`
}

type CreateMemberDto struct {
	IdNumber          string    `form:"id_number"`
	Surname           string    `form:"surname"`
	OtherName         string    `form:"other_name"`
	Dob               time.Time `form:"dob"`
	Gender            string    `form:"gender"`
	HomeTown          string    `form:"home_town"`
	Region            string    `form:"region"`
	Residency         string    `form:"residency"`
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
	BaptismDate       time.Time `form:"baptism_date"`
	MembershipYear    int       `form:"membership_year"`
	PhotoUrl          string    `form:"photo_url"`
	IsActive          bool      `form:"is_active"`

	validator.Validator `form:"-"`
}
