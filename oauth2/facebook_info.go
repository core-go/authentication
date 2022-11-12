package oauth2

type facebookInfo struct {
	Id        string
	Email     string
	FirstName string `json:"First_name"`
	LastName  string `json:"Last_name"`
	Name      string
	Gender    string
	Picture   facebookPicture `json:"Picture"`
}
type facebookPicture struct {
	Data facebookData `json:"Data"`
}
type facebookData struct {
	Height       int
	IsSilhouette bool `json:"Is_silhouette"`
	Url          string
	Width        int
}
