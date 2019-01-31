package tvsstructs

// ResponseResult Obj
type ResponseResult struct {
	ErrorCode    int    `json:"errorcode"`
	ErrorDesc    string `json:"errordesc"`
	CustomNum    int    `json:"customnum"`
	CustomString string `json:"customstring"`
}

// NewResponseResult Obj
func NewResponseResult() *ResponseResult {
	return &ResponseResult{
		ErrorCode:    1,
		ErrorDesc:    "Unexpected Error",
		CustomNum:    0,
		CustomString: "",
	}
}
