package main

import (
	"fmt"
	"time"
)

type Note struct {
	customerID        int32
	createdByUserID   int32
	actionUserID      int32
	categoryID        string
	completionStageID string
	body              string
	noteID            int32
	createDateTime    time.Time
}

func main() {
	fmt.Println("test")
}
