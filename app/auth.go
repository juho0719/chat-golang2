package main

import (
	"log"
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/codegangsta/negroni"
	"strings"
)

const (
	nextPageKey = "next_page"
	authSecurityKey = "auth_security_key"
)

func init()  {
	// set gomniauth information
	gomniauth.SetSecurityKey(authSecurityKey)
	gomniauth.WithProviders(
		google.New("292640455104-rdjbp61r50tinuer3f6lb66k878q0s76.apps.googleusercontent.com", "Bu89mYdhTHZwwB9USUlOPkrs", "http://127.0.0.1:3000/auth/callback/google"),
	)
}

func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	action := ps.ByName("action")
	provider := ps.ByName("provider")
	s := sessions.GetSession(r)

	switch action {
	case "login":
		//move the login page of gomniauth.Provider
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		loginUrl, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		http.Redirect(w, r, loginUrl, http.StatusFound)
	case "callback":
		//gomniauth callback process
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln(err)
		}

		//confirm user information from callback result
		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln(err)
		}

		u := &User{
			Uid: user.Data().Get("id").MustStr(),
			Name: user.Name(),
			Email: user.Email(),
			AvatarUrl: user.AvatarURL(),
		}

		SetCurrentUser(r, u)
		http.Redirect(w, r, s.Get(nextPageKey).(string), http.StatusFound)
	default:
		http.Error(w, "Auth action '"+action+"' is not supported", http.StatusNotFound)
	}
}

func LoginRequired(ignore ...string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// execute next handler if ignore url
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next(w, r)
				return
			}
		}

		//get CurrentUser information
		u := GetCurrentUser(r)

		// execute next handler after renew expired date if CurrentUser information valid
		if u != nil && u.Valid() {
			SetCurrentUser(r, u)
			next(w, r)
			return
		}

		// set nil for CurrentUser if CurrentUser information invalid
		SetCurrentUser(r, nil)

		// store moving url to session after login
		sessions.GetSession(r).Set(nextPageKey, r.URL.RequestURI())

		//redirect login page
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
