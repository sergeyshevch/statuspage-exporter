package statuspageio

// AtlassianStatusPageResponse is a response from Atlassian StatusPage API.
type AtlassianStatusPageResponse struct {
	Page struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"page"`
	Components []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
	} `json:"components"`
}
