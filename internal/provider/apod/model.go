package apod

type metaData struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Date        string `json:"date"`
	URL         string `json:"url"`
	Copyright   string `json:"Copyright"`
}
