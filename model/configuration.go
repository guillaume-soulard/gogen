package model

type Configuration struct {
	TemplateConfiguration interface{} `json:"template"`
	Options               Options     `json:"options"`
}

type Options struct {
	Seed       string `json:"seed"`
	Amount     int    `json:"amount"`
	Skip       int    `json:"skip"`
	Alias      map[string]Configuration
	Generation Generation `json:"generation"`
}

type Generation struct {
	Type    string            `json:"type"`
	Options GenerationOptions `json:"options"`
}

type GenerationOptions struct {
	Interval int64 `json:"interval"`
}
