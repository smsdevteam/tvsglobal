package structs

import (
	"time"
)

//Note is Note
type Note struct {
	CustomerID        int32
	CreatedByUserID   int32
	ActionUserID      int32
	CategoryID        string
	CompletionStageID string
	Body              string
	NoteID            int32
	CreateDateTime    time.Time
}
