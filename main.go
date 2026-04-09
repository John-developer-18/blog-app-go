package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("./templates/*.html"))

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
	switch r.Method {
	case "GET":
		err := tmpl.ExecuteTemplate(w, "login", nil)

		if err != nil {
			http.Error(w, "Internal Server Error : 500", 500)
			return
		}
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

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/posts", posts)
	http.HandleFunc("/make-posts", makePosts)
	http.HandleFunc("/post", post)

	fmt.Println("http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
