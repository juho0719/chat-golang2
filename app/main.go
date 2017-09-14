package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

const (
	sessionKey = "simple_chat_session"
	sessionSecret = "simple_chat_session_secret"
)

var renderer *render.Render

func init()  {
	//create renderer
	renderer = render.New()
}

func main()  {
	//create router
	router := httprouter.New()

	//definition handler
	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		//template rendering
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Simple Chat!"})
	})

	router.GET("/login", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// render login page
		renderer.HTML(w, http.StatusOK, "login", nil)
	})
	router.GET("/logout", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		//move the login page after remove user information in the session
		sessions.GetSession(req).Delete(keyCurrentUser)
		http.Redirect(w, req, "/login", http.StatusFound)
	})

	router.GET("/auth/:action/:provider", loginHandler)
	//create negroni middleware
	n := negroni.Classic()
	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))

	n.Use(LoginRequired("/login", "/auth"))

	//regist router to handler by negroni
	n.UseHandler(router)

	//run webserver
	n.Run(":3000")
}