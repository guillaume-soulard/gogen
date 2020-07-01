package model

type Configuration struct {
	TemplateConfiguration interface{} `json:"template"`
	Options               Options     `json:"options"`
}

type Options struct {
	Seed   string `json:"seed"`
	Amount int    `json:"amount"`
	Skip   int    `json:"skip"`
	Alias  map[string]Configuration
}
