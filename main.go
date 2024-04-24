package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./models"
)

var posts map[string]*models.Post
var clients map[string]*models.Client

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Println(posts)

	t.ExecuteTemplate(w, "index", posts)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	id := r.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(w, r)
		return
	}

	t.ExecuteTemplate(w, "write", post)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.ExecuteTemplate(w, "write", nil)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	if id != "" {
		// Обновление существующего поста
		post, found := posts[id]
		if !found {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		post.Title = title
		post.Content = content
	} else {
		// Создание нового поста
		id = GenerateID()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	http.Redirect(w, r, "/feed", http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	delete(posts, id)
	http.Redirect(w, r, "/feed", http.StatusFound)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Логика для получения информации о пользователе и его постах

	// Преобразуем карту постов в срез постов
	var userPosts []*models.Post
	for _, post := range posts {
		userPosts = append(userPosts, post)
	}

	t, err := template.ParseFiles("templates/profile.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	// Пример данных о постах пользователя для передачи в шаблон
	userData := struct {
		Username string
		Email    string
		Posts    []*models.Post
	}{
		Username: "John Doe",
		Email:    "johndoe@example.com",
		Posts:    userPosts, // Используем срез userPosts
	}

	t.ExecuteTemplate(w, "profile", userData)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/headLog.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	if r.Method == "POST" {
		http.Redirect(w, r, "/feed", http.StatusFound)
		return
	}

	t.ExecuteTemplate(w, "login", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		patronymic := r.FormValue("patronymic")
		birthdate := r.FormValue("birthdate")

		client := models.NewClient("", email, firstName, lastName, patronymic, birthdate, "")
		clients[client.Email] = client

		http.Redirect(w, r, "/feed", http.StatusFound)
		return
	}

	t, err := template.ParseFiles("templates/registration.html", "templates/headLog.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "registration", nil)
}

func main() {

	posts = make(map[string]*models.Post)
	clients = make(map[string]*models.Client)

	http.HandleFunc("/feed", indexHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/SavePost", savePostHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/", loginHandler)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	http.ListenAndServe(":3000", nil)

}
