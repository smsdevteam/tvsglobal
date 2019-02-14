package main

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"log"
	"strconv"
	"time"

	_ "gopkg.in/goracle.v2"

	cm "github.com/smsdevteam/tvsglobal/common"     // db
	st "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
)

//GetKeywordByKeywordID function
func GetKeywordByKeywordID(iKeywordID string) *st.GetKeywordResult {
	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSKeyword"
	l.FunctionName = "GetKeyword"
	l.Request = "KeywordID=" + iKeywordID
	l.Start = t0.Format(time.RFC3339Nano)
	l.InsertappLog("./log/tvskeywordapplog.log", "GetKeyword")

	oRes := st.NewGetKeywordResult()
	var oKeyword st.Keyword

	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer db.Close()
	var statement string
	statement = "begin PK_ICC_KEYWORD.GetKeyword(:0,:1); end;"
	var resultC driver.Rows
	intKeywordID, err := strconv.Atoi(iKeywordID)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 3
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	if _, err := db.Exec(statement, intKeywordID, sql.Out{Dest: &resultC}); err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 4
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resultC.Close()
	values := make([]driver.Value, len(resultC.Columns()))
	for {
		colmap := cm.Createmapcol(resultC.Columns())
		log.Println(colmap)

		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
			resp = err.Error()
			oRes.ErrorCode = 5
			oRes.ErrorDesc = err.Error()
			return oRes
		}

		if values[cm.Getcolindex(colmap, "ALLOWED_KEYWORD_ID")] != nil {
			oKeyword.AllowedKeywordID = values[cm.Getcolindex(colmap, "ALLOWED_KEYWORD_ID")].(int64)
		}
		oKeyword.Attribute = values[cm.Getcolindex(colmap, "ATTRIBUTE")].(string)
		if values[cm.Getcolindex(colmap, "COUNT_VALUE")] != nil {
			oKeyword.CountValue = values[cm.Getcolindex(colmap, "COUNT_VALUE")].(int64)
		}
		if values[cm.Getcolindex(colmap, "CUSTOMER_ID")] != nil {
			oKeyword.CustomerID = values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "DATE_VALUE")] != nil {
			oKeyword.DateValue = values[cm.Getcolindex(colmap, "DATE_VALUE")].(time.Time)
		}
		if values[cm.Getcolindex(colmap, "ID")] != nil {
			oKeyword.ID = values[cm.Getcolindex(colmap, "ID")].(int64)
		}
		oKeyword.KAAttribute = values[cm.Getcolindex(colmap, "KAATTRIBUTE")].(string)
		oKeyword.KAKeyword = values[cm.Getcolindex(colmap, "KAKEYWORD")].(string)
		oKeyword.KALongDescr = values[cm.Getcolindex(colmap, "KALONGDESCR")].(string)
		oKeyword.KTName = values[cm.Getcolindex(colmap, "KTNAME")].(string)
		oKeyword.KTUserKey = values[cm.Getcolindex(colmap, "KTUSERKEY")].(string)
		if values[cm.Getcolindex(colmap, "KEYTYPES_ID")] != nil {
			oKeyword.KeyTypesID = values[cm.Getcolindex(colmap, "KEYTYPES_ID")].(int64)
		}

	}
	oRes.MyKeyword = oKeyword
	if oRes.ErrorCode == 1 {
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "Success"
	}
	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSKeyword"
	l.FunctionName = "GetKeyword"
	l.Request = "KeywordID=" + iKeywordID
	l.Response = resp
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	l.InsertappLog("./log/tvskeywordapplog.log", "GetKeyword")
	//test

	return oRes
}

//GetListKeywordByCustomerID function
func GetListKeywordByCustomerID(iCustomerID string) *st.GetListKeywordResult {
	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSKeyword"
	l.FunctionName = "GetListKeywordByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Start = t0.Format(time.RFC3339Nano)
	l.InsertappLog("./log/tvskeywordapplog.log", "GetListKeywordByCustomerID")

	oRes := st.NewGetListKeywordResult()
	var oListKeyword st.ListKeyword

	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer db.Close()
	var statement string
	statement = "begin PK_ICC_KEYWORD.GetListKeywordByCustomerId(:0,:1); end;"
	var resultC driver.Rows
	intCustomerID, err := strconv.Atoi(iCustomerID)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 3
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	if _, err := db.Exec(statement, intCustomerID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
		resp = err.Error()
		oRes.ErrorCode = 4
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resultC.Close()
	values := make([]driver.Value, len(resultC.Columns()))
	var oLKeyword []st.Keyword
	for {
		colmap := cm.Createmapcol(resultC.Columns())

		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
			resp = err.Error()
			oRes.ErrorCode = 5
			oRes.ErrorDesc = err.Error()
			return oRes
		}

		var oKeyword st.Keyword
		if values[cm.Getcolindex(colmap, "ALLOWED_KEYWORD_ID")] != nil {
			oKeyword.AllowedKeywordID = values[cm.Getcolindex(colmap, "ALLOWED_KEYWORD_ID")].(int64)
		}
		oKeyword.Attribute = values[cm.Getcolindex(colmap, "ATTRIBUTE")].(string)
		if values[cm.Getcolindex(colmap, "COUNT_VALUE")] != nil {
			oKeyword.CountValue = values[cm.Getcolindex(colmap, "COUNT_VALUE")].(int64)
		}
		if values[cm.Getcolindex(colmap, "CUSTOMER_ID")] != nil {
			oKeyword.CustomerID = values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "DATE_VALUE")] != nil {
			oKeyword.DateValue = values[cm.Getcolindex(colmap, "DATE_VALUE")].(time.Time)
		}
		if values[cm.Getcolindex(colmap, "ID")] != nil {
			oKeyword.ID = values[cm.Getcolindex(colmap, "ID")].(int64)
		}
		oKeyword.KAAttribute = values[cm.Getcolindex(colmap, "KAATTRIBUTE")].(string)
		oKeyword.KAKeyword = values[cm.Getcolindex(colmap, "KAKEYWORD")].(string)
		oKeyword.KALongDescr = values[cm.Getcolindex(colmap, "KALONGDESCR")].(string)
		oKeyword.KTName = values[cm.Getcolindex(colmap, "KTNAME")].(string)
		oKeyword.KTUserKey = values[cm.Getcolindex(colmap, "KTUSERKEY")].(string)
		if values[cm.Getcolindex(colmap, "KEYTYPES_ID")] != nil {
			oKeyword.KeyTypesID = values[cm.Getcolindex(colmap, "KEYTYPES_ID")].(int64)
		}

		oLKeyword = append(oLKeyword, oKeyword)
		oListKeyword.Keywords = oLKeyword
	}

	oRes.MyListKeyword = oListKeyword
	if oRes.ErrorCode == 1 {
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "Success"
	}
	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSKeyword"
	l.FunctionName = "GetListKeywordByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Response = resp
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "GetListKeywordByCustomerID")
	//test
	return oRes
}
