package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
)

var (
	// cfg is the global configuration for the server. It's read in at startup from
	// the config.json file and enviornment variables, see config.go for more info.
	cfg *config

	// When was the last alert sent out?
	// Use this value to avoid bombing alerts
	lastAlertSent *time.Time

	// log output
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

	// application database connection
	appDB *sql.DB

	// cookie session storage
	sessionStore *sessions.CookieStore
)

func main() {
	var err error
	cfg, err = initConfig(os.Getenv("GOLANG_ENV"))
	if err != nil {
		// panic if the server is missing a vital configuration detail
		panic(fmt.Errorf("server configuration error: %s", err.Error()))
	}
	if err = initKeys(cfg); err != nil {
		panic(fmt.Errorf("server keys error: %s", err.Error()))
	}

	sessionStore = sessions.NewCookieStore([]byte(cfg.SessionSecret))

	connectToAppDb()

	s := &http.Server{}
	m := http.NewServeMux()

	m.HandleFunc("/.well-known/acme-challenge/", CertbotHandler)
	m.Handle("/", middleware(HealthCheckHandler))
	m.Handle("/session", middleware(SessionHandler))
	m.Handle("/session/keys", middleware(KeysHandler))
	m.Handle("/jwt/publickey", middleware(JwtPublicKeyHandler))
	m.Handle("/jwt/session", middleware(JwtHandler))

	// m.Handle("/session/groups", handler)
	m.Handle("/search", middleware(UsersSearchHandler))
	m.Handle("/users", middleware(UsersHandler))
	m.Handle("/users/", middleware(UserHandler))

	m.Handle("/groups", middleware(GroupsHandler))
	m.Handle("/groups/", middleware(GroupHandler))

	// m.Handle("/reset", middleware(ResetPasswordHandler))
	// m.Handle("/reset/", middleware(ResetPasswordHandler))

	// connect mux to server
	s.Handler = m

	// print notable config settings
	// printConfigInfo()

	// fire it up!
	fmt.Println("starting server on port", cfg.Port)

	// start server wrapped in a log.Fatal b/c http.ListenAndServe will not
	// return unless there's an error
	logger.Fatal(StartServer(cfg, s))
}
