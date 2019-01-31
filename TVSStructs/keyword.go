package tvsstructs

type Keyword struct {
	Attribute     string `json:"Attribute"`
	Description   string `json:"Description"`
	Hide          int32  `json:"Hide"`
	ID            int64  `json:"Id"`
	Keyword       string `json:"Keyword"`
	KeywordTypeID int64  `json:"KeywordTypeId"`
	Name          string `json:"Name"`
}
