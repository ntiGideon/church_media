package main

import (
	"context"
	"crypto/tls"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/joho/godotenv"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/internal/models"
	"google.golang.org/api/drive/v3"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger              *slog.Logger
	templateCache       map[string]*template.Template
	templateCacheAdmin  map[string]*template.Template
	db                  *ent.Client
	formDecoder         *form.Decoder
	sessionManager      *scs.SessionManager
	messageClient       *models.MessageModel
	memberClient        *models.MemberModel
	recordServiceClient *models.ServiceModel
	eventClient         *models.EventModel
	userClient          *models.UserModel
	uploadService       *drive.Service
}

func main() {
	gob.Register(map[string]interface{}{})

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := flag.String("addr", ":9000", "HTTP service address")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	// postgres database connection
	connectionString := fmt.Sprintf("host=%v port=%v user=%v dbname=%s password=%s sslmode=%v", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASS"), os.Getenv("DB_SSL_MODE"))
	dbString := flag.String("dsn", connectionString, "database connection string")
	flag.Parse()
	db, err := connectDb(*dbString)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	dbDriver, err := OpenDriver(connectionString)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	templateCacheAdmin, err := newAdminTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//form decode
	formDecode := form.NewDecoder()

	// session manager
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour // persist cookie for one day
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.Persist = false // we are setting to false since we want to be able to control it upon login when user decides
	sessionManager.Store = postgresstore.New(dbDriver)

	// upload service
	driveService, err := createDriveService(context.Background(), "credentials.json")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		logger:              logger,
		templateCache:       templateCache,
		templateCacheAdmin:  templateCacheAdmin,
		db:                  db,
		formDecoder:         formDecode,
		sessionManager:      sessionManager,
		messageClient:       &models.MessageModel{Db: db},
		memberClient:        &models.MemberModel{Db: db},
		recordServiceClient: &models.ServiceModel{Db: db},
		eventClient:         &models.EventModel{Db: db},
		userClient:          &models.UserModel{DB: db},
		uploadService:       driveService,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server on port", "addr", server.Addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
