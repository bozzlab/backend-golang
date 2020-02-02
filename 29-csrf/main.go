// package main

// import (
// 	"crypto/rand"
// 	"encoding/base64"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// )

// func main() {
// 	http.HandleFunc("/", index)
// 	http.HandleFunc("/signin", signIn)
// 	http.HandleFunc("/signout", signOut)
// 	http.HandleFunc("/transfer", transfer)
// 	http.ListenAndServe(":3333", nil)
// }

// var sessionStore = make(map[string]*session)

// type session struct {
// 	ID        string
// 	UserID    int
// 	CSRFToken string
// }

// func findSession(sessionID string) *session {
// 	return sessionStore[sessionID]
// }

// func storeSession(sessionID string, sess *session) {
// 	sessionStore[sessionID] = sess
// }

// func getSessionFromRequest(w http.ResponseWriter, r *http.Request) *session {
// 	var sess *session

// 	if c, err := r.Cookie("session"); err == nil {
// 		sess = findSession(c.Value)
// 	}

// 	if sess == nil {
// 		sess = &session{
// 			ID: generateSessionID(),
// 		}
// 		storeSession(sess.ID, sess)
// 		http.SetCookie(w, &http.Cookie{
// 			Name:     "session",
// 			Value:    sess.ID,
// 			Path:     "/",
// 			HttpOnly: true,
// 			SameSite: http.SameSiteLaxMode,
// 		})
// 	}
// 	return sess
// }

// func generateSessionID() string {
// 	b := make([]byte, 16)
// 	io.ReadFull(rand.Reader, b)
// 	return base64.RawURLEncoding.EncodeToString(b)
// }

// func generateCSRFToken() string {
// 	b := make([]byte, 32)
// 	io.ReadFull(rand.Reader, b)
// 	return base64.RawURLEncoding.EncodeToString(b)
// }

// func index(w http.ResponseWriter, r *http.Request) {
// 	var userID int

// 	sess := getSessionFromRequest(w, r)
// 	if sess != nil {
// 		userID = sess.UserID
// 	}

// 	if userID == 0 {
// 		w.Write([]byte(`
// 			<!doctype html>
// 			<a href=/signin>Sign In</a>
// 		`))
// 		return
// 	}

// 	// add csrf token to form
// 	fmt.Fprintf(w, `
// 		<!doctype html>
// 		<form method=POST action=/transfer>
// 			<input name=amount placeholder=amount>
// 			<input name=csrf value=%s type=hidden >
// 			<button type=submit>Transfer</button>
// 		</form>
// 		<a href=/signout>Sign Out</a>
// 	`, sess.CSRFToken)
// }

// func signIn(w http.ResponseWriter, r *http.Request) {
// 	// session fixation
// 	sess := getSessionFromRequest(w, r)
// 	sess.UserID = 1
// 	sess.CSRFToken = generateCSRFToken()

// 	// generate new csrf token

// 	http.Redirect(w, r, "/", http.StatusFound)
// }

// func signOut(w http.ResponseWriter, r *http.Request) {
// 	sess := getSessionFromRequest(w, r)
// 	sess.UserID = 0

// 	http.Redirect(w, r, "/", http.StatusFound)
// }

// func transfer(w http.ResponseWriter, r *http.Request) {
// 	// allow only POST
// 	if r.Method != "POST" {
// 		http.Error(w, "method not allow", 400)
// 		return
// 	}
// 	// check origin
// 	// if r.Header.Get("Origin") != "http://localhost:3333/" {
// 	// 	http.Error(w, "invalid origin", 400)
// 	// 	return
// 	// }
// 	// check referer
// 	// if !strings.HasPrefix(r.Referer(), "http://localhost:3333/") {
// 	// 	http.Error(w, "invalid referer", 400)
// 	// 	return
// 	// }
// 	// check is user sign in
// 	sess := getSessionFromRequest(w, r)
// 	// if sess.UserID == 0 {
// 	// 	http.Error(w, "Unauthorize", 400)
// 	// 	return
// 	// }
// 	// get amount from form
// 	amount := r.PostFormValue("amount")
// 	// get csrf token from form
// 	csrfToken := r.PostFormValue("csrf")
// 	log.Print(csrfToken)
// 	log.Print(sess.CSRFToken)

// 	// check form csrf token with session csrf token
// 	// if subtle.ConstantTimeCompare([]byte(csrfToken), []byte(sess.CSRFToken)) != 1 {
// 	// 	http.Error(w, "invalid csrf token", 400)
// 	// 	return
// 	// }
// 	if csrfToken != sess.CSRFToken {
// 		http.Error(w, "CSRF Detected!!!, Invalid CSRF Token", http.StatusBadRequest)
// 		return
// 	}
// 	// if transfer success print amount to console
// 	log.Print("transfer success: $", amount)
// 	// redirect to /
// 	http.Redirect(w, r, "/", http.StatusFound)
// }

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/signin", signIn)
	http.HandleFunc("/signout", signOut)
	http.HandleFunc("/transfer", transfer)
	http.ListenAndServe(":3333", nil)
}

var sessionStore = make(map[string]*session)

type session struct {
	ID        string
	UserID    int
	CSRFToken string
}

func findSession(sessionID string) *session {
	return sessionStore[sessionID]
}

func storeSession(sessionID string, sess *session) {
	sessionStore[sessionID] = sess
}

func getSessionFromRequest(w http.ResponseWriter, r *http.Request) *session {
	var sess *session

	if c, err := r.Cookie("session"); err == nil {
		sess = findSession(c.Value)
	}

	if sess == nil {
		sess = &session{
			ID: generateSessionID(),
		}
		storeSession(sess.ID, sess)
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sess.ID,
			Path:     "/",
			HttpOnly: true,
		})
	}
	return sess
}

func generateSessionID() string {
	b := make([]byte, 16)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func index(w http.ResponseWriter, r *http.Request) {
	var userID int

	sess := getSessionFromRequest(w, r)
	if sess != nil {
		userID = sess.UserID
	}

	// not sign in
	if userID == 0 {
		w.Write([]byte(`
			<!doctype html>
			<a href=/signin>Sign In</a>
		`))
		return
	}

	fmt.Fprintf(w, `
		<!doctype html>
		<form method=POST action=/transfer>
			<input name=amount placeholder=amount required>
			<input name=csrf_token value=%s type=hidden>
			<button type=submit>Transfer</button>
		</form>
		<a href=/signout>Sign Out</a>
	`, sess.CSRFToken)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	// session fixation
	sess := getSessionFromRequest(w, r)
	sess.UserID = 1
	sess.CSRFToken = generateCSRFToken()

	http.Redirect(w, r, "/", http.StatusFound)
}

func signOut(w http.ResponseWriter, r *http.Request) {
	sess := getSessionFromRequest(w, r)
	sess.UserID = 0

	http.Redirect(w, r, "/", http.StatusFound)
}

func transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// check origin
	// if s := r.Header.Get("Origin"); s != "" && s != "http://localhost:3333" {
	// 	http.Error(w, "CSRF Detected!!!, Invalid Origin", http.StatusBadRequest)
	// 	return
	// }

	// check referer
	// if s := r.Header.Get("Referer"); s != "" && !strings.HasPrefix(s, "http://localhost:3333/") {
	// 	http.Error(w, "CSRF Detected!!!, Invalid Referer", http.StatusBadRequest)
	// 	return
	// }

	sess := getSessionFromRequest(w, r)
	// if sess.UserID == 0 {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	amount := r.FormValue("amount")
	csrfToken := r.FormValue("csrf_token")

	if csrfToken != sess.CSRFToken {
		http.Error(w, "CSRF Detected!!!, Invalid CSRF Token", http.StatusBadRequest)
		return
	}

	log.Println("transfer", amount)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
