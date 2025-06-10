package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/story"
	"github.com/ogidi/church-media/ent/user"
	"github.com/ogidi/church-media/internal/models"
	"github.com/ogidi/church-media/internal/validator"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ChartConfig struct {
	Type    string                 `json:"type"`
	Data    map[string]interface{} `json:"data"`
	Options map[string]interface{} `json:"options"`
}

type ChartData struct {
	GrowthChart string
	AgeChart    string
	RegionChart string
	GenderChart string
}

type MessageDashboardData struct {
	Messages            []*ent.Message
	SelectedMessage     *ent.Message
	Responses           []*ent.Response
	UnreadMessagesCount int
	CurrentFilter       string
	SearchQuery         string
	User                *ent.User
}

func (app *application) getToast(r *http.Request) map[string]interface{} {
	val := app.sessionManager.Pop(r.Context(), "toast")
	if val == nil {
		return nil
	}

	toast, ok := val.(map[string]interface{})
	if !ok {
		app.logger.Error(models.InvalidToast)
		return nil
	}

	if toast["Message"] == nil {
		return nil
	}

	return toast
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateData(r)

	upcomingEvents, _ := app.eventClient.FeaturedEvents(r.Context(), 3)
	stories, _ := app.storiesClient.GetPublishedStories(r.Context(), 3)

	pageData.UpcomingEvents = upcomingEvents
	pageData.Stories = stories
	pageData.Toast = app.getToast(r)
	pageData.Form = models.CreateMessageDto{}
	app.render(w, r, http.StatusOK, "home.gohtml", pageData)
}

func (app *application) storyDetail(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateData(r)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	storyData, err := app.storiesClient.GetStoryByID(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	relatedStories, err := app.storiesClient.GetPublishedStories(r.Context(), 3)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	pageData.Story = storyData
	pageData.RelatedStories = relatedStories
	app.render(w, r, http.StatusOK, "story.gohtml", pageData)
}

func (app *application) stories(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateData(r)

	// Get page number from query params, default to 1
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Set items per page
	itemsPerPage := 9

	// Get total count of published stories
	totalStories, err := app.storiesClient.CountPublishedStories(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Calculate total pages
	totalPages := totalStories / itemsPerPage
	if totalStories%itemsPerPage != 0 {
		totalPages++
	}

	// Get paginated stories
	stories, err := app.storiesClient.GetPaginatedStories(r.Context(), page, itemsPerPage)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	storiesPagination := StoriesPagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}
	// Prepare pagination data
	pageData.Stories = stories
	pageData.StoriesPagination = storiesPagination

	app.render(w, r, http.StatusOK, "stories.gohtml", pageData)
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

func (app *application) aboutDeveloper(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}

	app.render(w, r, http.StatusOK, "aboutDeveloper.gohtml", pageData)
}

func (app *application) resetPasswordPage(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}

	app.render(w, r, http.StatusOK, "resetPassword.gohtml", pageData)
}

func (app *application) forgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "forgetPassword.gohtml", pageData)
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateAdmin(r)

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	userData, err := app.userClient.GetUserById(r.Context(), userId)
	if err != nil {
		if errors.Is(err, models.UserNotFoundError) {
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "The user you seek does not exist!",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.DashboardHome, http.StatusFound)
		}
	}
	pageData.User = *userData
	pageData.Toast = app.getToast(r)
	app.renderAdmin(w, r, http.StatusOK, "profile.gohtml", pageData)
}

func (app *application) forgetPasswordPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")

	emailExist, err := app.userClient.GetUserByEmail(r.Context(), email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if emailExist == nil {
		toastDto := map[string]interface{}{
			"Type":    "success",
			"Message": "An email has been sent to reset your password",
		}
		app.sessionManager.Put(r.Context(), "toast", toastDto)
		http.Redirect(w, r, models.LoginURL, http.StatusFound)
		return
	}

	cacheToken, err := app.userClient.GenerateToken(email, 3)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	token := fmt.Sprintf("ResetPassword%v", cacheToken)
	app.cache.Set(token, emailExist.ID, 3)

	//// send reset email with code
	//token, err := app.userClient.UpdateResetToken(r.Context(), email)
	//if err != nil {
	//	app.serverError(w, r, err)
	//	return
	//}
	registrationUrl := fmt.Sprintf("%v/reset-password?auth_token=%v", os.Getenv("APP_URL"), cacheToken)

	emailDto := EmailDto{
		Firstname:       "",
		To:              email,
		Subject:         "Reset Password",
		Token:           cacheToken,
		Expiration:      3,
		RegistrationURL: registrationUrl,
		TemplatePath:    "./ui/email-templates/resetPassword.html",
	}

	err = app.sendInvitationEmail(&emailDto, emailDto.Subject)
	if err != nil {
		app.logger.Error(models.ErrorSendingEmails, "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "An email has been sent to reset your password",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, models.LoginURL, http.StatusFound)
}

func (app *application) editAdminProfilePage(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateAdmin(r)

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	userData, err := app.userClient.GetUserById(r.Context(), userId)
	if err != nil {
		if errors.Is(err, models.UserNotFoundError) {
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "The user you seek does not exist!",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.DashboardHome, http.StatusFound)
		}
	}
	pageData.User = *userData
	app.renderAdmin(w, r, http.StatusOK, "editProfile.gohtml", pageData)
}

func (app *application) editAdminProfilePost(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	file, _, _ := r.FormFile("profile_picture")
	dateOfBirth, _ := time.Parse(models.DateFormat, r.PostForm.Get("date_of_birth"))

	dto := models.UpdateProfileRequest{
		FirstName:       r.PostForm.Get("first_name"),
		Surname:         r.PostForm.Get("surname"),
		PhoneNumber:     r.PostForm.Get("phone_number"),
		DateOfBirth:     dateOfBirth,
		Gender:          r.PostForm.Get("gender"),
		Department:      r.PostForm.Get("department"),
		Occupation:      r.PostForm.Get("occupation"),
		MaritalStatus:   r.PostForm.Get("marital_status"),
		Address:         r.PostForm.Get("address"),
		CurrentPassword: r.PostForm.Get("current_password"),
		NewPassword:     r.PostForm.Get("new_password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
	}

	dto.CheckField(validator.NotBlank(dto.FirstName), "first_name", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Surname), "surname", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Gender), "gender", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Department), "department", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Occupation), "occupation", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Address), "address", models.CannotBeBlankField)

	if !dto.Valid() {
		data := app.newTemplateAdmin(r)
		data.Form = dto
		app.renderAdmin(w, r, http.StatusUnprocessableEntity, "editProfile.gohtml", data)
		return
	}

	userData, _ := app.userClient.GetUserById(r.Context(), userId)

	// password checks
	if dto.NewPassword != "" {

		if dto.NewPassword != dto.ConfirmPassword {
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Password should be the same!",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.UserProfile, http.StatusFound)
		}

		correctPassword, err := app.userClient.PasswordCheck(r.Context(), userData.Password, dto.NewPassword)
		if err != nil {
			if errors.Is(err, models.InvalidCredentialsError) {
				toastDto := map[string]interface{}{
					"Type":    "error",
					"Message": "Password mis-matched, please try again!",
				}
				app.sessionManager.Put(r.Context(), "toast", toastDto)
				http.Redirect(w, r, models.UserProfile, http.StatusFound)
			}
		}
		if !correctPassword {
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Wrong password!",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.UserProfile, http.StatusFound)
		}
	}

	// upload image
	if file != nil {
		userData, _ := app.userClient.GetUserById(r.Context(), userId)
		if userData.Edges.ContactProfile.ProfilePicture != nil {
			err := deleteFileFromDrive(app.uploadService, *userData.Edges.ContactProfile.ProfilePicture)
			if err != nil {
				app.logger.Error("an error occured while deleting file")
				app.serverError(w, r, err)
				return
			}
		}

		uploadFile, err := app.validateImageUploads(r, "profile_picture")
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		dto.ProfilePicture = uploadFile
	}

	err = app.userClient.UpdateProfile(r.Context(), &dto, userId)
	if err != nil {
		if errors.Is(err, models.UserNotFoundError) {
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Could not find the user you seek!",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.UserProfile, http.StatusFound)
		} else if errors.Is(err, models.ConstraintError) {
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "There was a constraint error, please try again!",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.UserProfile, http.StatusFound)
		} else {
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Something bad happened, please try again",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.UserProfile, http.StatusFound)
		}

	}
	err = app.userClient.LastUpdated(r.Context(), userId) //last update made
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	userData, _ = app.userClient.GetUserById(r.Context(), userId)
	userIPAddress := r.Header.Get("X-Forwarded-For")
	if userIPAddress == "" {
		userIPAddress = r.Header.Get("X-Real-Ip")
	}
	if userIPAddress == "" {
		userIPAddress = r.RemoteAddr
	}

	entry := models.AuditLogEntry{
		Action:     "User successfully updated profile",
		EntityType: "USER",
		EntityID:   userId,
		EntityData: map[string]interface{}{
			"id":    userData.ID,
			"email": userData.Email,
			"role":  userData.Role,
			"state": userData.State,
		},
		CreatedBy: &userData.ID,
		IPAddress: userIPAddress,
		UserAgent: r.UserAgent(),
		RequestID: r.Header.Get("X-Request-Id"),
		Metadata: map[string]interface{}{
			"path":   r.URL.Path,
			"status": http.StatusText(http.StatusOK),
		},
	}
	_, err = app.logAudit.CreateLogs(r.Context(), entry)

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Successfully updated profile data!",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, "/user-profile", http.StatusSeeOther)
}

func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateData(r)
	pageData.Form = models.LoginDto{}

	var toast map[string]interface{}
	val := app.sessionManager.Pop(r.Context(), "toast")
	if val != nil {
		var ok bool
		toast, ok = val.(map[string]interface{})
		if !ok {
			app.logger.Error("toast is not of type map[string]interface{}")
			toast = nil
		}
	}

	if toast != nil && toast["Message"] != nil {
		pageData.Toast = toast
	}

	app.render(w, r, http.StatusOK, "login.gohtml", pageData)
}

func (app *application) registerPage(w http.ResponseWriter, r *http.Request) {
	registrationToken := r.URL.Query().Get("auth_token")

	var toast map[string]interface{}
	val := app.sessionManager.Pop(r.Context(), "toast")
	if val != nil {
		var ok bool
		toast, ok = val.(map[string]interface{})
		if !ok {
			app.logger.Error("toast is not of type map[string]interface{}")
			toast = nil
		}
	}

	pageData := app.newTemplateData(r)
	pageData.Form = models.RegisterDto{}
	if registrationToken != "" {
		pageData.Form = models.RegisterDto{
			RegistrationToken: registrationToken,
		}
	}

	if toast != nil && toast["Message"] != nil {
		pageData.Toast = toast
	}

	app.render(w, r, http.StatusOK, "register.gohtml", pageData)
}

func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateData(r)
	pageData.Form = models.CreateMessageDto{
		Subject: "GENERAL_ENQUIRY",
	}
	app.render(w, r, http.StatusOK, models.ContactURL, pageData)
}

func (app *application) createStoryForm(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateAdmin(r)
	pageData.Form = models.StoriesDto{}

	app.renderAdmin(w, r, http.StatusOK, "create_stories.gohtml", pageData)
}

func (app *application) aboutPastor(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "aboutPastor.gohtml", pageData)
}

func (app *application) successPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Toast = app.getToast(r)
	app.render(w, r, http.StatusOK, "verify-success.gohtml", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form models.LoginDto

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.logger.Error("an error occured while decoding")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.logger.Error(models.ErrorParsingForms, "error", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto := models.LoginDto{
		EmailUsername: r.PostForm.Get("email_username"),
		Password:      r.PostForm.Get("password"),
		RememberMe:    r.PostForm.Get("remember_me") == "on",
	}
	dto.CheckField(validator.NotBlank(dto.EmailUsername), "email_username", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Password), "password", models.CannotBeBlankField)

	if !dto.Valid() {
		app.clientError(w, http.StatusBadRequest)
		data := app.newTemplateData(r)
		data.Form = dto
		app.render(w, r, http.StatusOK, "login.gohtml", data)
		return
	}

	userFound, err := app.userClient.Login(r.Context(), &dto)
	if err != nil {
		if errors.Is(err, models.EmailUsernameInvalidError) {
			form.AddFieldError("email_username", "Invalid email address or username")
			data := app.newTemplateData(r)
			data.Form = dto
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Invalid email address or username",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.LoginURL, http.StatusFound)
		} else if errors.Is(err, models.InvalidCredentialsError) {
			data := app.newTemplateData(r)
			data.Form = dto
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Invalid email or password",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.LoginURL, http.StatusFound)
		} else if errors.Is(err, models.AccountNotVerifiedError) {
			data := app.newTemplateData(r)
			data.Form = dto
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Account not verfied yet, check your inbox and verify account!",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, models.LoginURL, http.StatusFound)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	if dto.RememberMe {
		app.sessionManager.RememberMe(r.Context(), true)
	}

	// renewing the token
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var welcomeMessage string

	if userFound.LastLogin.IsZero() {
		welcomeMessage = fmt.Sprintf("Welcome %v 😊", userFound.Edges.ContactProfile.FirstName)
	}
	welcomeMessage = fmt.Sprintf("Welcome back %v 😊", userFound.Edges.ContactProfile.FirstName)

	// audit log entry
	userIPAddress := r.Header.Get("X-Forwarded-For")
	if userIPAddress == "" {
		userIPAddress = r.Header.Get("X-Real-Ip")
	}
	if userIPAddress == "" {
		userIPAddress = r.RemoteAddr
	}

	entry := models.AuditLogEntry{
		Action:     "User successfully login",
		EntityType: "USER",
		EntityID:   userFound.ID,
		EntityData: map[string]interface{}{
			"id":    userFound.ID,
			"email": userFound.Email,
			"role":  userFound.Role,
			"state": userFound.State,
		},
		CreatedBy: &userFound.ID,
		IPAddress: userIPAddress,
		UserAgent: r.UserAgent(),
		RequestID: r.Header.Get("X-Request-Id"),
		Metadata: map[string]interface{}{
			"path":   r.URL.Path,
			"status": http.StatusText(http.StatusOK),
		},
	}
	_, err = app.logAudit.CreateLogs(r.Context(), entry)

	//add their id to the session
	app.sessionManager.Put(r.Context(), "authenticatedUserID", userFound.ID)
	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": welcomeMessage,
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, models.DashboardHome, http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.userClient.UpdateLastLogin(r.Context(), userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	//log entry
	userData, err := app.userClient.GetUserById(r.Context(), userId)
	userIPAddress := r.Header.Get("X-Forwarded-For")
	if userIPAddress == "" {
		userIPAddress = r.Header.Get("X-Real-Ip")
	}
	if userIPAddress == "" {
		userIPAddress = r.RemoteAddr
	}

	entry := models.AuditLogEntry{
		Action:     "User logout",
		EntityType: "USER",
		EntityID:   userId,
		EntityData: map[string]interface{}{
			"id":    userData.ID,
			"email": userData.Email,
			"role":  userData.Role,
			"state": userData.State,
		},
		CreatedBy: &userData.ID,
		IPAddress: userIPAddress,
		UserAgent: r.UserAgent(),
		RequestID: r.Header.Get("X-Request-Id"),
		Metadata: map[string]interface{}{
			"path":   r.URL.Path,
			"status": http.StatusText(http.StatusOK),
		},
	}
	_, err = app.logAudit.CreateLogs(r.Context(), entry)

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "User logged out successfully.",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) adminEditUserForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	userData, errr := app.userClient.GetUserById(r.Context(), id)
	if errr != nil {
		app.logger.Error("an error occured while fetching user data: ", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	pageData := app.newTemplateAdmin(r)
	pageData.User = *userData

	app.renderAdmin(w, r, http.StatusOK, "edit_role.gohtml", pageData)
}

func (app *application) adminUpdateUserRole(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	role := r.PostForm.Get("role")
	notifyUser := r.PostForm.Get("notify_user") == "on"

	err = app.userClient.UpdateUserRole(r.Context(), id, role)
	if err != nil {
		app.logger.Error("an error occured while updating user role: ", "error", err.Error())
		app.serverError(w, r, err)
		return
	}
	userData, _ := app.userClient.GetUserById(r.Context(), id)

	if notifyUser {
		emailDto := EmailDto{
			Firstname:       userData.Edges.ContactProfile.FirstName,
			To:              userData.Email,
			Subject:         "Account Role Update",
			RegistrationURL: role,
			TemplatePath:    "./ui/email-templates/role-update.html",
		}
		err = app.sendInvitationEmail(&emailDto, emailDto.Subject)
		if err != nil {
			app.logger.Error("an error occured while fetching user data: ", "error", err.Error())

			app.serverError(w, r, err)
			return
		}
	}

	userIPAddress := r.Header.Get("X-Forwarded-For")
	if userIPAddress == "" {
		userIPAddress = r.Header.Get("X-Real-Ip")
	}
	if userIPAddress == "" {
		userIPAddress = r.RemoteAddr
	}

	entry := models.AuditLogEntry{
		Action:     "User role changed!",
		EntityType: "USER",
		EntityID:   userData.ID,
		EntityData: map[string]interface{}{
			"id":    userData.ID,
			"email": userData.Email,
			"role":  userData.Role,
			"state": userData.State,
		},
		CreatedBy: &userData.ID,
		IPAddress: userIPAddress,
		UserAgent: r.UserAgent(),
		RequestID: r.Header.Get("X-Request-Id"),
		Metadata: map[string]interface{}{
			"path":   r.URL.Path,
			"status": http.StatusText(http.StatusOK),
		},
	}
	_, err = app.logAudit.CreateLogs(r.Context(), entry)

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Successfully Updated User Role",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	w.Header().Set("HX-Trigger", "roleUpdated")
	http.Redirect(w, r, "/admin-list", http.StatusSeeOther)
}

func (app *application) registerPost(w http.ResponseWriter, r *http.Request) {
	var form models.RegisterDto

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.logger.Error("an error occured while decoding")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.logger.Error("could not parse form: ", "error", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto := models.RegisterDto{
		RegistrationToken: r.PostForm.Get("token"),
		Username:          r.PostForm.Get("username"),
		Password:          r.PostForm.Get("password"),
		ConfirmPassword:   r.PostForm.Get("confirm_password"),
		Email:             r.PostForm.Get("email"),
	}
	dto.CheckField(validator.NotBlank(dto.RegistrationToken), "token", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Username), "username", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Password), "password", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Email), "email", models.CannotBeBlankField)
	dto.CheckField(validator.Matches(dto.Email, validator.EmailRX), "email", models.ValidEmail)

	if !dto.Valid() {
		app.clientError(w, http.StatusBadRequest)
		data := app.newTemplateData(r)
		data.Form = dto
		app.render(w, r, http.StatusOK, models.RegisterURL, data)
		return
	}

	isValidToken, err := app.userClient.VerifyToken(dto.RegistrationToken, dto.Email)
	if err != nil {
		if errors.Is(err, models.TokenMissMatchError) {
			dto.AddFieldError("token", "Token mismatched, please verify!")
			data := app.newTemplateData(r)
			data.Form = dto
			app.render(w, r, http.StatusOK, models.RegisterURL, data)
		} else if errors.Is(err, models.TokenExpiredError) {
			dto.AddFieldError("token", "Token expired, please reach out to support team!")
			data := app.newTemplateData(r)
			data.Form = dto
			app.render(w, r, http.StatusOK, models.RegisterURL, data)
		} else if errors.Is(err, models.TokenValidationError) {
			dto.AddFieldError("token", "Token validation failed, please try again!")
			data := app.newTemplateData(r)
			data.Form = dto
			app.render(w, r, http.StatusOK, models.RegisterURL, data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	if !isValidToken {
		dto.AddFieldError("token", "Registration token invalid, please verify!")
		data := app.newTemplateData(r)
		data.Form = dto
		app.render(w, r, http.StatusOK, models.RegisterURL, data)
	}

	err = app.userClient.Register(r.Context(), &dto)
	if err != nil {
		if errors.Is(err, models.BcryptError) {
			dto.AddFieldError("password", "an issue occured bcrypting password")
			data := app.newTemplateData(r)
			data.Form = dto
			app.render(w, r, http.StatusOK, models.RegisterURL, data)
		} else if errors.Is(err, models.UsernameExistsError) {
			dto.AddFieldError("username", "Username already exists")
			data := app.newTemplateData(r)
			data.Form = dto
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Username already exist",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			app.render(w, r, http.StatusOK, models.RegisterURL, data)
		} else if errors.Is(err, models.EmailExistsError) {
			dto.AddFieldError("email", "Email already exists")
			data := app.newTemplateData(r)
			data.Form = dto
			app.render(w, r, http.StatusOK, models.RegisterURL, data)
		} else if errors.Is(err, models.RegistrationTokenInvalidError) {
			dto.AddFieldError("token", "Registration token invalid")
			data := app.newTemplateData(r)
			data.Form = dto
			app.render(w, r, http.StatusOK, models.RegisterURL, data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	verifyToken, err := app.userClient.UpdateVerificationToken(r.Context(), dto.Email)
	if err != nil {
		app.logger.Error("an error occured while generate verification token ", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	registrationUrl := fmt.Sprintf("%v/verify?auth_token=%v", os.Getenv("APP_URL"), verifyToken)

	emailDto := EmailDto{
		Firstname:       dto.Username,
		To:              dto.Email,
		Subject:         "Account Verification",
		Token:           verifyToken,
		Expiration:      3,
		RegistrationURL: registrationUrl,
		TemplatePath:    "./ui/email-templates/verify.html",
	}
	err = app.sendInvitationEmail(&emailDto, emailDto.Subject)
	if err != nil {
		app.logger.Error("an error occured while sending email: ", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Successfully Registered Administrative Account",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)

	requestedLoginPath := app.sessionManager.PopString(r.Context(), "requestedPathAfterLogin")
	if requestedLoginPath != "" {
		http.Redirect(w, r, requestedLoginPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) activateAccount(w http.ResponseWriter, r *http.Request) {
	registrationToken := r.URL.Query().Get("auth_token")
	if registrationToken == "" {
		app.sessionManager.Put(r.Context(), "toast", map[string]interface{}{
			"Type":    "error",
			"Message": "Invalid activation link",
		})
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	registrationToken = fmt.Sprintf("ResetPassword%v", registrationToken)
	token, exist := app.cache.Get(registrationToken)
	if !exist {
		app.sessionManager.Put(r.Context(), "toast", map[string]interface{}{
			"Type":    "error",
			"Message": "Invalid or expired activation token",
		})
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	id, _ := strconv.Atoi(token.(string))
	userExistById, _ := app.userClient.GetUserById(r.Context(), id)
	if userExistById == nil {
		app.sessionManager.Put(r.Context(), "toast", map[string]interface{}{
			"Type":    "error",
			"Message": "Account activation failed",
		})
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	err := app.userClient.AccountActivation(r.Context(), userExistById.ID)
	if err != nil {
		app.sessionManager.Put(r.Context(), "toast", map[string]interface{}{
			"Type":    "error",
			"Message": "Account activation failed",
		})
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	app.sessionManager.Put(r.Context(), "toast", map[string]interface{}{
		"Type":    "success",
		"Message": "Account activated successfully! You can now login",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) adminInvitePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("could not parse form: ", "error", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, _ := strconv.Atoi(r.PostForm.Get("expires_at"))
	dto := models.InviteDto{
		Firstname: r.PostForm.Get("firstname"),
		Lastname:  r.PostForm.Get("lastname"),
		Email:     r.PostForm.Get("email"),
		Role:      user.Role(r.PostForm.Get("role")),
		ExpiresAt: expires,
	}

	err = app.decodePostForm(r, &dto)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto.CheckField(validator.NotBlank(dto.Firstname), "firstname", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Lastname), "lastname", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Email), "email", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(string(dto.Role)), "role", models.CannotBeBlankField)
	dto.CheckField(validator.Matches(dto.Email, validator.EmailRX), "email", models.ValidEmail)

	if !dto.Valid() {
		data := app.newTemplateAdmin(r)
		data.Form = dto
		app.renderAdmin(w, r, http.StatusUnprocessableEntity, models.InviteURL, data)
		return
	}

	token, firstname, err := app.userClient.InviteUser(r.Context(), &dto)
	if err != nil {
		app.logger.Error(models.AnErrorOccured, "error", err.Error())
		if errors.Is(err, models.EmailExistsError) {
			dto.AddFieldError("email", "Email already exists")
			data := app.newTemplateAdmin(r)
			data.Form = dto
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Email already exist",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			app.renderAdmin(w, r, http.StatusUnprocessableEntity, models.InviteURL, data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	registrationUrl := fmt.Sprintf("%v/register?auth_token=%v", os.Getenv("APP_URL"), token)

	emailDto := EmailDto{
		Firstname:       firstname,
		To:              dto.Email,
		Subject:         "Administrative Account Registration",
		Token:           token,
		Expiration:      dto.ExpiresAt,
		RegistrationURL: registrationUrl,
		TemplatePath:    "./ui/email-templates/registration-token.html",
	}

	// send email here
	err = app.sendInvitationEmail(&emailDto, emailDto.Subject)
	if err != nil {
		app.logger.Error("an error occured while sending email: ", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Successfully invited church administrator",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, models.DashboardHome, http.StatusSeeOther)
}

func (app *application) subscribe(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	err = app.subscribeClient.Subscribe(r.Context(), email)
	if err != nil {
		if errors.Is(err, models.EmailExistsError) {
			toastDto := map[string]interface{}{
				"Type":    "error",
				"Message": "Email already subscribed to our newsletter!",
			}
			app.sessionManager.Put(r.Context(), "toast", toastDto)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			app.logger.Error(err.Error())
			app.serverError(w, r, err)
			return
		}
	}

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Successfully subscribed to our newsletter!",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	dto.CheckField(validator.NotBlank(dto.Name), "name", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Email), "email", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Subject), "subject", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Description), "description", models.CannotBeBlankField)
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

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Message successfully sent!",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) handleLike(w http.ResponseWriter, r *http.Request) {
	storyID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check cookie to prevent duplicate votes
	cookieName := fmt.Sprintf("story_%d_vote", storyID)
	if _, err := r.Cookie(cookieName); err == nil {
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}

	// Increment like count only
	likes, err := app.storiesClient.IncrementLike(r.Context(), storyID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Get current dislikes count
	story, err := app.storiesClient.GetStoryByID(r.Context(), storyID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Set cookie to prevent duplicate votes (expires in 30 days)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "liked",
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"likes":    likes,
		"dislikes": story.Dislikes,
	})
}

func (app *application) handleDisLike(w http.ResponseWriter, r *http.Request) {
	storyID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check cookie to prevent duplicate votes
	cookieName := fmt.Sprintf("story_%d_vote", storyID)
	if _, err := r.Cookie(cookieName); err == nil {
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}

	// Increment dislike count only
	dislikes, err := app.storiesClient.IncrementDislike(r.Context(), storyID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Get current likes count
	story, err := app.storiesClient.GetStoryByID(r.Context(), storyID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Set cookie to prevent duplicate votes (expires in 30 days)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "disliked", // Changed to "disliked" for distinction
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"likes":    story.Likes,
		"dislikes": dislikes,
	})
}

func (app *application) giveOnline(w http.ResponseWriter, r *http.Request) {
	pageData := templateData{}
	app.render(w, r, http.StatusOK, "giveOnline.gohtml", pageData)
}

// pagesAdmin interfaces
func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "all"
	}

	messages, err := app.messageClient.GetMessages(r.Context(), filter)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	unreadCount, err := app.messageClient.GetUnreadCount(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	respondedCount, err := app.messageClient.GetRespondedCount(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	prayerCount, err := app.messageClient.GetPrayerCount(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Messages = messages
	pageData.UnreadMessagesCount = unreadCount
	pageData.RespondedCount = respondedCount
	pageData.PrayerCount = prayerCount
	pageData.CurrentFilter = filter

	app.renderAdmin(w, r, http.StatusOK, "admin.gohtml", pageData)
}

func (app *application) viewMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	msg, err := app.messageClient.GetMessageWithResponses(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := struct {
		SelectedMessage *ent.Message
		Responses       []*ent.Response
	}{
		SelectedMessage: msg,
		Responses:       msg.Edges.Responses,
	}

	pageData := app.newTemplateAdmin(r)
	pageData.SelectedMessage = data.SelectedMessage
	pageData.Responses = data.Responses

	// Render just the message-view partial
	err = app.templateCacheAdmin["partialsAdmin/message_view.gohtml"].Execute(w, pageData)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) replyToMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	subject := r.FormValue("subject")
	description := r.FormValue("description")
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	_, err = app.messageClient.CreateResponse(r.Context(), id, userId, email, subject, description)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.messageClient.MarkAsResponded(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	msg, err := app.messageClient.GetMessageWithResponses(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := struct {
		SelectedMessage *ent.Message
		Responses       []*ent.Response
	}{
		SelectedMessage: msg,
		Responses:       msg.Edges.Responses,
	}

	pageData := app.newTemplateAdmin(r)
	pageData.SelectedMessage = data.SelectedMessage
	pageData.Responses = data.Responses

	app.renderAdmin(w, r, http.StatusOK, "message_view.gohtml", pageData)
}

func (app *application) filterMessages(w http.ResponseWriter, r *http.Request) {
	filter := r.PathValue("filter")
	if filter == "" {
		filter = "all"
	}

	messages, err := app.messageClient.GetMessages(r.Context(), filter)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := struct {
		Messages []*ent.Message
	}{
		Messages: messages,
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Messages = data.Messages

	app.renderAdmin(w, r, http.StatusOK, models.MessageList, pageData)
}

func (app *application) searchMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")

	messages, err := app.messageClient.SearchMessages(r.Context(), filter, query)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := struct {
		Messages []*ent.Message
	}{
		Messages: messages,
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Messages = data.Messages

	app.renderAdmin(w, r, http.StatusOK, models.MessageList, pageData)
}

func (app *application) createStory(w http.ResponseWriter, r *http.Request) {
	var form models.StoriesDto

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.logger.Error("an error occured while decoding")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.logger.Error("could not parse form: ", "error", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto := models.StoriesDto{
		Title:   r.PostForm.Get("title"),
		Body:    r.PostForm.Get("body"),
		Excerpt: r.PostForm.Get("excerpt"),
		Status:  story.Status(r.PostForm.Get("status")),
	}
	dto.CheckField(validator.NotBlank(dto.Title), "title", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Body), "body", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.Excerpt), "excerpt", models.CannotBeBlankField)

	if !dto.Valid() {
		data := app.newTemplateData(r)
		data.Form = dto
		app.render(w, r, http.StatusUnprocessableEntity, "create_stories.gohtml", data)
		return
	}

	file, _, _ := r.FormFile("imageUrl")

	if file != nil {
		uploadImage, err := app.validateImageUploads(r, "imageUrl")
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		dto.Image = uploadImage
	}

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	dto.AuthorId = userId

	err = app.storiesClient.CreateStory(r.Context(), dto)
	if err != nil {
		app.logger.Error("could not create story: ", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Story successfully created!",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) markAsRead(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.messageClient.MarkAsRead(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Return updated message list
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "all"
	}

	messages, err := app.messageClient.GetMessages(r.Context(), filter)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := struct {
		Messages []*ent.Message
	}{
		Messages: messages,
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Messages = data.Messages

	app.renderAdmin(w, r, http.StatusOK, models.MessageList, pageData)
}

func (app *application) deleteMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.messageClient.DeleteMessage(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Return updated message list
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "all"
	}

	messages, err := app.messageClient.GetMessages(r.Context(), filter)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := struct {
		Messages []*ent.Message
	}{
		Messages: messages,
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Messages = data.Messages

	app.renderAdmin(w, r, http.StatusOK, "message_list.gohtml", pageData)
}

func (app *application) getUnreadCount(w http.ResponseWriter, r *http.Request) {
	count, err := app.messageClient.GetUnreadCount(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}

func (app *application) updateMessageState(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	state := r.PathValue("state")
	if state == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.messageClient.UpdateState(r.Context(), id, state)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) getMessageResponses(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	responses, err := app.messageClient.GetResponses(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := struct {
		Responses []*ent.Response
	}{
		Responses: responses,
	}

	pageData := app.newTemplateAdmin(r)
	pageData.Responses = data.Responses

	app.renderAdmin(w, r, http.StatusOK, "message_list.gohtml", pageData)
}

func (app *application) adminList(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateAdmin(r)
	pageData.Toast = app.getToast(r)

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize := 10 // You can make this configurable

	filters := models.Filters{
		Role:   r.URL.Query().Get("role"),
		Status: r.URL.Query().Get("status"),
		Search: r.URL.Query().Get("search"),
	}

	sort := models.Sort{
		Field: r.URL.Query().Get("sort"),
		Order: r.URL.Query().Get("order"),
	}
	if sort.Order != "asc" && sort.Order != "desc" {
		sort.Order = "desc" // default order
	}

	// Get invitations with filters, sorting, and pagination
	admins, pagination, err := app.userClient.GetAllInvitedAdmins(r.Context(), page, pageSize, filters, sort)
	if err != nil {
		app.logger.Error("an error occurred while fetching admins: ", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	pageData.Users = admins
	pageData.Filters = filters
	pageData.Sort = sort
	pageData.Pagination = Pagination{
		CurrentPage: pagination.CurrentPage,
		PageSize:    pagination.PageSize,
		TotalPages:  pagination.TotalPages,
		TotalItems:  pagination.TotalItems,
	}
	pageData.Toast = app.getToast(r)
	app.renderAdmin(w, r, http.StatusOK, "adminsList.gohtml", pageData)

}

func (app *application) adminInvites(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateAdmin(r)
	pageData.Form = models.InviteDto{}
	pageData.Toast = app.getToast(r)

	app.renderAdmin(w, r, http.StatusOK, models.InviteURL, pageData)
}

func (app *application) members(w http.ResponseWriter, r *http.Request) {
	pageData := app.newTemplateAdmin(r)
	pageData.IsEdit = false
	pageData.Form = models.CreateMemberDto{
		HasTitheCard: true,
	}

	app.renderAdmin(w, r, http.StatusOK, models.MembersURL, pageData)
}

func (app *application) editMember(w http.ResponseWriter, r *http.Request) {
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
	pageData.Form = models.CreateMemberDto{}
	pageData.Member = *memberData

	app.renderAdmin(w, r, http.StatusOK, models.MemberEdits, pageData)
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

	app.renderAdmin(w, r, http.StatusOK, models.MembersURL, pageData)
}

func (app *application) deleteMember(w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE method
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	msg := fmt.Sprintf("Successfully deleted %v", member.OtherNames)
	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": msg,
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)

	// Return a JSON response with redirect URL
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"redirect": "/list-members",
	})
}

func (app *application) deleteAdminAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = app.userClient.DeleteUser(r.Context(), id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Successfully deleted church administrator",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, "/admin-list", http.StatusSeeOther)
}

func (app *application) dashboardData(w http.ResponseWriter, r *http.Request) {
	// Get filter parameters
	timeFilter := r.URL.Query().Get("timeFilter")
	if timeFilter == "" {
		timeFilter = "365" // default to "This Year"
	}
	growthChartType := r.URL.Query().Get("growthChartType")
	if growthChartType == "" {
		growthChartType = "yearly" // default to yearly
	}

	// Get statistics
	stats, err := app.memberClient.GetMemberStatistics(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Calculate derived stats
	stats.GrowthRate = calculateGrowthRate(stats.TotalMembers, stats.PreviousYearMembers)
	stats.BaptismPercentage = calculatePercentage(stats.BaptizedMembers, stats.TotalMembers)
	stats.AttendanceRate = calculateAttendanceRate(stats.ActiveMembers, stats.TotalMembers)
	stats.MalePercentage = calculatePercentage(stats.MaleMembers, stats.TotalMembers)
	stats.FemalePercentage = calculatePercentage(stats.FemaleMembers, stats.TotalMembers)

	// Get birthdays
	stats.BirthdaysThisMonth, err = app.memberClient.GetBirthdaysThisMonth(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Get recent members
	recentMembers, err := app.memberClient.GetRecentMembers(r.Context(), 5)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Prepare chart data
	chartData := ChartData{
		GrowthChart: app.createGrowthChart(r, growthChartType),
		AgeChart:    app.createAgeChart(r),
		RegionChart: app.createRegionChart(r),
		GenderChart: app.createGenderChart(stats.MaleMembers, stats.FemaleMembers, stats.OtherGenderMembers),
	}

	// Prepare template data
	data := app.newTemplateAdmin(r)
	data.Stats = *stats
	data.RecentMembers = recentMembers
	data.ChartData = chartData
	data.Toast = app.getToast(r)

	// Render template
	app.renderAdmin(w, r, http.StatusOK, "dashboards.gohtml", data)
}

func (app *application) createGrowthChart(r *http.Request, chartType string) string {
	var labels []string
	var data []int
	var err error

	// Get actual data from database
	if chartType == "monthly" {
		labels, data, err = app.memberClient.GetMonthlyGrowth2(r.Context())
		if err != nil {
			app.logger.Error("failed to get monthly growth data", "error", err.Error())
			labels = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
			data = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		}
	} else {
		labels, data, err = app.memberClient.GetYearlyGrowth(r.Context())
		if err != nil {
			app.logger.Error("failed to get yearly growth data", "error", err.Error())
			currentYear := time.Now().Year()
			labels = []string{
				strconv.Itoa(currentYear - 3),
				strconv.Itoa(currentYear - 2),
				strconv.Itoa(currentYear - 1),
				strconv.Itoa(currentYear),
			}
			data = []int{0, 0, 0, 0}
		}
	}

	// Determine x-axis title based on chart type
	xAxisTitle := "Years"
	if chartType == "monthly" {
		xAxisTitle = "Months"
	}

	chartConfig := ChartConfig{
		Type: "line",
		Data: map[string]interface{}{
			"labels": labels,
			"datasets": []map[string]interface{}{
				{
					"label":           "Membership Growth",
					"data":            data,
					"borderColor":     models.ColorCode,
					"backgroundColor": "rgba(54, 162, 235, 0.1)",
					"tension":         0.4,
					"fill":            true,
					"borderWidth":     2,
				},
			},
		},
		Options: map[string]interface{}{
			"responsive":          true,
			"maintainAspectRatio": false,
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
					"title": map[string]interface{}{
						"display": true,
						"text":    "Number of Members",
					},
				},
				"x": map[string]interface{}{
					"title": map[string]interface{}{
						"display": true,
						"text":    xAxisTitle,
					},
				},
			},
			"plugins": map[string]interface{}{
				"tooltip": map[string]interface{}{
					"callbacks": map[string]interface{}{
						"label": `function(context) {
                            return 'Members: ' + context.parsed.y;
                        }`,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(chartConfig)
	if err != nil {
		app.logger.Error("failed to marshal growth chart config", "error", err.Error())
		return "{}"
	}
	return string(jsonData)
}

func (app *application) createAgeChart(r *http.Request) string {
	// Get actual data from database
	labels, data, err := app.memberClient.GetMemberAgeDistribution(r.Context())
	if err != nil {
		app.logger.Error("failed to get age distribution", "error", err.Error())
		// Fallback to empty data
		labels = []string{"0-18", "19-35", "36-50", "51-65", "65+"}
		data = []int{0, 0, 0, 0, 0}
	}

	config := ChartConfig{
		Type: "pie",
		Data: map[string]interface{}{
			"labels": labels,
			"datasets": []map[string]interface{}{
				{
					"label":           "Age Distribution",
					"data":            data,
					"backgroundColor": []string{"#ff6384", "#36a2eb", "#ffce56", "#4bc0c0", "#9966ff"},
					"borderWidth":     1,
				},
			},
		},
		Options: map[string]interface{}{
			"responsive":          true,
			"maintainAspectRatio": false,
			"plugins": map[string]interface{}{
				"legend": map[string]interface{}{
					"position": "right",
				},
				"tooltip": map[string]interface{}{
					"callbacks": map[string]interface{}{
						"label": "function(context) {\n" +
							"    let label = context.label || '';\n" +
							"    let value = context.raw || 0;\n" +
							"    let total = context.dataset.data.reduce((a, b) => a + b, 0);\n" +
							"    let percentage = Math.round((value / total) * 100);\n" +
							"    return `${label}: ${value} (${percentage}%)`;\n" +
							"}",
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		app.logger.Error("failed to marshal chart config", "error", err.Error())
		return "{}"
	}
	return string(jsonData)
}

func (app *application) createRegionChart(r *http.Request) string {
	regions, counts, err := app.memberClient.GetDistinctRegions(r.Context())
	if err != nil {
		// Fallback to empty data if there's an error
		regions = []string{}
		counts = []int{}
		app.logger.Error("failed to get region data", "error", err.Error())
	}

	config := ChartConfig{
		Type: "bar",
		Data: map[string]interface{}{
			"labels": regions,
			"datasets": []map[string]interface{}{
				{
					"label":           "Members by Region",
					"data":            counts,
					"backgroundColor": "#4bc0c0",
				},
			},
		},
		Options: map[string]interface{}{
			"responsive":          true,
			"maintainAspectRatio": false,
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
				},
			},
		},
	}

	jsonData, _ := json.Marshal(config)
	return string(jsonData)
}

func (app *application) createGenderChart(male, female, other int) string {
	labels := []string{"Male", "Female", "Other"}
	data := []int{male, female, other}

	config := ChartConfig{
		Type: "doughnut",
		Data: map[string]interface{}{
			"labels": labels,
			"datasets": []map[string]interface{}{
				{
					"label": "Gender Distribution",
					"data":  data,
					"backgroundColor": []string{
						"#36a2eb", "#ff6384", "#ffce56",
					},
				},
			},
		},
		Options: map[string]interface{}{
			"responsive":          true,
			"maintainAspectRatio": false,
		},
	}

	jsonData, _ := json.Marshal(config)
	return string(jsonData)
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

func (app *application) editEventForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	data := app.newTemplateAdmin(r)

	eventData, err := app.eventClient.GetEventByID(r.Context(), id)
	if err != nil {
		app.logger.Error("an error occured when fetching event")
		app.serverError(w, r, err)
		return
	}
	data.Event = *eventData

	app.renderAdmin(w, r, http.StatusOK, "edit_event.gohtml", data)
}

func (app *application) churchEvents(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateAdmin(r)
	data.Form = models.EventModel{}

	allEvents, err := app.eventClient.AllEvents(r.Context())
	if err != nil {
		app.logger.Error("upcoming events: %v", "error", err.Error())
		app.serverError(w, r, err)
		return
	}
	data.Events = allEvents
	app.renderAdmin(w, r, http.StatusOK, "events.gohtml", data)
}

func (app *application) editChurchEvent(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("could not parse form: %v", "error", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// limit file size to 10mb
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	pageData := app.newTemplateAdmin(r)
	pageData.Toast = app.getToast(r)

	eventDto := models.CreateEventDto{
		Title:       r.PostForm.Get("title"),
		Description: r.PostForm.Get("description"),
		StartTime:   r.PostForm.Get("startTime"),
		EndTime:     r.PostForm.Get("endTime"),
		Location:    r.PostForm.Get("location"),
		Featured:    r.PostForm.Get("featured") == "on",
	}
	eventDto.CheckField(validator.NotBlank(eventDto.Title), "title", "This field is required")
	file, _, _ := r.FormFile("image_url")

	if !eventDto.Valid() {
		toastDto := map[string]interface{}{
			"Type":    "error",
			"Message": "An error occured while updating the event.",
		}
		app.sessionManager.Put(r.Context(), "toast", toastDto)
		pageData.Form = eventDto
		app.renderAdmin(w, r, http.StatusBadRequest, "edit_event.gohtml", pageData)
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	if file != nil {
		uploadImage, err := app.validateImageUploads(r, "image_url")
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		eventDto.ImageUrl = uploadImage
		//existingEvent, _ := app.eventClient.GetEventByID(r.Context(), id)
		//err = deleteFileFromDrive(app.uploadService, existingEvent.ImageURL)
		//if err != nil {
		//	app.logger.Error("could not update event", err.Error())
		//	toastDto := map[string]interface{}{
		//		"Type":    "error",
		//		"Message": "Could not delete and update existing image URL!",
		//	}
		//	app.sessionManager.Put(r.Context(), "toast", toastDto)
		//	pageData.Form = eventDto
		//	app.renderAdmin(w, r, http.StatusBadRequest, "edit_event.gohtml", pageData)
		//}

	}

	err = app.eventClient.EditEvent(r.Context(), &eventDto, id)
	if err != nil {
		app.logger.Error("could not edit event: %v", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Successfully updated event data",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, models.DashboardHome, http.StatusSeeOther)
}

func (app *application) churchEventForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("could not parse form: %v", "error", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// limit file size to 10mb
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	eventDto := models.CreateEventDto{
		Title:       r.PostForm.Get("title"),
		Description: r.PostForm.Get("description"),
		StartTime:   r.PostForm.Get("startTime"),
		EndTime:     r.PostForm.Get("endTime"),
		Location:    r.PostForm.Get("location"),
		Featured:    r.PostForm.Get("featured") == "on",
	}
	eventDto.CheckField(validator.NotBlank(eventDto.Title), "title", "This field is required")

	file, _, _ := r.FormFile("imageUrl")

	// validate and upload file
	if file != nil {
		uploadImage, err := app.validateImageUploads(r, "imageUrl")
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		eventDto.ImageUrl = uploadImage
	}

	err = app.eventClient.CreateEvent(r.Context(), &eventDto)
	if err != nil {
		app.logger.Error("could not create event: %v", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Successfully created event!",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, "/dashboards", http.StatusSeeOther)
}

func (app *application) viewEventDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	data := app.newTemplateData(r)

	eventDetails, err := app.eventClient.ViewEvent(r.Context(), id)
	if err != nil {
		app.logger.Error("viewEventDetails: error %v", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	data.Event = eventDetails
	app.render(w, r, http.StatusOK, "viewEventDetails.gohtml", data)
}

func (app *application) listEvents(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 9

	eventList, total, err := app.eventClient.ListEvents(r.Context(), page, pageSize)
	if err != nil {
		app.logger.Error("an error occured %v", "error", err.Error())
		app.serverError(w, r, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	pages := make([]int, 0)
	startPage := max(1, page-2)
	endPage := min(totalPages, page+2)
	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}
	data.Events = eventList
	data.Pagination = struct {
		CurrentPage int
		TotalPages  int
		Pages       []int
	}{CurrentPage: page, TotalPages: totalPages, Pages: pages}

	app.render(w, r, http.StatusOK, "listEvents.gohtml", data)
}

func (app *application) deleteEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	existingEvent, err := app.eventClient.EventExist(r.Context(), id)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverError(w, r, err)
		return
	}

	if existingEvent.ImageURL != "" {
		err = deleteFileFromDrive(app.uploadService, existingEvent.ImageURL)
		if err != nil {
			app.logger.Error(err.Error())
			app.serverError(w, r, err)
			return
		}
	}

	err = app.eventClient.DeleteEvent(r.Context(), id)
	if err != nil {
		app.logger.Error("could not delete event: %v", "error", err.Error())
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Event deleted successfully")
	http.Redirect(w, r, "/dashboards", http.StatusSeeOther)
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
		app.logger.Error(models.AnErrorOccured, "error", err.Error())
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
		app.logger.Error("an error occured", "error", err.Error())
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
		app.logger.Error("an error occured", "error", err.Error())
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
	date, _ := time.Parse(models.DateFormat, r.PostForm.Get("date"))

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
		app.logger.Error("an error occured here ", "error", err.Error())
		app.serverError(w, r, err)
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

	pageData.Toast = app.getToast(r)
	app.renderAdmin(w, r, http.StatusOK, "listMembers.gohtml", pageData)
}

func (app *application) editMemberForm(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	id, err := strconv.Atoi(memberID)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
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

	dto.CheckField(validator.NotBlank(dto.Surname), "surname", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.OtherName), "other_name", models.CannotBeBlankField)
	if dto.Email != "" {
		dto.CheckField(validator.Matches(dto.Email, validator.EmailRX), "email", "This field must be a valid email address")
	}
	dto.CheckField(validator.NotBlank(dto.Occupation), "occupation", models.CannotBeBlankField)

	if !dto.Valid() {
		data := app.newTemplateAdmin(r)
		data.Form = dto
		app.renderAdmin(w, r, http.StatusUnprocessableEntity, models.MemberEdits, data)
		return
	}

	if file != nil {
		uploadedFileURL, err := app.validateImageUploads(r, "photo_url")
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		dto.PhotoURL = uploadedFileURL

	}

	err = app.memberClient.UpdateMember(r.Context(), id, dto)

	if err != nil {
		if errors.Is(err, models.PhoneInUse) {
			dto.AddFieldError("mobile", "Phone in use!")
			dataAdmin := app.newTemplateAdmin(r)
			dataAdmin.Form = dto
			app.renderAdmin(w, r, http.StatusUnprocessableEntity, models.MemberEdits, dataAdmin)
			return
		} else if dto.Email != "" && errors.Is(err, models.EmailConstraint) {
			dto.AddFieldError("email", "Email in use!")
			dataAdmin := app.newTemplateAdmin(r)
			dataAdmin.Form = dto
			app.renderAdmin(w, r, http.StatusUnprocessableEntity, models.MemberEdits, dataAdmin)
			return
		} else {
			app.serverError(w, r, err)
			return
		}
	}
	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "successfully updated member data!",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
	http.Redirect(w, r, "/list-members", http.StatusSeeOther)
}

func (app *application) memberForm(w http.ResponseWriter, r *http.Request) {

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

	dto.CheckField(validator.NotBlank(dto.Surname), "surname", models.CannotBeBlankField)
	dto.CheckField(validator.NotBlank(dto.OtherName), "other_name", models.CannotBeBlankField)
	if dto.Email != "" {
		dto.CheckField(validator.Matches(dto.Email, validator.EmailRX), "email", "This field must be a valid email address")
	}
	dto.CheckField(validator.NotBlank(dto.Occupation), "occupation", models.CannotBeBlankField)

	if !dto.Valid() {
		data := app.newTemplateAdmin(r)
		data.Form = dto
		app.renderAdmin(w, r, http.StatusUnprocessableEntity, models.MembersURL, data)
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

	data := app.newTemplateAdmin(r)

	dto.FormNumber = data.FormNumber
	err = app.memberClient.CreateMember(r.Context(), dto)

	if err != nil {
		if errors.Is(err, models.PhoneInUse) {
			dto.AddFieldError("mobile", "This phone number is already in use!")
			dataAdmin := app.newTemplateAdmin(r)
			dataAdmin.Form = dto
			app.renderAdmin(w, r, http.StatusUnprocessableEntity, models.MembersURL, dataAdmin)
		} else if dto.Email != "" && errors.Is(err, models.EmailConstraint) {
			dto.AddFieldError("email", "This email address is already in use!")
			dataAdmin := app.newTemplateAdmin(r)
			dataAdmin.Form = dto

			app.renderAdmin(w, r, http.StatusUnprocessableEntity, models.MembersURL, dataAdmin)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	toastDto := map[string]interface{}{
		"Type":    "success",
		"Message": "Member data succesfully created!",
	}
	app.sessionManager.Put(r.Context(), "toast", toastDto)
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

// serveDriveImageHandler this handler helps to render the image within the img tag
func (app *application) serveDriveImageHandler(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("fileId")
	if fileID == "" {
		http.Error(w, "Missing fileId", http.StatusBadRequest)
		return
	}

	// download from Google Drive
	resp, err := app.uploadService.Files.Get(fileID).Download()
	if err != nil {
		app.logger.Error("Failed to download image from Drive", "error", err.Error())
		http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the full body into memory
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		app.logger.Error("Failed to read image body", "error", err.Error())
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
		app.logger.Error("Failed to stream image", "error", err.Error())
	}
}
