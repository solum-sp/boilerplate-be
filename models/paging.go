package model
type Paging struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p *Paging) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
}