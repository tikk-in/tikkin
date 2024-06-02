package dto

// ResultList Result list (generics)
type ResultList struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}
