// Package crowd provides methods for interacting with the
// Atlassian Crowd authentication, directory integration, and
// Single Sign-On system.
package crowd

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// Crowd represents your Crowd (client) Application settings
type Crowd struct {
	user    string
	passwd  string
	url     string
	cookies http.CookieJar
	Client *http.Client
}

// New initializes & returns a Crowd object.
func New(appuser string, apppass string, baseurl string) (Crowd, error) {
	if !strings.HasSuffix(baseurl, "/") {
		baseurl += "/"
	}

	cr := Crowd{
		Client: http.DefaultClient,
		user:   appuser,
		passwd: apppass,
		url:    baseurl,
	}

	cj, err := cookiejar.New(nil)
	if err != nil {
		return cr, err
	}

	cr.cookies = cj

	return cr, nil
}

func (c *Crowd) get() {

}
