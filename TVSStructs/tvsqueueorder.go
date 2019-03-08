package tvsstructs

import "time"

// TVSQueuSubmitOrderRequest object  call from client
type TVSQueuSubmitOrderRequest struct {
	Orderid     string
	OrderType   string
	ChannelCode string
	OrderDate   time.Time
	TVSCustNo   int64
	Custinfo    CustomerInfo
}

// TVSSubmitOrderToQueue object for send to rabbit mq
type TVSSubmitOrderToQueue struct {
	Trackingno  string
	TVSOrdReq   TVSQueuSubmitOrderRequest
	TVSTaskList []TVSTaskinfo
}

// TVSQueueSubmitOrderReponse object for response to client
type TVSQueueSubmitOrderReponse struct {
	Orderid           string
	Trackingno        string
	ResponseResultobj ResponseResult
}

// TVSTaskinfo object for task process info
type TVSTaskinfo struct {
	Taskid    string
	Taskname  string
	MSname    string
	Resultobj ResponseResult
}
