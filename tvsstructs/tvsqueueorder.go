package tvsstructs

import "time"

// TVSSubmitOrdReqData object  call from client
type TVSSubmitOrdReqData struct {
	Orderid     string
	OrderType   string
	ChannelCode string
	OrderDate   time.Time
	TVSCustNo   int64
	Custinfo    CustomerInfo
}

// TVSSubmitOrdResData object for response to client
type TVSSubmitOrdResData struct {
	Orderid           string
	Trackingno        string
	ResponseResultobj ResponseResult
}

// TVSSubmitOrderProcess object for send to rabbit mq
type TVSSubmitOrderProcess struct {
	Orderdata   TVSSubmitOrderData
	TVSTaskList []TVSTaskinfo
}

// TVSSubmitOrderData object for send to rabbit mq
type TVSSubmitOrderData struct {
	Trackingno string
	TVSOrdReq  TVSSubmitOrdReqData
}

// Tvsmsresp object for task process info
type Tvsmsresp struct {
	Resultobj ResponseResult
}

// TVSTaskinfo object for task process info
type TVSTaskinfo struct {
	Taskid    string
	Taskname  string
	MSname    string
	Servurl   string
	Responseobjname string
	Resultobj ResponseResult
}
