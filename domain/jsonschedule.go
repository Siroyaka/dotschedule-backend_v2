package domain

type JsonSchedule struct {
	VideoID      string          `json:"VideoID"`
	VideoTitle   string          `json:"VideoTitle"`
	VideoLink    string          `json:"VideoLink"`
	VideoStatus  int             `json:"VideoStatus"`
	StreamerName string          `json:"StreamerName"`
	StreamerID   string          `json:"StreamerID"`
	StartDate    string          `json:"StartDate"`
	Thumbnail    string          `json:"Thumbnail"`
	Duration     int             `json:"Duration"`
	Participants map[string]bool `json:"Participants"`
}

type JsonScheduleList []JsonSchedule
