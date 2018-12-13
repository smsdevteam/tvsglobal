package tvsstructs

import (
	"time"
)

type Note struct {
	CustomerID        int       `json:"customerid"`
	CreatedByUserID   int       `json:"createdbyuserid"`
	ActionUserID      int       `json:"actionuserid"`
	CategoryID        string    `json:"categoryid"`
	CompletionStageID string    `json:"completionstageid"`
	Body              string    `json:"body"`
	NoteID            int       `json:"noteid"`
	CreateDateTime    time.Time `json:"createdatetime"`
}
