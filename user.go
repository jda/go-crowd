package crowd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// User represents a user in Crowd
type User struct {
	XMLName     struct{} `xml:"user"`
	UserName    string   `xml:"name,attr"`
	FirstName   string   `xml:"first-name"`
	LastName    string   `xml:"last-name"`
	DisplayName string   `xml:"display-name"`
	Email       string   `xml:"email"`
	Active      bool     `xml:"active"`
	Key         string   `xml:"key"`
}

// GetUser retrieves user information
func (c *Crowd) GetUser(user string) (User, error) {
	u := User{}

	v := url.Values{}
	v.Set("username", user)
	url := c.url + "rest/usermanagement/1/user?" + v.Encode()
	c.Client.Jar = c.cookies
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return u, err
	}
	req.SetBasicAuth(c.user, c.passwd)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml")
	resp, err := c.Client.Do(req)
	if err != nil {
		return u, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return u, fmt.Errorf("user not found")
	case 200:
		// fall through switch without returning
	default:
		return u, fmt.Errorf("request failed: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return u, err
	}

	err = xml.Unmarshal(body, &u)
	if err != nil {
		return u, err
	}

	return u, nil
}

// GetUserByEmail searches for a user by email
func (c *Crowd) GetUserByEmail(email string) (User, error) {
	u := User{}

	names, err := c.SearchUsers("email=" + email)
	if err != nil {
		return u, err
	}

	switch len(names) {
	case 1:
		return c.GetUser(names[0])
	case 0:
		return u, fmt.Errorf("user not found")
	default:
		return u, fmt.Errorf("multiple users found")
	}
}

// SearchUsers retrieves a slice of user names that match the given search restriction (as Crowd Query Language)
func (c *Crowd) SearchUsers(restriction string) ([]string, error) {
	var names []string

	v := url.Values{}
	v.Set("entity-type", "user")
	v.Set("restriction", restriction)
	url := c.url + "rest/usermanagement/1/search?" + v.Encode()
	c.Client.Jar = c.cookies
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return names, nil
	}
	req.SetBasicAuth(c.user, c.passwd)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml")
	resp, err := c.Client.Do(req)
	if err != nil {
		return names, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		// fall through switch without returning
	default:
		return names, fmt.Errorf("request failed: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return names, err
	}

	r := struct{
		Users      []struct{
			Name      string `xml:"name,attr"`
		} `xml:"user"`
	}{}

	err = xml.Unmarshal(body, &r)
	if err != nil {
		return names, err
	}

	for _, u := range r.Users {
		names = append(names, u.Name)
	}

	return names, nil
}