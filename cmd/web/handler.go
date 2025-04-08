package main

import (
	"errors"
	"github.com/ogidi/church-media/internal/models"
	"github.com/ogidi/church-media/internal/validator"
	"net/http"
	"strconv"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "home.gohtml", pageData)
}

func (app *application) ministries(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "ministries.gohtml", pageData)
}

func (app *application) mensMinistry(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "mens.gohtml", pageData)
}

func (app *application) womenMinistry(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "womens.gohtml", pageData)
}

func (app *application) youthMinistry(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "youth.gohtml", pageData)
}

func (app *application) childrenMinistry(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "children.gohtml", pageData)
}

func (app *application) teensMinistry(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "teens.gohtml", pageData)
}

func (app *application) musicMinistry(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "music.gohtml", pageData)
}

func (app *application) ushersMinistry(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "ushers.gohtml", pageData)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "about.gohtml", pageData)
}

func (app *application) events(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "events.gohtml", pageData)
}

func (app *application) sermon(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "sermon.gohtml", pageData)
}

func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateData(r)
	pageData.Form = models.CreateMessageDto{
		Subject: "GENERAL_ENQUIRY",
	}
	app.render(w, r, http.StatusOK, "contact.gohtml", pageData)
}

func (app *application) contactForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto := models.CreateMessageDto{
		Name:        r.PostForm.Get("name"),
		Email:       r.PostForm.Get("email"),
		Phone:       r.PostForm.Get("phone"),
		Subject:     r.PostForm.Get("subject"),
		Description: r.PostForm.Get("description"),
	}

	err = app.decodePostForm(r, &dto)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto.CheckField(validator.NotBlank(dto.Name), "name", "This field cannot be blank!")
	dto.CheckField(validator.NotBlank(dto.Email), "email", "This field cannot be blank!")
	dto.CheckField(validator.NotBlank(dto.Subject), "subject", "This field cannot be blank!")
	dto.CheckField(validator.NotBlank(dto.Description), "description", "This field cannot be blank!")
	dto.CheckField(validator.Matches(dto.Email, validator.EmailRX), "email", "This field must be a valid email address")

	if !dto.Valid() {
		data := app.newTemplateData(r)
		data.Form = dto
		app.render(w, r, http.StatusUnprocessableEntity, "contact.gohtml", data)
		return
	}

	err = app.messageClient.CreateMessage(r.Context(), dto)
	if err != nil {
		if errors.Is(err, models.PhoneInUse) {
			dto.AddFieldError("phone", "Phone is already in use")
			data := app.newTemplateData(r)
			data.Form = dto
			app.render(w, r, http.StatusUnprocessableEntity, "contact.gohtml", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Message successfully sent!")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) giveOnline(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "giveOnline.gohtml", pageData)
}

// pagesAdmin interfaces
func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	messages, err := app.messageClient.GetMessages(r.Context(), "all") //TODO make it dynamic to support different filters
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Messages = messages
	app.renderAdmin(w, r, http.StatusOK, "admin.gohtml", pageData)
}

func (app *application) members(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateAdmin(r)
	pageData.Form = models.CreateMemberDto{
		HasTitheCard: true,
	}

	app.renderAdmin(w, r, http.StatusOK, "members.gohtml", pageData)
}

func (app *application) memberForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	_, err = strconv.ParseBool(r.PostForm.Get("has_tithe_card"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	_, err = strconv.ParseBool(r.PostForm.Get("is_baptised"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	_, err = strconv.ParseBool(r.PostForm.Get("has_spouse"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	_, err = strconv.ParseBool(r.PostForm.Get("is_active"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	membershipYear, err := strconv.Atoi(r.PostForm.Get("membership_year"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	dateOfBirth, err := time.Parse("2006-01-02", r.PostForm.Get("dob"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	baptismDate, err := time.Parse("2006-01-02", r.PostForm.Get("baptism_date"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto := models.CreateMemberDto{
		IdNumber:          r.PostForm.Get("id_number"),
		Surname:           r.PostForm.Get("surname"),
		OtherName:         r.PostForm.Get("other_name"),
		Dob:               dateOfBirth,
		Gender:            r.PostForm.Get("gender"),
		HomeTown:          r.PostForm.Get("home_town"),
		Region:            r.PostForm.Get("region"),
		Residency:         r.PostForm.Get("residency"),
		Address:           r.PostForm.Get("address"),
		Mobile:            r.PostForm.Get("mobile"),
		Email:             r.PostForm.Get("email"),
		SundaySchoolClass: r.PostForm.Get("sunday_school_class"),
		Occupation:        r.PostForm.Get("occupation"),
		HasTitheCard:      r.PostForm.Get("has_tithe_card") == "on",
		TitheCardNumber:   r.PostForm.Get("tithe_card_number"),
		DayBorn:           r.PostForm.Get("day_born"),
		HasSpouse:         r.PostForm.Get("has_spouse") == "on",
		SpouseIdNumber:    r.PostForm.Get("spouse_id_number"),
		SpouseName:        r.PostForm.Get("spouse_name"),
		SpouseOccupation:  r.PostForm.Get("spouse_occupation"),
		SpouseContact:     r.PostForm.Get("spouse_contact"),
		IsBaptised:        r.PostForm.Get("is_baptised") == "on",
		BaptisedBy:        r.PostForm.Get("baptised_by"),
		BaptismChurch:     r.PostForm.Get("baptism_church"),
		BaptismCertNumber: r.PostForm.Get("baptism_cert_number"),
		BaptismDate:       baptismDate,
		MembershipYear:    membershipYear,
		PhotoUrl:          r.PostForm.Get("member_photo"),
		IsActive:          r.PostForm.Get("is_active") == "on",
	}

	err = app.decodePostForm(r, &dto)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto.CheckField(validator.NotBlank(dto.Surname), "surname", "This field cannot be blank!")
	dto.CheckField(validator.NotBlank(dto.OtherName), "other_name", "This field cannot be blank!")
	dto.CheckField(validator.Matches(dto.Email, validator.EmailRX), "email", "This field must be a valid email address")
	dto.CheckField(validator.NotBlank(dto.Occupation), "occupation", "This field cannot be blank!")
	//TODO: add the rest of the validations

	if !dto.Valid() {
		data := app.newTemplateAdmin(r)
		data.Form = dto
		app.renderAdmin(w, r, http.StatusUnprocessableEntity, "members.gohtml", data)
		return
	}

	err = app.memberClient.CreateMember(r.Context(), dto)
	if err != nil {
		if errors.Is(err, models.PhoneInUse) {
			dto.AddFieldError("mobile", "This phone number is already in use!")
			dataAdmin := app.newTemplateAdmin(r)
			dataAdmin.Form = dto
			app.renderAdmin(w, r, http.StatusUnprocessableEntity, "members.gohtml", dataAdmin)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Member created successfully!")
	//http.Redirect(w, r, fmt.Sprintf("/member/success/%d", newMember.ID), http.StatusSeeOther)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) memberSuccess(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateAdmin(r)
	data.Flash = flash
	data.NewMemberID = id

	app.renderAdmin(w, r, http.StatusOK, "members.gohtml", data)
}
