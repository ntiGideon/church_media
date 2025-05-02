package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/migrate"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func connectDb(connectionString string) (*ent.Client, error) {
	client, err := ent.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	err = client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true))
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to the database")
	return client, nil
}

func OpenDriver(connectionString string) (*sql.DB, error) {
	drv, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return drv, nil
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1>An Error occurred at our side!</h1>"))
	app.logger.Error("An error occurred at our side!", "status", status, "trace", string(debug.Stack()))
	http.Error(w, http.StatusText(status), status)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		return
	}
}

// validateImageUploads validates a specific image before upload
func (app *application) validateImageUploads(r *http.Request, image string) (string, error) {
	file, handler, err := r.FormFile(image)
	if err != nil {
		app.logger.Error("Error retrieving file", err)
		return "", err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".jpeg": true,
	}

	ext := filepath.Ext(handler.Filename)
	if !allowedExtensions[strings.ToLower(ext)] {
		app.logger.Error("Invalid file extension", handler.Filename)
		return "", err
	}

	uploadedFileURL, err := uploadImageToDrive(app.uploadService, file, handler.Filename)
	if err != nil {
		app.logger.Error("Error uploading file:", err)
		return "", err
	}

	return uploadedFileURL, nil
}

func (app *application) renderAdmin(w http.ResponseWriter, r *http.Request, status int, page string, data templateDataAdmin) {
	ts, ok := app.templateCacheAdmin[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		return
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination that we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// snippetCreatePost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err
	}

	return nil
}

// A constructor for the TemplateData
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		CSRFToken:       nosurf.Token(r),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func (app *application) newTemplateAdmin(r *http.Request) templateDataAdmin {
	return templateDataAdmin{
		UnreadMessagesCount: app.messageClient.GetUnreadMessagesCount(r.Context()),
		Flash:               app.sessionManager.PopString(r.Context(), "flash"),
		CSRFToken:           nosurf.Token(r),
		FormNumber:          randomString(11),
		IsAuthenticated:     app.isAuthenticated(r),
	}
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func getPaginationPages(currentPage, totalPages int) []int {
	var pages []int

	startPage := maxFunc(1, currentPage-2)
	endPage := minFunc(currentPage, totalPages+2)

	if currentPage <= 3 {
		endPage = minFunc(5, totalPages)
	} else if currentPage >= totalPages-2 {
		startPage = maxFunc(totalPages-4, 1)
	}

	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}
	return pages
}

func maxFunc(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minFunc(a, b int) int {
	if a < b {
		return a
	}
	return b
}
