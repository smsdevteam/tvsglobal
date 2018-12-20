package common

import (
	"os"
	"time"
)

//Writelogfile is write log file
func Writelogfile() {

	filenamea := "tvs-applicationlog"
	currentdate := time.Now()
	filenamea = filenamea + currentdate.Format("20060102") + ".txt"
	f, err := os.OpenFile("d:/"+filenamea, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString("new data that wasn't there originally\n"); err != nil {
		panic(err)
	}

}
