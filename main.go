package main

import (
	"log"
	"net"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/namsral/flag"
	"github.com/startdusk/tiny-pastebin/controller"
	"github.com/startdusk/tiny-pastebin/view"
)

var (
	databaseURL = flag.String("database-url", "postgres://postgres:postgres@localhost:5432/paste?sslmode=disable", "Database URL")
	host        = flag.String("host", "0.0.0.0", "host for listen")
	port        = flag.String("port", "8088", "port for listen")
)

func main() {
	flag.Parse()
	conn, err := sqlx.Open("postgres", *databaseURL)
	if err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(32)
	addr := net.JoinHostPort(*host, *port)
	log.Println("Tiny code listen on:", addr)
	if err := run(conn, addr); err != nil {
		panic(err)
	}
}

func run(conn *sqlx.DB, addr string) error {
	mux := http.NewServeMux()
	view := view.CreatePasteView("./view/static")
	pasteHandler, err := controller.CreatePasteHandler(conn, view)
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(
		http.Dir("tiny-pastebin/view/static")))

	mux.HandleFunc("GET /", pasteHandler.Index)
	mux.HandleFunc("POST /", pasteHandler.CreatePaste)
	mux.HandleFunc("GET /{code}", pasteHandler.GetPaste)
	err = http.ListenAndServe(addr, mux)
	return err
}
