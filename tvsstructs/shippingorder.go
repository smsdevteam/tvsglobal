package tvsstructs

// import (
// 	"time"
// )

// ShippingOrderReq Header Obj
type ShippingOrderReq struct {
	ID                   int64                  `json:"id"`
	Comments             string                 `json:"comments"`
	CustomerID           int64                  `json:"customerid"`
	ShipFromStockhandler int64                  `json:"shipfromstockhandler"`
	Status               int                    `json:"status"`
	ShippingMethod       int64                  `json:"shippingmethod"`
	OrderType            int64                  `json:"ordertype"`
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

// ShippingOrderData Obj
type ShippingOrderData struct {
	AgreementID    int64  `xml:"agreementId" json:"agreementId"`
	Comment        string `xml:"Comment" json:"Comment"`
	CreateDateTime string `xml:"CreateDateTime" json:"CreateDateTime"`
	CustomFields   struct {
		CustomFields []CustomFieldValue
	} `xml:"CustomFields,omitempty" json:"CustomFields,omitempty"`
	CustomerID                int64  `xml:"CustomerId" json:"CustomerId"`
	Destination               string `xml:"Destination" json:"Destination"`
	Extended                  string `xml:"Extended" json:"Extended"`
	FinancialAccountID        int64  `xml:"FinancialAccountId" json:"FinancialAccountId"`
	FullyReceiveReturnedOrder bool   `xml:"FullyReceiveReturnedOrder" json:"FullyReceiveReturnedOrder"`
	ID                        int64  `xml:"Id" json:"Id"`
	IgnoreAgreementID         bool   `xml:"IgnoreAgreementId" json:"IgnoreAgreementId"`
	OldStatusID               int64  `xml:"OldStatusId" json:"OldStatusId"`
	ParentOrderID             int64  `xml:"ParentOrderId" json:"ParentOrderId"`
	ReceivedQuantity          int64  `xml:"ReceivedQuantity" json:"ReceivedQuantity"`
	Reference                 string `xml:"Reference" json:"Reference"`
	ReturnedQuantity          int64  `xml:"ReturnedQuantity" json:"ReturnedQuantity"`
	SandboxID                 int64  `xml:"SandboxId" json:"SandboxId"`
	SandboxSkipValidation     bool   `xml:"SandboxSkipValidation" json:"SandboxSkipValidation"`
	ShipByDate                string `xml:"ShipByDate" json:"ShipByDate"`
	ShipFromStockHandlerID    int64  `xml:"ShipFromStockHandlerId" json:"ShipFromStockHandlerId"`
	ShipToAddressID           int64  `xml:"ShipToAddressId" json:"ShipToAddressId"`
	ShipToPartyID             int64  `xml:"ShipToPartyId" json:"ShipToPartyId"`
	ShipToPostalCode          string `xml:"ShipToPostalCode" json:"ShipToPostalCode"`
	ShippedDate               string `xml:"ShippedDate" json:"ShippedDate"`
	ShippedQuantity           int64  `xml:"ShippedQuantity" json:"ShippedQuantity"`
	ShippingMethodID          int64  `xml:"ShippingMethodId" json:"ShippingMethodId"`
	ShippingOrderLines        struct {
		Items struct {
			ShippingOrderLine []ShippingOrderLineStruct
		} `xml:"Items" json:"Items"`
		More       bool  `xml:"More" json:"More"`
		Page       int64 `xml:"Page" json:"Page"`
		TotalCount int64 `xml:"TotalCount" json:"TotalCount"`
	} `xml:"ShippingOrderLines" json:"ShippingOrderLines"`
	StatusID        int64 `xml:"StatusId" json:"StatusId"`
	TotalQuantity   int64 `xml:"TotalQuantity" json:"TotalQuantity"`
	TrackingNumbers struct {
		Items struct {
			TrackingNumbers []TrackingNumber
		} `xml:"Items" json:"Items"`
		More       bool  `xml:"More" json:"More"`
		Page       int64 `xml:"Page" json:"Page"`
		TotalCount int64 `xml:"TotalCount" json:"TotalCount"`
	} `xml:"TrackingNumbers,omitempty" json:"TrackingNumbers,omitempty"`
	TypeID int64 `xml:"TypeId" json:"TypeId"`
}

// ShippingOrderLineStruct Obj
type ShippingOrderLineStruct struct {
	AgreeementDetailID        int64 `xml:"AgreeementDetailId" json:"AgreeementDetailId"`
	CorrelatedHardwareModelID int64 `xml:"CorrelatedHardwareModelId" json:"CorrelatedHardwareModelId"`
	CustomFields              struct {
		CustomFields []CustomFieldValue
	} `xml:"CustomFields,omitempty" json:"CustomFields,omitempty"`
	DevicePerAgreementDetailID int64  `xml:"DevicePerAgreementDetailId" json:"DevicePerAgreementDetailId"`
	Extended                   string `xml:"Extended" json:"Extended"`
	ExternalID                 string `xml:"ExternalId" json:"ExternalId"`
	FinanceOptionID            int64  `xml:"FinanceOptionId" json:"FinanceOptionId"`
	HardwareModelID            int64  `xml:"HardwareModelId" json:"HardwareModelId"`
	ID                         int64  `xml:"Id" json:"Id"`
	NonSubstitutableModel      bool   `xml:"NonSubstitutableModel" json:"NonSubstitutableModel"`
	OrderLineNumber            int64  `xml:"OrderLineNumber" json:"OrderLineNumber"`
	Quantity                   int64  `xml:"Quantity" json:"Quantity"`
	ReceivedQuantity           int64  `xml:"ReceivedQuantity" json:"ReceivedQuantity"`
	ReturnedQuantity           int64  `xml:"ReturnedQuantity" json:"ReturnedQuantity"`
	SandboxID                  int64  `xml:"SandboxId" json:"SandboxId"`
	SandboxSkipValidation      bool   `xml:"SandboxSkipValidation" json:"SandboxSkipValidation"`
	SerializedStock            bool   `xml:"SerializedStock" json:"SerializedStock"`
	ShippingOrderID            int64  `xml:"ShippingOrderId" json:"ShippingOrderId"`
	TechnicalProductID         int64  `xml:"TechnicalProductId" json:"TechnicalProductId"`
	TotalLinkedDevices         int64  `xml:"TotalLinkedDevices" json:"TotalLinkedDevices"`
	TotalUnlinkedDevices       int64  `xml:"TotalUnlinkedDevices" json:"TotalUnlinkedDevices"`
}

// TrackingNumber Obj
type TrackingNumber struct {
	Extended        string `xml:"Extended" json:"Extended"`
	ID              int64  `xml:"Id" json:"Id"`
	Number          string `xml:"Number" json:"Number"`
	ShippingOrderID int64  `xml:"ShippingOrderId" json:"ShippingOrderId"`
}

// CustomFieldValue Obj
type CustomFieldValue struct {
	Extended string `xml:"Extended" json:"Extended"`
	ID       int64  `xml:"Id" json:"Id"`
	Name     string `xml:"Name" json:"Name"`
	Sequence int64  `xml:"Sequence" json:"Sequence"`
	Value    string `xml:"Value" json:"Value"`
}

// ShippingOrderDataReq Obj
type ShippingOrderDataReq struct {
	SODetail   ShippingOrderData `xml:"SODetail" json:"SODetail"`
	Reasonnr   int64             `xml:"Reasonnr" json:"Reasonnr"`
	ByUsername string            `xml:"ByUsername" json:"ByUsername"`
}

// SOResult Obj
type SOResult struct {
	ProcessResult ResponseResult
	SODetail      ShippingOrderData
}
