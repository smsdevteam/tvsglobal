package tvsstructs

// ShippingOrderReq Header Obj
type ShippingOrderReq struct {
	ID                   int64                  `json:"id"`
	Comments             string                 `json:"comments"`
	CustomerID           int64                  `json:"customerid"`
	ShipByDate           string                 `json:"shipbydate"`
	ShipFromStockhandler int64                  `json:"shipfromstockhandler"`
	Status               int                    `json:"status"`
	ShippingMethod       int64                  `json:"shippingmethod"`
	OrderType            int64                  `json:"ordertype"`
	FinancialAccount     int64                  `json:"financialaccount"`
	ShippedDate          string                 `json:"shippeddate"`
	CreateDateTime       string                 `json:"createdatetime"`
	Reference            string                 `json:"reference"`
	ExternalAgent        string                 `json:"externalagent"`
	ShippingOrderLines   []ShippingOrderLineReq `json:"shippingorderlines"`
	ShippingDevices      []ShippingDeviceReq    `json:"shippingdevices"`
}

// ShippingOrderLineReq Obj
type ShippingOrderLineReq struct {
	ShippingOrderID         int64 `json:"shippingorderid"`
	ShippingOrderLineID     int64 `json:"shippingorderlineid"`
	DeviceAgreementDetailID int64 `json:"deviceagreementdetailid"`
	AgreementDetailID       int64 `json:"agreementdetailid"`
	TechnicalProductID      int64 `json:"technicalproductid"`
	OrderLineNumber         int64 `json:"orderlinenumber"`
	Quantity                int64 `json:"quantity"`
	HardwareModelID         int64 `json:"hardwaremodelid"`
	FinanceOptionID         int64 `json:"financeoptionid"`
}

// ShippingDeviceReq Obj
type ShippingDeviceReq struct {
	ShippingOrderLineID int64  `json:"shippingorderlineid"`
	SerialNumber        string `json:"sn"`
}

// ShippingOrderRes Header Obj
type ShippingOrderRes struct {
	ID                 int64                  `json:"id"`
	DepotFrom          int64                  `json:"depotfrom"`
	DepotTo            int64                  `json:"depotto"`
	StatusID           int64                  `json:"statusid"`
	StatusDesc         string                 `json:"statusname"`
	TypeID             int64                  `json:"typeid"`
	TypeDesc           string                 `json:"typename"`
	CreateComments     string                 `json:"createcomments"`
	CreateReference    string                 `json:"createreference"`
	CreateDateTime     string                 `json:"createdatetime"`
	CreateBy           int64                  `json:"createby"`
	CreateByName       string                 `json:"createbyname"`
	ShippingOrderLines []ShippingOrderLineRes `json:"shippingorderlines"`
	ShippingDevices    []ShippingDeviceRes    `json:"shippingdevices"`
}

// ShippingOrderLineRes Obj
type ShippingOrderLineRes struct {
	LineID          int64  `json:"lineid"`
	ShippingOrderID int64  `json:"shippingorderid"`
	LineNr          int64  `json:"linenr"`
	ProductID       int64  `json:"productid"`
	ProductKey      string `json:"productkey"`
	ModelID         int64  `json:"modelid"`
	ModelKey        string `json:"modelkey"`
	Qty             int64  `json:"qty"`
}

// ShippingDeviceRes Obj
type ShippingDeviceRes struct {
	ShippingOrderID int64  `json:"shippingorderid"`
	LineID          int64  `json:"lineid"`
	SerialNumber    string `json:"sn"`
	StatusID        int64  `json:"statusid"`
	DVResult        string `json:"deviceresult"`
}
