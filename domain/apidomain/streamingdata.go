package apidomain

type StreamingData struct {
	ID          string `json:"ID"`
	URL         string `json:"URL"`
	Platform    string `json:"Platform"`
	Status      int    `json:"Status"`
	StartDate   string `json:"StartDate"`
	Duration    int    `json:"Duration"`
	Thumbnail   string `json:"Thumbnail"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}
