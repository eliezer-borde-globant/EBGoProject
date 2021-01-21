package controller

type UpdateParams struct {
	Repo    string                              `json:"repo" xml:"repo" form:"repo"`
	Owner   string                              `json:"owner" xml:"owner" form:"owner"`
	Changes map[string][]map[string]interface{} `json:"changes" xml:"changes" form:"changes"`
}

type CreateParams struct {
	Repo    string `json:"repo" xml:"repo" form:"repo"`
	Owner   string `json:"owner" xml:"owner" form:"owner"`
	Content string `json:"content" xml:"content" form:"content"`
}
