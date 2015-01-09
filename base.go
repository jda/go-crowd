package crowd

import (
	"net/http"
	"net/http/cookiejar"
)

type Crowd struct {
	user    string
	passwd  string
	url     string
	cookies http.CookieJar
}

func New(appuser string, apppass string, baseurl string) (Crowd, error) {
	cr := Crowd{
		user:   appuser,
		passwd: apppass,
		url:    baseurl,
	}

	// TODO make sure URL ends with '/'

	cj, err := cookiejar.New(nil)
	if err != nil {
		return cr, err
	}

	cr.cookies = cj

	return cr, nil
}

func (c *Crowd) get() {

}
