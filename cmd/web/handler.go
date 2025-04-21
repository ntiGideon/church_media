package main

import (
	"errors"
	"fmt"
	"github.com/ogidi/church-media/internal/models"
	"github.com/ogidi/church-media/internal/validator"
	"io"
	"math"
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
	pageData.IsEdit = false
	pageData.Form = models.CreateMemberDto{
		HasTitheCard: true,
	}

	app.renderAdmin(w, r, http.StatusOK, "members.gohtml", pageData)
}

func (app *application) membersEdit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	memberData, err := app.memberClient.GetMemberById(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	pageData := app.newTemplateAdmin(r)
	pageData.IsEdit = true
	pageData.Form = models.CreateMemberDto{
		IdNumber:          memberData.IDNumber,
		Surname:           memberData.Surname,
		OtherName:         memberData.OtherNames,
		Dob:               memberData.Dob,
		Gender:            string(memberData.Gender),
		HomeTown:          memberData.Hometown,
		Region:            memberData.Region,
		Residence:         memberData.Residence,
		Address:           memberData.Address,
		Mobile:            memberData.Mobile,
		Email:             memberData.Email,
		SundaySchoolClass: memberData.SundaySchoolClass,
		Occupation:        memberData.Occupation,
		HasTitheCard:      memberData.HasTitleCard,
		TitheCardNumber:   memberData.TitleCardNumber,
		DayBorn:           memberData.DayBorn,
		HasSpouse:         memberData.HasSpouse,
		SpouseIdNumber:    memberData.SpouseIDNumber,
		SpouseName:        memberData.SpouseName,
		SpouseOccupation:  memberData.SpouseOccupation,
		SpouseContact:     memberData.SpouseContact,
		IsBaptised:        memberData.IsBaptized,
		BaptisedBy:        memberData.BaptizedBy,
		BaptismChurch:     memberData.BaptismChurch,
		BaptismCertNumber: memberData.BaptismCertNumber,
		BaptismDate:       memberData.BaptismDate,
		MembershipYear:    memberData.MembershipYear,
		PhotoURL:          memberData.PhotoURL,
		PhotoData:         memberData.PhotoURL,
		FormNumber:        memberData.FormNumber,
	}

	app.renderAdmin(w, r, http.StatusOK, "members.gohtml", pageData)
}

func (app *application) deleteMember(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	member, err := app.memberClient.GetMemberById(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = app.memberClient.DeleteMember(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Member %s %s was successfully deleted", member.Surname, member.OtherNames))
	http.Redirect(w, r, "/list-members", http.StatusSeeOther)
}

func (app *application) dashboardData(w http.ResponseWriter, r *http.Request) {
	stats, err := app.memberClient.GetMemberStatistics(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// calculate additional stats needed by the template
	stats.GrowthRate = calculateGrowthRate(stats.TotalMembers, stats.PreviousYearMembers)
	stats.BaptizedMembers = int(calculatePercentage(stats.BaptizedMembers, stats.TotalMembers))
	stats.AttendanceRate = calculateAttendanceRate(stats.ActiveMembers, stats.TotalMembers)
	stats.MalePercentage = calculatePercentage(stats.MaleMembers, stats.TotalMembers)
	stats.FemalePercentage = calculatePercentage(stats.FemaleMembers, stats.TotalMembers)

	stats.BirthdaysThisMonth, err = app.memberClient.GetBirthdaysThisMonth(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	recentMembers, err := app.memberClient.GetRecentMembers(r.Context(), 5)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	chartData := struct {
		GrowthLabels []string
		GrowthValues []int
		AgeGroups    []string
		AgeCounts    []int
		Regions      []string
		RegionCounts []int
		GenderLabels []string
		GenderValues []int
	}{
		GrowthLabels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}, //getMonthlyLabels(), // ["Jan", "Feb", ...]
		GrowthValues: []int{34, 56, 67, 79},                                                                        //getMonthlyGrowth(), // Actual data from DB
		AgeGroups:    []string{"0-18", "19-35", "36-50", "51-65", "65+"},
		AgeCounts:    nil, //getAgeDistribution(), // Actual data from DB
		Regions:      nil, //getRegions(),         // Actual regions from DB
		RegionCounts: nil, //getRegionCounts(),    // Actual counts from DB
		GenderLabels: []string{"Male", "Female", "Other"},
		GenderValues: []int{stats.MaleMembers, stats.FemaleMembers, stats.OtherGenderMembers},
	}

	data := app.newTemplateAdmin(r)
	data.Stats = *stats
	data.RecentMembers = recentMembers
	data.ChartData = chartData

	app.renderAdmin(w, r, http.StatusOK, "dashboards.gohtml", data)
}

func calculateGrowthRate(current, previous int) float64 {
	if previous == 0 {
		return 0
	}
	return float64(current-previous) / float64(previous) * 100
}

func calculateAttendanceRate(current, previous int) float64 {
	if previous == 0 {
		return 0
	}
	return float64(current-previous) / float64(previous) * 100
}

func calculatePercentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func (app *application) serviceRecordForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateAdmin(r)
	data.Form = models.CreateServiceDto{
		//Types: "second_service",
	}
	app.renderAdmin(w, r, http.StatusOK, "serviceRecord.gohtml", data)
}

func (app *application) listServiceRecords(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")
	serviceType := r.URL.Query().Get("type")
	dateFilter := r.URL.Query().Get("date")

	validateSortFields := map[string]bool{
		"date":     true,
		"type":     true,
		"males":    true,
		"females":  true,
		"offering": true,
		"tithe":    true,
	}

	if !validateSortFields[sortField] {
		sortField = "id"
	}

	if sortOrder != "desc" && sortOrder != "asc" {
		sortOrder = "asc"
	}

	records, total, err := app.recordServiceClient.GetAllServiceRecords(r.Context(), page, pageSize, sortField, sortOrder, serviceType, dateFilter)
	if err != nil {
		app.logger.Error("an error occured", err)
		app.serverError(w, r, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	pages := make([]int, 0)
	startPage := maxFunc(1, page-2)
	endPage := minFunc(totalPages, page+2)

	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}

	data := app.newTemplateAdmin(r)
	data.Records = records
	data.Pagination = Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		Pages:       pages,
		SortField:   sortField,
		SortOrder:   sortOrder,
		PageSize:    pageSize,
	}

	data.CurrentServiceType = serviceType
	data.CurrentDateFilter = dateFilter

	stats, err := app.recordServiceClient.GetServiceStatistics(r.Context(), serviceType)
	if err != nil {
		app.logger.Error("an error occured", err)
		app.serverError(w, r, err)
		return
	}

	data.ServiceStatistics = *stats
	app.renderAdmin(w, r, http.StatusOK, "records.gohtml", data)
}

func (app *application) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = app.recordServiceClient.DeleteRecord(r.Context(), id)
	if err != nil {
		app.logger.Error("an error occured", err)
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Successfully deleted record")
	http.Redirect(w, r, "/list-records", http.StatusSeeOther)
}

func (app *application) CreateServiceRecord(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	date, _ := time.Parse("2006-01-02", r.PostForm.Get("date"))

	// TODO: remember to validate the inputs like strings

	recordsData := models.CreateServiceDto{
		Types:    r.PostForm.Get("service_type"),
		Males:    parseInt(r.PostForm.Get("males")),
		Females:  parseInt(r.PostForm.Get("females")),
		Offering: parseFloat(r.PostForm.Get("offering")),
		Tithe:    parseFloat(r.PostForm.Get("tithe")),
		Date:     date,
	}

	if recordsData.Types == "friday_service" || recordsData.Types == "wednesday_service" {
		recordsData.Name = "WeekDay Service"
	} else {
		recordsData.Name = "Sunday Service"
	}

	err = app.recordServiceClient.CreateServiceRecord(r.Context(), &recordsData)
	if err != nil {
		app.logger.Error("an error occured here ", err)
		app.serverError(w, r, err)
		// TODO validate the type of errors
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Service Record Created Successfully")
	http.Redirect(w, r, "/list-records", http.StatusSeeOther)
}

func (app *application) memberDetailPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	memberData, err := app.memberClient.GetMemberById(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Member = *memberData

	app.renderAdmin(w, r, http.StatusOK, "memberDetail.gohtml", pageData)
}

func (app *application) listMembers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	// do this to prevent sql inject
	validateSortFields := map[string]bool{
		"id":          true,
		"surname":     true,
		"other_names": true,
		"dob":         true,
		"residence":   true,
	}

	if !validateSortFields[sortField] {
		sortField = "id"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	listOfMembers, totalCount, err := app.memberClient.ListMembers(r.Context(), page, pageSize, sortField, sortOrder)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	totalPages := totalCount / pageSize
	if totalCount%pageSize > 0 {
		totalPages++
	}

	var pages []int
	startPage := maxFunc(1, page-2)
	endPage := minFunc(totalPages, pageSize+2)

	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Members = listOfMembers
	pageData.Pagination = struct {
		CurrentPage int
		PageSize    int
		TotalItems  int
		TotalPages  int
		OffsetStart int
		OffsetEnd   int
		SortField   string
		SortOrder   string
		Pages       []int
	}{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
		OffsetStart: (page-1)*pageSize + 1,
		OffsetEnd:   minFunc(page*pageSize, totalCount),
		SortField:   sortField,
		SortOrder:   sortOrder,
		Pages:       getPaginationPages(page, totalPages),
	}

	app.renderAdmin(w, r, http.StatusOK, "listMembers.gohtml", pageData)
}

func (app *application) memberForm(w http.ResponseWriter, r *http.Request) {
	// Get the ID from path if it exists (for update)
	memberID := r.PathValue("id")

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	membershipYear, _ := strconv.Atoi(r.PostForm.Get("membership_year"))
	isBaptised, _ := strconv.ParseBool(r.PostForm.Get("is_baptised"))
	hasTithe, _ := strconv.ParseBool(r.PostForm.Get("has_tithe_card"))
	hasSpouse, _ := strconv.ParseBool(r.PostForm.Get("has_spouse"))
	dateOfBirth, _ := time.Parse("2006-01-02", r.PostForm.Get("dob"))
	baptismDate, _ := time.Parse("2006-01-02", r.PostForm.Get("baptism_date"))

	// limit file size to 10mb
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	file, _, _ := r.FormFile("photo_url")

	dto := models.CreateMemberDto{
		IdNumber:          r.PostForm.Get("id_number"),
		FormNumber:        r.PostForm.Get("form_number"),
		Surname:           r.PostForm.Get("surname"),
		OtherName:         r.PostForm.Get("other_name"),
		Dob:               dateOfBirth,
		Gender:            r.PostForm.Get("gender"),
		HomeTown:          r.PostForm.Get("home_town"),
		Region:            r.PostForm.Get("region"),
		Residence:         r.PostForm.Get("residency"),
		Address:           r.PostForm.Get("address"),
		Mobile:            r.PostForm.Get("mobile"),
		Email:             r.PostForm.Get("email"),
		SundaySchoolClass: r.PostForm.Get("sunday_school_class"),
		Occupation:        r.PostForm.Get("occupation"),
		HasTitheCard:      hasTithe,
		TitheCardNumber:   r.PostForm.Get("tithe_card_number"),
		DayBorn:           r.PostForm.Get("day_born"),
		HasSpouse:         hasSpouse,
		SpouseIdNumber:    r.PostForm.Get("spouse_id_number"),
		SpouseName:        r.PostForm.Get("spouse_name"),
		SpouseOccupation:  r.PostForm.Get("spouse_occupation"),
		SpouseContact:     r.PostForm.Get("spouse_contact"),
		IsBaptised:        isBaptised,
		BaptisedBy:        r.PostForm.Get("baptised_by"),
		BaptismChurch:     r.PostForm.Get("baptism_church"),
		BaptismCertNumber: r.PostForm.Get("baptism_cert_number"),
		BaptismDate:       baptismDate,
		MembershipYear:    membershipYear,
	}

	dto.CheckField(validator.NotBlank(dto.Surname), "surname", "This field cannot be blank!")
	dto.CheckField(validator.NotBlank(dto.OtherName), "other_name", "This field cannot be blank!")
	if dto.Email != "" {
		dto.CheckField(validator.Matches(dto.Email, validator.EmailRX), "email", "This field must be a valid email address")
	}
	dto.CheckField(validator.NotBlank(dto.Occupation), "occupation", "This field cannot be blank!")

	if !dto.Valid() {
		data := app.newTemplateAdmin(r)
		data.Form = dto
		if memberID != "" {
			data.IsEdit = true
		}
		app.renderAdmin(w, r, http.StatusUnprocessableEntity, "members.gohtml", data)
		return
	}

	//uploadFile
	if file != nil {
		uploadedFileURL, err := app.validateImageUploads(r, "photo_url")
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		dto.PhotoURL = uploadedFileURL
	}

	var operationResult string
	data := app.newTemplateAdmin(r)
	if memberID == "" {
		// Create new member
		dto.FormNumber = data.FormNumber
		err = app.memberClient.CreateMember(r.Context(), dto)
		operationResult = "Member created successfully!"
	} else {
		// Update existing member
		id, err := strconv.Atoi(memberID)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		err = app.memberClient.UpdateMember(r.Context(), id, dto)
		operationResult = "Member updated successfully!"
	}

	if err != nil {
		if errors.Is(err, models.PhoneInUse) {
			dto.AddFieldError("mobile", "This phone number is already in use!")
			dataAdmin := app.newTemplateAdmin(r)
			dataAdmin.Form = dto
			if memberID != "" {
				dataAdmin.IsEdit = true
			}
			app.renderAdmin(w, r, http.StatusUnprocessableEntity, "members.gohtml", dataAdmin)
		} else if dto.Email != "" && errors.Is(err, models.EmailConstraint) {
			dto.AddFieldError("email", "This email address is already in use!")
			dataAdmin := app.newTemplateAdmin(r)
			dataAdmin.Form = dto
			if memberID != "" {
				dataAdmin.IsEdit = true
			}
			app.renderAdmin(w, r, http.StatusUnprocessableEntity, "members.gohtml", dataAdmin)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", operationResult)
	http.Redirect(w, r, "/list-members", http.StatusSeeOther)
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

// serveDriveImageHandler this handler helps to render the image in the img tag
func (app *application) serveDriveImageHandler(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("fileId")
	if fileID == "" {
		http.Error(w, "Missing fileId", http.StatusBadRequest)
		return
	}

	// download from Google Drive
	resp, err := app.uploadService.Files.Get(fileID).Download()
	if err != nil {
		app.logger.Error("Failed to download image from Drive", "error", err)
		http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the full body into memory
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		app.logger.Error("Failed to read image body", "error", err)
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		return
	}

	// Detect content type
	contentType := http.DetectContentType(bodyBytes)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)

	// Stream the image to the client
	_, err = w.Write(bodyBytes)
	if err != nil {
		app.logger.Error("Failed to stream image", "error", err)
	}
}
