package main

import (
	"github.com/justinas/alice"
	"github.com/ogidi/church-media/ui"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// routes here
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// middlewares
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	//protected routes initialization
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /", dynamic.ThenFunc(app.home))
	mux.Handle("GET /ministries", dynamic.ThenFunc(app.ministries))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /events", dynamic.ThenFunc(app.events))
	mux.Handle("GET /sermons", dynamic.ThenFunc(app.sermon))
	mux.Handle("GET /contact", dynamic.ThenFunc(app.contact))
	mux.Handle("GET /about-pastor", dynamic.ThenFunc(app.aboutPastor))
	mux.Handle("POST /contact", dynamic.ThenFunc(app.contactForm))
	mux.Handle("GET /give-online", dynamic.ThenFunc(app.giveOnline))
	mux.Handle("GET /dashboard", protected.ThenFunc(app.dashboard))
	mux.Handle("GET /admin/messages/filter-by/{filter}", protected.ThenFunc(app.filterMessages))
	mux.Handle("GET /admin/messages/search", protected.ThenFunc(app.searchMessages))
	mux.Handle("POST /admin/messages/{id}/read", protected.ThenFunc(app.markAsRead))
	mux.Handle("DELETE /admin/messages/{id}", protected.ThenFunc(app.deleteMessage))
	mux.Handle("GET /admin/messages/unread-count", protected.ThenFunc(app.getUnreadCount))
	mux.Handle("PUT /admin/messages/{id}/state/{state}", protected.ThenFunc(app.updateMessageState))
	mux.Handle("GET /admin/messages/details/{id}/responses", protected.ThenFunc(app.getMessageResponses))
	mux.Handle("GET /dashboards", protected.ThenFunc(app.dashboardData))
	mux.Handle("POST /member", protected.ThenFunc(app.memberForm))
	mux.Handle("DELETE /members/{id}/delete", protected.ThenFunc(app.deleteMember))
	mux.Handle("POST /member/{id}", protected.ThenFunc(app.memberForm))
	mux.Handle("GET /member/success/{id}", protected.ThenFunc(app.memberSuccess))
	mux.Handle("GET /member", protected.ThenFunc(app.members))
	mux.Handle("GET /member/{id}", protected.ThenFunc(app.membersEdit))
	mux.Handle("GET /list-members", protected.ThenFunc(app.listMembers))
	mux.Handle("GET /members/{id}/view", protected.ThenFunc(app.memberDetailPage))
	mux.Handle("GET /mens", dynamic.ThenFunc(app.mensMinistry))
	mux.Handle("GET /womens", dynamic.ThenFunc(app.womenMinistry))
	mux.Handle("GET /youth", dynamic.ThenFunc(app.youthMinistry))
	mux.Handle("GET /childrens", dynamic.ThenFunc(app.childrenMinistry))
	mux.Handle("GET /teens", dynamic.ThenFunc(app.teensMinistry))
	mux.Handle("GET /music", dynamic.ThenFunc(app.musicMinistry))
	mux.Handle("GET /ushers", dynamic.ThenFunc(app.ushersMinistry))
	mux.Handle("GET /service-record", protected.ThenFunc(app.serviceRecordForm))
	mux.Handle("POST /service-record", protected.ThenFunc(app.CreateServiceRecord))
	mux.Handle("GET /list-records", protected.ThenFunc(app.listServiceRecords))
	mux.Handle("DELETE /service-record/{id}/delete", protected.ThenFunc(app.DeleteRecord))
	mux.Handle("GET /church-events", protected.ThenFunc(app.churchEvents))
	mux.Handle("POST /event-create", protected.ThenFunc(app.churchEventForm))
	mux.Handle("POST /event/{id}/delete", protected.ThenFunc(app.deleteEvent))
	mux.Handle("GET /event/{id}/view", dynamic.ThenFunc(app.viewEventDetails))
	mux.Handle("GET /events-list", dynamic.ThenFunc(app.listEvents))
	mux.Handle("GET /about-developer", dynamic.ThenFunc(app.aboutDeveloper))
	mux.Handle("GET /register", dynamic.ThenFunc(app.registerPage))
	mux.Handle("POST /register", dynamic.ThenFunc(app.registerPost))
	mux.Handle("GET /invite", protected.ThenFunc(app.adminInvites))
	mux.Handle("GET /admin-list", protected.ThenFunc(app.adminList))
	mux.Handle("POST /invite", protected.ThenFunc(app.adminInvitePost))
	mux.Handle("GET /verify", dynamic.ThenFunc(app.activateAccount))
	mux.Handle("GET /success", dynamic.ThenFunc(app.successPage))
	mux.Handle("GET /login", dynamic.ThenFunc(app.loginPage))
	mux.Handle("POST /login", dynamic.ThenFunc(app.loginPost))
	mux.Handle("POST /logout", protected.ThenFunc(app.logout))
	mux.Handle("GET /user-profile", protected.ThenFunc(app.userProfile))
	mux.Handle("GET /user-profile/edit", protected.ThenFunc(app.editAdminProfilePage))
	mux.Handle("POST /user-profile/edit", protected.ThenFunc(app.editAdminProfilePost))
	mux.Handle("GET /reset-password", dynamic.ThenFunc(app.resetPasswordPage))
	mux.Handle("GET /forget-password", dynamic.ThenFunc(app.forgotPasswordPage))
	mux.Handle("POST /forget-password", dynamic.ThenFunc(app.forgetPasswordPost))
	mux.Handle("GET /admin/users/{id}/edit", protected.ThenFunc(app.adminEditUserForm))
	mux.Handle("POST /admin/users/{id}", protected.ThenFunc(app.adminUpdateUserRole))
	mux.Handle("POST /delete/{id}", protected.ThenFunc(app.deleteAdminAccount))
	mux.Handle("GET /admin/event/{id}/edit", protected.ThenFunc(app.editEventForm))
	mux.Handle("POST /admin/events/{id}/edit", protected.ThenFunc(app.editChurchEvent))

	mux.Handle("GET /image", dynamic.ThenFunc(app.serveDriveImageHandler))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
