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
