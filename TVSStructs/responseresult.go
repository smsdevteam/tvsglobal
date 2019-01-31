package tvsstructs

// Response Obj
type ResponseResult struct {
	ErrorCode    int    `json:"errorcode"`
	ErrorDesc    string `json:"errordesc"`
	CustomNum    int    `json:"customnum"`
	CustomString string `json:"customstring"`
}
