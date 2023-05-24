package main

import (
	auth "forum/authentication"
	"forum/controller"
	"forum/dbmanagement"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/posts" {
		controller.PageErrors(w, r, tmpl, 400, "Bad requests")
		return
	}
	controller.AllPosts(w, r, tmpl)

}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	tags := dbmanagement.SelectAllTags()
	tagexists := false
	var url string
	for _, v := range tags {
		if r.URL.Path == "/categories/"+v.TagName {
			url = v.TagName
			tagexists = true
		}
	}
	if r.URL.Path == "/categories/" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	if !tagexists && r.URL.Path != "/categories/" {
		controller.PageErrors(w, r, tmpl, 400, "Bad requests")
		return
	}

	if tagexists && r.URL.Path != "/" {
		controller.SubForum(w, r, tmpl, url)
	}
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := dbmanagement.SelectAllPosts()
	if err != nil {
		controller.PageErrors(w, r, tmpl, 500, "Internal Server Error")
		return
	}
	postexists := false
	var url string
	for _, v := range posts {
		if r.URL.Path == "/posts/"+v.UUID {
			url = v.UUID
			postexists = true
		}
	}
	if r.URL.Path == "/posts/" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	if !postexists && r.URL.Path != "/posts/" {
		controller.PageErrors(w, r, tmpl, 400, "Bad requests")
		return
	}
	if postexists && r.URL.Path != "/" {
		controller.Post(w, r, tmpl, url)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	auth.Login(w, r, tmpl)
}

func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	auth.Authenticate(w, r, tmpl)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	auth.Logout(w, r, tmpl)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	auth.Register(w, r, tmpl)
}

func RegisterAccountHandler(w http.ResponseWriter, r *http.Request) {
	auth.RegisterAcount(w, r, tmpl)
}

// oauth handlers
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	auth.GoogleLogin(w, r, tmpl)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	auth.GoogleCallback(w, r, tmpl)
}

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	auth.GithubLogin(w, r, tmpl)
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	auth.GithubCallback(w, r, tmpl)
}

func FacebookLoginHandler(w http.ResponseWriter, r *http.Request) {
	auth.FacebookLogin(w, r, tmpl)
}

func FacebookCallbackHandler(w http.ResponseWriter, r *http.Request) {
	auth.FacebookCallback(w, r, tmpl)
}

// forum handlers
func ForumHandler(w http.ResponseWriter, r *http.Request) {
	controller.AllPosts(w, r, tmpl)
}

func SubmitPostHandler(w http.ResponseWriter, r *http.Request) {
	controller.SubmitPost(w, r, tmpl)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	controller.Admin(w, r, tmpl)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	controller.User(w, r, tmpl)
}

func PrivacyPolicyHandler(w http.ResponseWriter, r *http.Request) {
	controller.PrivacyPolicy(w, r, tmpl)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	controller.PageErrors(w, r, tmpl, 429, "Too much requests")
}

func OauthErrorHandler(w http.ResponseWriter, r *http.Request) {
	controller.PageErrors(w, r, tmpl, 500, "Something went wrong...")
}
