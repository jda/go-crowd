package crowd

type Error struct {
	XMLName struct{} `xml:"error"`
	Reason  string   `xml:"reason"`
	Message string   `xml:"message"`
}
