package crowd

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
