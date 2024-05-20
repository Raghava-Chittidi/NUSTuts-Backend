package api

type Response struct {
	Data interface{} `json:"data,omitempty"`
	Message string `json:"message"`
	Error error `json:"error"`
}