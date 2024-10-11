package model

type XMLResponse struct {
	Status   int    `xml:"Status"`
	Message  string `xml:"Message"`
	Resource string `xml:"Resource"`
}
