package oauth2

type linkedInInfo struct {
	Id        string            `json:"Id"`
	Elements  []linkedInHandle1 `json:"Elements"`
	FirstName string            `json:"LocalizedFirstName"`
	LastName  string            `json:"LocalizedLastName"`
}
type linkedInElements struct {
	Elements []linkedInHandle1 `json:"Elements"`
}
type linkedInEmail struct {
	EmailAddress string
}
type linkedInHandle1 struct {
	Handle string
	Email  linkedInEmail `json:"Handle~"`
}
