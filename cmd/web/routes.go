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
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	mux.Handle("GET /", dynamic.ThenFunc(app.home))
	mux.Handle("GET /ministries", dynamic.ThenFunc(app.ministries))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /events", dynamic.ThenFunc(app.events))
	mux.Handle("GET /sermons", dynamic.ThenFunc(app.sermon))
	mux.Handle("GET /contact", dynamic.ThenFunc(app.contact))
	mux.Handle("POST /contact", dynamic.ThenFunc(app.contactForm))
	mux.Handle("GET /give-online", dynamic.ThenFunc(app.giveOnline))
	mux.Handle("GET /dashboard", dynamic.ThenFunc(app.dashboard))
	mux.Handle("GET /dashboards", dynamic.ThenFunc(app.dashboardData))
	mux.Handle("POST /member", dynamic.ThenFunc(app.memberForm))
	mux.Handle("DELETE /members/{id}/delete", dynamic.ThenFunc(app.deleteMember))
	mux.Handle("POST /member/{id}", dynamic.ThenFunc(app.memberForm))
	mux.Handle("GET /member/success/{id}", dynamic.ThenFunc(app.memberSuccess))
	mux.Handle("GET /member", dynamic.ThenFunc(app.members))
	mux.Handle("GET /member/{id}", dynamic.ThenFunc(app.membersEdit))
	mux.Handle("GET /list-members", dynamic.ThenFunc(app.listMembers))
	mux.Handle("GET /members/{id}/view", dynamic.ThenFunc(app.memberDetailPage))
	mux.Handle("GET /mens", dynamic.ThenFunc(app.mensMinistry))
	mux.Handle("GET /womens", dynamic.ThenFunc(app.womenMinistry))
	mux.Handle("GET /youth", dynamic.ThenFunc(app.youthMinistry))
	mux.Handle("GET /childrens", dynamic.ThenFunc(app.childrenMinistry))
	mux.Handle("GET /teens", dynamic.ThenFunc(app.teensMinistry))
	mux.Handle("GET /music", dynamic.ThenFunc(app.musicMinistry))
	mux.Handle("GET /ushers", dynamic.ThenFunc(app.ushersMinistry))
	mux.Handle("GET /service-record", dynamic.ThenFunc(app.serviceRecordForm))
	mux.Handle("POST /service-record", dynamic.ThenFunc(app.CreateServiceRecord))
	mux.Handle("GET /list-records", dynamic.ThenFunc(app.listServiceRecords))
	mux.Handle("DELETE /service-record/{id}/delete", dynamic.ThenFunc(app.DeleteRecord))
	mux.Handle("GET /church-events", dynamic.ThenFunc(app.churchEvents))
	mux.Handle("POST /event-create", dynamic.ThenFunc(app.churchEventForm))

	mux.Handle("GET /image", dynamic.ThenFunc(app.serveDriveImageHandler))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
