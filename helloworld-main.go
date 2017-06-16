package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	_ "github.com/lib/pq"
)

type server struct {
	DB *sql.DB
}

// Increment and display a counter for the IP address making the request
func (server *server) incrementer(w http.ResponseWriter, r *http.Request) {
	ip := guessIPofRequester(r)

	// Clearly the following is not transaction safe...
	var curCount = 0
	err := server.DB.QueryRow("select value from counter where name = $1;", ip).Scan(&curCount)
	switch err {
	case nil:
		curCount++
		_, err = server.DB.Exec("update counter set value = $1 where name = $2;", curCount, ip)
		if err != nil {
			log.Panic(err)
		}
	case sql.ErrNoRows:
		_, err = server.DB.Exec("insert into counter(name, value) values($1, $2);", ip, 1)
		if err != nil {
			log.Panic(err)
		}
	default:
		log.Panic(err)
	}

	fmt.Fprintf(w, "Hello %s. You have visited %d times.", ip, curCount)
}

// Initialize the database (TODO remove to separate tool)
func (server *server) bootstrap(w http.ResponseWriter, r *http.Request) {
	_, err := server.DB.Exec("CREATE TABLE IF NOT EXISTS counter (name varchar(255), value integer);")
	if err != nil {
		fmt.Fprintf(w, "Error creating db: %s", err.Error())
		return
	}

	fmt.Fprintln(w, "Success!")
}

// Return a database object, using the CloudFoundry environment data
func postgresDBFromCF() (*sql.DB, error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		return nil, err
	}

	dbEnv, err := appEnv.Services.WithTag("postgres")
	if err != nil {
		return nil, err
	}

	if len(dbEnv) != 1 {
		return nil, errors.New("expecting 1 database")
	}

	dbURI, ok := dbEnv[0].CredentialString("uri")
	if !ok {
		return nil, errors.New("no uri in creds for db")
	}
	// Service broker adds a reconnect=true that RDS doesn't seem to understand.
	// TODO - fix service broker?
	idx := strings.Index(dbURI, "?")
	if idx >= 0 {
		dbURI = dbURI[:idx]
	}

	return sql.Open("postgres", dbURI)
}

// Get an approximation of the IP address of the requestor
func guessIPofRequester(r *http.Request) string {
	forwardedIPs, ok := r.Header["X-Forwarded-For"]
	if !ok {
		forwardedIPs = nil
	}

	// Since some proxies add comma's, and others headers, handle both
	forwardedIPs = strings.Split(strings.Join(forwardedIPs, ","), ",")

	// Last one added should be added by our reverse proxy
	for idx := len(forwardedIPs) - 1; idx >= 0; idx-- {
		ip := strings.TrimSpace(forwardedIPs[idx])
		if len(ip) == 0 {
			continue
		}
		if strings.HasPrefix(ip, "10.") { // but apparently we must be behind at least 2, so skip any 10. addresses
			// TODO - be more precise about this
			continue
		}
		return ip
	}
	return "unknown"
}

func (server *server) CreateHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/bootstrap", server.bootstrap)
	mux.HandleFunc("/favicon.ico", http.NotFound) // if we don't handle this, then the "/" handler matches, and we get double-counts
	mux.HandleFunc("/", server.incrementer)
	return mux
}

func main() {
	// Get the database
	db, err := postgresDBFromCF()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start the app
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), (&server{
		DB: db,
	}).CreateHTTPHandler()))
}
