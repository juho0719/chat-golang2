package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/gorilla/websocket"

	"gopkg.in/mgo.v2"
	"log"
)

const (
	sessionKey = "simple_chat_session"
	sessionSecret = "simple_chat_session_secret"

	socketBufferSize = 1024
)

var (
	renderer *render.Render
	mongoSession *mgo.Session

	upgrader = &websocket.Upgrader{
		ReadBufferSize: socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

func init()  {
	//create renderer
	renderer = render.New()

	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	mongoSession = s
}

func main()  {
	//create router
	router := httprouter.New()

	//definition handler
	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		//template rendering
		renderer.HTML(w, http.StatusOK, "index", map[string]interface{}{"host": req.Host})
	})

	router.GET("/info", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		u := GetCurrentUser(req)
		info := map[string]interface{}{"current_user": u, "clients": clients}
		renderer.JSON(w, http.StatusOK, info)
	})

	router.GET("/login", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// render login page
		renderer.HTML(w, http.StatusOK, "login", nil)
	})
	router.GET("/logout", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		//move the login page after remove user information in the session
		sessions.GetSession(req).Delete(currentUserKey)
		http.Redirect(w, req, "/login", http.StatusFound)
	})
	router.GET("/auth/:action/:provider", loginHandler)
	router.POST("/rooms", createRoom)
	router.GET("/rooms", retrieveRooms)
	router.GET("/rooms/:id/messages", retrieveMessages)
	router.GET("/ws/:room_id", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		socket, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Fatal("ServeHTTP:", err)
			return
		}
		newClient(socket, ps.ByName("room_id"), GetCurrentUser(req))
	})

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