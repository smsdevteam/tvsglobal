package common

import (
	"os"
	"time"
)

func writelogfile() {

	filenamea := "tvs-applicationlog"
	currentdate := time.Now()
	filenamea := +currentdate.Format("20060102")
	f, err := os.OpenFile("d:/"+filenamea, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString("new data that wasn't there originally\n"); err != nil {
		panic(err)
	}

}
