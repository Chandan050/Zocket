package models

type Request struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Value string `json:"value"`
}
