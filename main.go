package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// 1. получаем все статьи с сервера и выводим их на главную страницу
// 	1.1 создаем структуру, которая описывает таблицу статей на сервере
type Article struct {
	ID uint16
	Title, Anons, FullText string
}

// создаем список
var posts = []Article{}
var showPost = Article{}


// функция, срабатывает каждый раз когда переходим на главну страницу
func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// подключение к БД
	db, err := sql.Open("mysql", "u1036535_default:sJbUbb6_@tcp(https://riportfolio.ru/:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// выборка данных из БД
	result, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for result.Next() {
		var post Article
		err = result.Scan(&post.ID, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	tmpl.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	if title == "" || anons == "" || fullText == "" {
		fmt.Fprintf(w, "Заполните все поля формы")
	} else {
		// подключение к базе данных
		db, err := sql.Open("mysql", "u1036535_default:sJbUbb6_@tcp(https://riportfolio.ru/:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		// установка данных
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, fullText))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		// переадресация пользователя после добавления статьи
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func showArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tmpl, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "u1036535_default:sJbUbb6_@tcp(https://riportfolio.ru/:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// выборка данных из БД
	result, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for result.Next() {
		var post Article
		err = result.Scan(&post.ID, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		showPost = post
	}

	tmpl.ExecuteTemplate(w, "show", showPost)
}

func handleFunc() {
	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/create", create).Methods("GET")
	router.HandleFunc("/save_article", saveArticle).Methods("POST")
	router.HandleFunc("/post/{id:[0-9]+}", showArticle).Methods("GET")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}