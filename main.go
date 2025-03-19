package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/priyanka-choubey/stock-trade/handlers"
)

type App struct {
	Port         string
	StaticBase   string
	ErrorCode    int
	ErrorMessage string
}

func (a App) Start() {
	if a.StaticBase == "/static" {
		log.Printf("serving static assets")
		http.Handle("/static/", logreq(staticHandler("static")))
	}
	http.Handle("/", logreq(a.index))
	http.Handle("/error", logreq(a.error))
	http.Handle("/login", logreq(a.login))
	http.Handle("/signup", logreq(a.signup))
	http.Handle("/login_user", logreq(a.login_user))
	http.Handle("/signup_user", logreq(a.signup_user))
	addr := fmt.Sprintf(":%s", a.Port)
	log.Printf("Starting app on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func logreq(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("path: %s", r.URL.Path)

		f(w, r)
	})
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// This is inefficient - it reads the templates from the
	// filesystem every time. This makes it much easier to
	// develop though, so I can edit my templates and the
	// changes will be reflected without having to restart
	// the app.
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), 500)
		return
	}

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), 500)
		return
	}
}

func (a App) index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", struct {
		Name       string
		StaticBase string
	}{
		Name:       "Sonic The Hedgehog",
		StaticBase: a.StaticBase,
	})
}

func (a App) error(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "error.html", struct {
		Code    int
		Message string
	}{
		Code:    a.ErrorCode,
		Message: a.ErrorMessage,
	})
}

func (a App) signup(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "signup.html", struct {
		StaticBase string
	}{
		StaticBase: a.StaticBase,
	})
}

func (a App) signup_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var Username = r.Form["uname"]
	var Password = r.Form["token"]

	err := handlers.CreateUser(Username, Password)
	if err != nil {
		a.ErrorCode = err.code
		a.ErrorMessage = err.message
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (a App) login(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", struct {
		StaticBase string
	}{
		StaticBase: a.StaticBase,
	})
}

func (a App) login_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var Username = r.Form["uname"]
	var Password = r.Form["token"]

	err := handlers.AuthenticateUser(Username, Password)
	if err != nil {
		a.ErrorCode = err.code
		a.ErrorMessage = err.message
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func staticHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
	}
}

func main() {
	server := App{
		Port:       env("PORT", "8080"),
		StaticBase: env("STATIC_BASE", "/static"),
	}
	server.Start()
}
