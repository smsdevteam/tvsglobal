package tvsstructs

// Device Obj
type Device struct {
	DeviceID            int64  `json:"deviceid"`
	SerialNumber        string `json:"sn"`
	StatusID            int64  `json:"statusid"`
	StatusDesc          string `json:"statusdesc"`
	ModelID             int64  `json:"modelid"`
	ModelDesc           string `json:"modeldesc"`
	ProductID           int64  `json:"productid"`
	ProductDesc         string `json:"productdesc"`
	StockhandlerID      int64  `json:"stockhandlerid"`
	StockhandlerDesc    string `json:"stockhandlerdesc"`
	AllowSystem         string `json:"allowsystem"`
	FactoryWarrantyDate string `json:"factorywrdate"`
	AgentKey            string `json:"agentkey"`
	AgentName           string `json:"agentname"`
	ReturnDate          string `json:"returndate"`
}

// StockReceiveDetails Obj
type StockReceiveDetails struct {
	StockReceiveDetailsID int64  `json:"id"`
	BatchNumber           string `json:"batchnumber"`
	FromStockHanderId     int64  `json:"fromdepotid"`
	ToStockHanderId       int64  `json:"todepotid"`
	IBSModelId            int64  `json:"modelid"`
	WarrantyEndDate       string `json:"wrenddate"`
}

// NewDeviceReq Obj
type NewDeviceReq struct {
	StockReceiveDetail StockReceiveDetails `json:"stockreceivedetails"`
	SerialNumber       []string            `json:"serialno"`
	Reason             int64               `json:"reason"`
	ByUser             string              `json:"byusername"`
}

// NewDeviceRes Obj
type NewDeviceRes struct {
	SerialNumber string `json:"serialno"`
	ResultCode   int    `json:"resultcode"`
	ResultDesc   string `json:"resultdesc"`
}
