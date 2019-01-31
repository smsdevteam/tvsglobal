package common

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

//Applog strucure
type Applog struct {
	OrderNo         string `json:"orderno"`
	TrackingNo      string `json:"trackingno"`
	ApplicationName string `json:"applicationname"`
	FunctionName    string `json:"functionname"`
	OrderDate       string `json:"orderdate"`
	TVSNo           string `json:"tvsno"`
	MobileNo        string `json:"mobileno"`
	SerialNo        string `json:"serialno"`
	Reference1      string `json:"reference1"`
	Reference2      string `json:"reference2"`
	Reference3      string `json:"reference3"`
	Reference4      string `json:"reference4"`
	Reference5      string `json:"reference5"`
	Request         string `json:"request"`
	Response        string `json:"response"`
	Start           string `json:"start"`
	End             string `json:"end"`
	Duration        string `json:"duration"`
}

type Processconfig struct {
	CallFunction string `json:"callfunction"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Duration     string `json:"duration"`
	ResultCode   string `json:"resultcode"`
	ResultDesc   string `json:"resultdesc"`
}

//Workflowlog strucure
type Workflowlog struct {
	OrderNo         string          `json:"orderno"`
	TrackingNo      string          `json:"trackingno"`
	ApplicationName string          `json:"applicationname"`
	FunctionName    string          `json:"functionname"`
	OrderDate       string          `json:"orderdate"`
	TVSNo           string          `json:"tvsno"`
	MobileNo        string          `json:"mobileno"`
	SerialNo        string          `json:"serialno"`
	Reference1      string          `json:"reference1"`
	Reference2      string          `json:"reference2"`
	Reference3      string          `json:"reference3"`
	Reference4      string          `json:"reference4"`
	Reference5      string          `json:"reference5"`
	InputData       string          `json:"InputData"`
	Start           string          `json:"start"`
	End             string          `json:"end"`
	Duration        string          `json:"duration"`
	ProcessConfig   []Processconfig `json:"processconfig"`
}

func (a *Applog) InsertappLog(logfile string, msg string) error {
	// open a file
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	log.WithFields(log.Fields{
		"OrderNo":         a.OrderNo,
		"TrackingNo":      a.TrackingNo,
		"ApplicationName": a.ApplicationName,
		"FunctionName":    a.FunctionName,
		"OrderDate":       a.OrderDate,
		"TVSNo":           a.TVSNo,
		"MobileNo":        a.MobileNo,
		"SerialNo":        a.SerialNo,
		"Reference1":      a.Reference1,
		"Reference2":      a.Reference2,
		"Reference3":      a.Reference3,
		"Reference4":      a.Reference4,
		"Reference5":      a.Reference5,
		"Request":         a.Request,
		"Response":        a.Response,
		"Start":           a.Start,
		"End":             a.End,
		"Duration":        a.Duration,
	}).Info(msg)

	return nil
}

func (a *Workflowlog) Insertworkflowlog(logfile string) error {
	// open a file
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	log.WithFields(log.Fields{
		"OrderNo":         a.OrderNo,
		"TrackingNo":      a.TrackingNo,
		"ApplicationName": a.ApplicationName,
		"FunctionName":    a.FunctionName,
		"OrderDate":       a.OrderDate,
		"TVSNo":           a.TVSNo,
		"MobileNo":        a.MobileNo,
		"SerialNo":        a.SerialNo,
		"Reference1":      a.Reference1,
		"Reference2":      a.Reference2,
		"Reference3":      a.Reference3,
		"Reference4":      a.Reference4,
		"Reference5":      a.Reference5,
		"Input Data":      a.InputData,
		"Start":           a.Start,
		"End":             a.End,
		"Duration":        a.Duration,
		"ProcessConfig":   a.ProcessConfig,
	}).Info("")
	return nil
}
