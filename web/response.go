package web

const STATUS_FAIL = "fail"
const STATUS_SUCCESS = "success"

type Response struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   any    `json:"data"`
	Error  Error  `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
	Detail  any    `json:"details"`
}

type ErrorDetail struct {
	Path    string `json:"path"`
	Value   string `json:"value"`
	Message string `json:"message"`
}
