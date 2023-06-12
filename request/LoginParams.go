package request

type LoginParams struct {
	Mobile int `json:"mobile"`
	Code   int `json:"code"`
}
