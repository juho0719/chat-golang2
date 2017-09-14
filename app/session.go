package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/goincremental/negroni-sessions"
)

const (
	currentUserKey = "oauth2_current_user"	// stored CurrentUser's key to session
	sessionDuration = time.Hour
)

type User struct {
	Uid string `json:"uid"`
	Name string `json:"name"`
	Email string `json:"user"`
	AvatarUrl string `json:"avatar_url"`
	Expired time.Time `json:"expired"`
}

func (u *User) Valid() bool {
	//confirm expired date based on current time.
	return u.Expired.Sub(time.Now()) > 0
}

func (u *User) Refresh() {
	//extend expired date
	u.Expired = time.Now().Add(sessionDuration)
}

func GetCurrentUser(r *http.Request) *User {
	//get a CurrentUser information from session
	s := sessions.GetSession(r)

	if s.Get(currentUserKey) == nil {
		return nil
	}

	data := s.Get(currentUserKey).([]byte)
	var u User
	json.Unmarshal(data, &u)
	return &u
}

func SetCurrentUser(r *http.Request, u *User)  {
	if u != nil {
		//renew expired date
		u.Refresh()
	}

	// stored CurrentUser information in the session as json
	s := sessions.GetSession(r)
	val, _ := json.Marshal(u)
	s.Set(currentUserKey, val)
}
