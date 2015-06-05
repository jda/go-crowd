// Package crowdauth provides middleware for Crowd SSO logins
//
// Goals:
//  1) drop-in authentication against Crowd SSO
//  2) make it easy to use Crowd SSO as part of your own auth
package crowdauth // import "go.jona.me/crowd/crowdauth"

import (
	"go.jona.me/crowd"
	"html/template"
	"log"
	"net"
	"net/http"
)

type SSO struct {
	CrowdApp            *crowd.Crowd
	LoginPage           AuthLoginPage
	LoginTemplate       *template.Template
	ClientAddressFinder ClientAddressFinder
}

// The AuthLoginPage type extends the normal http.HandlerFunc type
// with a boolean return to indicate login success or failure.
type AuthLoginPage func(http.ResponseWriter, *http.Request, *SSO) bool

type ClientAddressFinder func(*http.Request) (string, error)

var authErr string = "unauthorized, login required"

func DefaultClientAddressFinder(r *http.Request) (addr string, err error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return host, nil
}

// New returns creates a new instance of SSO
func New(user string, password string, crowdURL string) (s SSO, err error) {
	s.LoginPage = loginPage
	s.ClientAddressFinder = DefaultClientAddressFinder
	s.LoginTemplate = template.Must(template.New("authPage").Parse(defLoginPage))

	cwd, err := crowd.New(user, password, crowdURL)
	if err != nil {
		return s, err
	}
	s.CrowdApp = &cwd

	return s, nil
}

// Handler provides HTTP middleware using http.Handler chaining
// that requires user authentication via Atlassian Crowd SSO.
func (s *SSO) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.loginHandler(w, r) == false {
			return
		}
		h.ServeHTTP(w, r)
	})

}

func (s *SSO) loginHandler(w http.ResponseWriter, r *http.Request) bool {
	cc, err := s.CrowdApp.GetCookieConfig()
	if err != nil {
		http.Error(w, "service error", http.StatusInternalServerError)
		return false
	}

	ck, err := r.Cookie(cc.Name)
	if err == http.ErrNoCookie {
		// no cookie so show login page if GET
		// if POST check if login and handle
		// if fail, show login page with message
		if r.Method == "GET" {
			s.LoginPage(w, r, s)
			return false
		} else if r.Method == "POST" {
			authOK := s.LoginPage(w, r, s)
			if authOK != true {
				log.Printf("crowdauth: authentication failed\n")
				return false
			}
		} else {
			http.Error(w, authErr, http.StatusUnauthorized)
			return false
		}
	} else {
		// validate cookie or show login page
		host, err := s.ClientAddressFinder(r)
		if err != nil {
			log.Printf("crowdauth: could not get remote addr: %s\n", err)
		}

		_, err = s.CrowdApp.ValidateSession(ck.Value, host)
		if err != nil {
			log.Printf("crowdauth: could not validate cookie, deleting because: %s\n", err)
			ck.MaxAge = -1
			http.SetCookie(w, ck)
			s.LoginPage(w, r, s)
			return false
		}

		// valid cookie so fallthrough
	}
	return true
}

func (s *SSO) Login(user string, pass string, addr string) (cs crowd.Session, err error) {
	cs, err = s.CrowdApp.NewSession(user, pass, addr)
	return cs, err
}

func (s *SSO) StartSession(w http.ResponseWriter, cs crowd.Session) (err error) {
	cc, err := s.CrowdApp.GetCookieConfig()
	if err != nil {
		return err
	}

	ck := http.Cookie{
		Name:    cc.Name,
		Domain:  cc.Domain,
		Secure:  false,
		Value:   cs.Token,
		Expires: cs.Expires,
	}
	http.SetCookie(w, &ck)
	return nil
}

func loginPage(w http.ResponseWriter, r *http.Request, s *SSO) bool {
	if r.Method == "GET" { // show login page and bail
		showLoginPage(w, s)
		return false
	} else if r.Method == "POST" {
		user := r.FormValue("username")
		pass := r.FormValue("password")
		host, err := s.ClientAddressFinder(r)
		if err != nil {
			log.Printf("crowdauth: could not get remote addr: %s\n", err)
			showLoginPage(w, s)
			return false
		}

		sess, err := s.Login(user, pass, host)
		if err != nil {
			log.Printf("crowdauth: login/new session failed: %s\n", err)
			showLoginPage(w, s)
			return false
		}

		err = s.StartSession(w, sess)
		if err != nil {
			log.Printf("crowdauth: could not save session: %s\n", err)
			showLoginPage(w, s)
			return false
		}
	} else {
		return false
	}

	return true
}

func showLoginPage(w http.ResponseWriter, s *SSO) {
	err := s.LoginTemplate.ExecuteTemplate(w, "authPage", nil)
	if err != nil {
		log.Printf("crowdauth: could not exec template: %s\n", err)
	}
}
