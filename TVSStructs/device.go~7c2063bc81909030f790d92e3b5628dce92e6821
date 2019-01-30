package tvsstructs

// Device Obj
type DeviceInfo struct {
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
	FromStockHanderID     int64  `json:"fromdepotid"`
	ToStockHanderID       int64  `json:"todepotid"`
	IBSModelID            int64  `json:"modelid"`
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
	ProcessRes ResponseResult `json:"processres"`
	NewSNRes   []CreateSNRes  `json:"newsnres"`
}

// CreateSNRes Obj
type CreateSNRes struct {
	SerialNumber string `json:"serialno"`
	ResultCode   int    `json:"resultcode"`
	ResultDesc   string `json:"resultdesc"`
}

// ReceiveExchangeDeviceReq Obj
type ReceiveExchangeDeviceReq struct {
	DeviceID               int64  `json:"deviceid"`
	StockHandlerID         int64  `json:"stockhandlerid"`
	PalletID               int64  `json:"palletid"`
	ReasonID               int64  `json:"reasonid"`
	DeviceExchangeReasonID int64  `json:"deviceexchangereasonid"`
	ShipDate               string `json:"shipdate"`
}

// Device Obj
type Device struct {
	CaReferenceNumber int `json:"CaReferenceNumber"`
	//CustomFields
	ExternalID            string `json:"externalid"`
	FinanceOptionID       int    `json:"financeoptionid"`
	FromStockHandlerID    int    `json:"fromstockhandlerid"`
	ID                    int    `json:"id"`
	MACAddress1           string `json:"macaddress1"`
	MACAddress2           string `json:"macaddress2"`
	ModelID               int    `json:"modelid"`
	OrderID               int    `json:"orderid"`
	PalletID              int    `json:"palletid"`
	Provisioned           bool   `json:"provisioned"`
	SerialNumber          string `json:"serialnumber"`
	ShipDate              string `json:"shipdate"`
	StatusID              int    `json:"statusid"`
	StockHandlerID        int    `json:"stockhandlerid"`
	StockReceiveDetailsID int    `json:"stockreceivedetailsid"`
	WarrantyEndDate       string `json:"warrantyenddate"`
}

// GetDeviceResponse Obj
type GetDeviceResponse struct {
	ProcessResult ResponseResult
	DeviceResult  Device
}
