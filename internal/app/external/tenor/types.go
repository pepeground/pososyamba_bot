package tenor

type Size string

const (
	Gif Size = "gif"
)

type Response struct {
	Results []Result `json:"results"`
}

type Result struct {
	Tags  []string `json:"tags,omitempty"`
	Url   string   `json:"url,omitempty"`
	Media []Media  `json:"media,omitempty"`
}

type Media struct {
	Gif struct {
		Url     string `json:"url"`
		Dims    []int  `json:"dims"`
		Preview string `json:"preview"`
		Size    int    `json:"size"`
	} `json:"gif"`
}
