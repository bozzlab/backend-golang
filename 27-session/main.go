package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/signin", signIn)
	http.HandleFunc("/signout", signOut)
	http.ListenAndServe(":3333", nil)
}

var sessionStore = make(map[string]*session)

type session struct {
	UserID int

	//
}

func findSession(sessionID string) *session {
	return sessionStore[sessionID]
}

func storeSession(sessionID string, sess *session) {
	//
	sessionStore[sessionID] = sess
}

func getSessionFromRequest(w http.ResponseWriter, r *http.Request) *session {
	// read session id from cookie
	if cookie, _ := r.Cookie("session"); cookie != nil {
		// find session from database
		sess := findSession(cookie.Value)
		if sess != nil {
			return sess
		}
	}
	// if session not found create new session, and set session id to cookie
	sess := &session{}
	sessID := generateSessionID()

	storeSession(sessID, sess)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessID,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})
	return sess
}

func generateSessionID() string {
	// create new slice
	b := make([]byte, 32)
	// read data from rand.Reader from package crypto/rand
	rand.Read(b)
	// encode to base64
	return base64.StdEncoding.EncodeToString(b)
}

func index(w http.ResponseWriter, r *http.Request) {
	// get session from request
	sess := getSessionFromRequest(w, r)
	// if not sign in, return sign in page

	if sess.UserID == 0 {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
		<!doctype html>
		<h1>hello</h1>
		<a href=/signin>Sign in </a>`))
		return
	}
	// if already sign in return sign out page
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
	<!doctype html>
	<h1>xxx</h1>
	<a href=/signout>Sign out </a>`))
	return

}

func signIn(w http.ResponseWriter, r *http.Request) {
	sess := &session{}
	sessID := generateSessionID()
	// TODO: session fixation hack
	// sess := getSessionFromRequest(w, r)
	// get session from request
	// set user id to `1`
	sess.UserID = 1
	storeSession(sessID, sess)
	storeSession(sessID, sess)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessID,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})
	// redirect to /
	http.Redirect(w, r, "/", http.StatusFound)

}

func signOut(w http.ResponseWriter, r *http.Request) {
	// get session from request
	sess := getSessionFromRequest(w, r)

	// set user id to `0`
	sess.UserID = 0

	// redirect to /
	http.Redirect(w, r, "/", http.StatusFound)

}
