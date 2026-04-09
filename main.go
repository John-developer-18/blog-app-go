package main

import (
	models "blog-app/models"
	"blog-app/storage"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var tmpl = template.Must(template.ParseGlob("./templates/*.html"))

//Middleware for authentication and redirection

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")

		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}
		session_id := session.Value

		_, ok := storage.Sessions[session_id]

		if !ok {
			http.Redirect(w, r, "/login", 303)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func guestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := r.Cookie("session_id")

		if err == nil {
			session_id := session.Value

			_, ok := storage.Sessions[session_id]

			if ok {
				http.Redirect(w, r, "/make-posts", http.StatusSeeOther)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/home" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := tmpl.ExecuteTemplate(w, "home", nil)

	if err != nil {
		http.Error(w, "Internal Server Error : 500", 500)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	session := uuid.New().String()

	switch r.Method {
	case "GET":
		err := tmpl.ExecuteTemplate(w, "login", nil)

		if err != nil {
			http.Error(w, "Internal Server Error : 500", 500)
			return
		}
	case "POST":

		errors := r.ParseForm()

		if errors != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		value, ok := storage.Users[username]

		if !ok || value.Password != password {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		storage.Sessions[session] = username
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/make-posts", 303)

	}
}

func register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := tmpl.ExecuteTemplate(w, "register", nil)

		if err != nil {
			http.Error(w, "Internal Server Error : 500", 500)
			return
		}
	case "POST":
		err := r.ParseForm()

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		newUser := models.Users{
			ID:       uuid.New().ID(),
			Username: username,
			Password: password,
		}

		_, ok := storage.Users[username]

		//I might not need this line because I might handle it later in middleware
		if ok {
			http.Redirect(w, r, "/make-posts", http.StatusSeeOther)
			return
		}
		//take note of this line

		storage.Users[username] = newUser

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
}

func posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := tmpl.ExecuteTemplate(w, "posts", nil)

		if err != nil {
			http.Error(w, "Internal Server Error : 500", 500)
			return
		}
	}
}

func makePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := tmpl.ExecuteTemplate(w, "make-posts", nil)

		if err != nil {
			http.Error(w, "Internal Server Error : 500", 500)
			return
		}
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := tmpl.ExecuteTemplate(w, "post", nil)

		if err != nil {
			http.Error(w, "Internal Server Error : 500", 500)
			return
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if err == nil {
		sessionID := session.Value
		delete(storage.Sessions, sessionID)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", 303)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/make-posts", authMiddleware(http.HandlerFunc(makePosts)))
	http.Handle("/login", guestMiddleware(http.HandlerFunc(login)))
	http.Handle("/register", guestMiddleware(http.HandlerFunc(register)))
	http.HandleFunc("/", home)
	http.HandleFunc("/posts", posts)
	http.HandleFunc("/post", post)
	http.HandleFunc("/logout", logout)

	fmt.Println("http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
