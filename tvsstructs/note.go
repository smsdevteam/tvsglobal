package tvsstructs

import (
	"time"
)

type Note struct {
	CustomerID        int64     `json:"customerid"`
	CreatedByUserID   int64     `json:"createdbyuserid"`
	ActionUserID      int64     `json:"actionuserid"`
	CategoryID        string    `json:"categoryid"`
	CompletionStageID string    `json:"completionstageid"`
	Body              string    `json:"body"`
	NoteID            int64     `json:"noteid"`
	CreateDateTime    time.Time `json:"createdatetime"`
}

type GetNote struct {
	CustomerID        string `json:"customerid"`
	CreatedByUserID   string `json:"createdbyuserid"`
	ActionUserID      string `json:"actionuserid"`
	CategoryID        string `json:"categoryid"`
	CompletionStageID string `json:"completionstageid"`
	Body              string `json:"body"`
	NoteID            string `json:"noteid"`
	CreateDateTime    string `json:"createdatetime"`
}

type GetNoteResult struct {
	MyNote    Note   `json:"note"`
	ErrorCode int    `json:"errorcode"`
	ErrorDesc string `json:"errordesc"`
}

// NewGetNoteResult Obj
func NewGetNoteResult() *GetNoteResult {
	return &GetNoteResult{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}

type ListNote struct {
	Notes []Note `json:"notes"`
}

type GetListNoteResult struct {
	MyListNote ListNote `json:"getlistnoteresult"`
	ErrorCode  int      `json:"errorcode"`
	ErrorDesc  string   `json:"errordesc"`
}

// NewGetListNoteResult Obj
func NewGetListNoteResult() *GetListNoteResult {
	return &GetListNoteResult{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}

type CreateNoteRequest struct {
	InNote struct {
		ActionUserKey      int64     `json:"ActionUserKey"`
		CustomerID         int64     `json:"CustomerId"`
		CreatedByUserID    int64     `json:"CreatedByUserId"`
		CategoryKey        string    `json:"CategoryKey"`
		CompletionStageKey string    `json:"CompletionStageKey"`
		Body               string    `json:"Body"`
		NoteID             int64     `json:"Id"`
		CreateDate         time.Time `json:"CreateDate"`
		Extended           string    `json:"Extended"`
	} `json:"inNote"`
	InReason int64 `json:"InReason"`
	ByUser   struct {
		ByUser    string `json:"byUser"`
		ByChannel string `json:"byChannel"`
		ByProject string `json:"byProject"`
		ByHost    string `json:"byHost"`
	} `json:"byUser"`
}

type CreateNoteResponse struct {
	ErrorCode   int    `json:"errorcode"`
	ErrorDesc   string `json:"errordesc"`
	ResultValue string `json:"resultvalue"`
}

// NewCreateNoteResponse Obj
func NewCreateNoteResponse() *CreateNoteResponse {
	return &CreateNoteResponse{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}

type UpdateNoteRequest struct {
	InNote struct {
		ActionUserKey      int64     `json:"ActionUserKey"`
		CustomerID         int64     `json:"CustomerId"`
		CreatedByUserID    int64     `json:"CreatedByUserId"`
		CategoryKey        string    `json:"CategoryKey"`
		CompletionStageKey string    `json:"CompletionStageKey"`
		Body               string    `json:"Body"`
		NoteID             int64     `json:"Id"`
		CreateDate         time.Time `json:"CreateDate"`
		Extended           string    `json:"Extended"`
	} `json:"inNote"`
	InReason int64 `json:"InReason"`
	ByUser   struct {
		ByUser    string `json:"byUser"`
		ByChannel string `json:"byChannel"`
		ByProject string `json:"byProject"`
		ByHost    string `json:"byHost"`
	} `json:"byUser"`
}

type UpdateNoteResponse struct {
	ErrorCode   int    `json:"errorcode"`
	ErrorDesc   string `json:"errordesc"`
	ResultValue string `json:"resultvalue"`
}

// NewUpdateNoteResponse Obj
func NewUpdateNoteResponse() *UpdateNoteResponse {
	return &UpdateNoteResponse{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}
