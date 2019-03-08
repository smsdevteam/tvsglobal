package tvsstructs

import "time"

type TVSQueuSubmitOrderRequest struct {
	Orderid     string
	OrderType   string
	ChannelCode string
	OrderDate   time.Time
	TVSCustNo   int64
	Custinfo    CustomerInfo
}
type TVSSubmitOrderToQueue struct {
	Trackingno        string
	TVSOrdReq TVSQueuSubmitOrderRequest
	
}

type TVSQueueSubmitOrderReponse struct {
	Orderid           string
	Trackingno        string
	ResponseResultobj ResponseResult
}
