package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"

	//"fmt"
	"io"
	//"strconv"
	"time"
	cm "tvsglobal/common"
	tg "tvsglobal/tvsccglobal"
	st "tvsglobal/tvsstructs"

	//"github.com/jmoiron/sqlx"
	"net/http"

	_ "gopkg.in/goracle.v2"
)

func preparejobdata(customerid int, ccbsfn string) string {
	var trnseqno string
	cm.ExcutestoreDS("ICC", ` begin :result := tvs_goccbsbn.create_omxjob(:p_custno,
	 		 :p_ccbsfn,:p_shno,:p_sheventno,:p_processflag,
	 		 :p_allowsendomx,:p_ccbs_activityreason); end;`, sql.Out{Dest: &trnseqno}, customerid, ccbsfn, 0, 0, 1, 1, "CREQ")
	return trnseqno
}
func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}
func createmapcol(data []string) map[string]int {
	var colmap = map[string]int{}

	for k, v := range data {
		colmap[v] = k
	}
	return colmap
}
func getjobinfo(trnseqno string) st.TVSBNProperty {
	var TVSBNP st.TVSBNProperty
	var resultI driver.Rows
	var err error
	cm.ExcutestoreDS("ICC", ` begin  tvs_goccbsbn.getjobinfo(:p_transeq,
		:p_rs); end;`, trnseqno, sql.Out{Dest: &resultI})
	defer resultI.Close()
	values := make([]driver.Value, len(resultI.Columns()))
	colmap := createmapcol(resultI.Columns())
	for {
		print(colmap)
		err = resultI.Next(values)
		if err == nil {
			if err == io.EOF {
				break
			}
		} else {
			break
		}
		TVSBNP.TRNSEQNO = values[cm.Getcolindex(colmap, "TRNSEQNO")].(string) // values[colmap["TRNSEQNO"]].(string)
		TVSBNP.CCBSorderno = values[colmap["CCBS_ORDERNO"]].(string)
		//custno, err := strconv.Atoi(values[colmap["TVS_CUSTOMERNO"]].(string)) // strconv.ParseInt( values[2].(string), 10, 64)
		//if err != nil {
		//	TVSBNP.TVSCUSTOMERNO = custno
		//	}

		TVSBNP.TVSAccountno = values[colmap["TVS_ACCOUNTNO"]].(string)
		TVSBNP.SHNO = values[colmap["SHNO"]].(string)
		TVSBNP.CCBSFN = values[colmap["CCBSFN"]].(string)
		TVSBNP.RECALLOFFER = values[colmap["RECALLOFFER"]].(string)
		TVSBNP.HAVEOCCHARE = values[colmap["HAVEOCCHARE"]].(string)
		TVSBNP.OMXOrderType = values[colmap["OMXORDERTYPE"]].(string)
		TVSBNP.TVSCustomerType = values[colmap["TVS_CUSTOMERTYPE"]].(string)
		TVSBNP.CCBSACTIVITYREON = values[colmap["CCBS_ACTIVITYREASON"]].(string)
		TVSBNP.CCBSUSERTEXT = values[colmap["CCBS_USERTEXT"]].(string)
		TVSBNP.CCBSCustomerno = values[colmap["CCBS_CUSTOMERNO"]].(string)
		TVSBNP.CCBSAccountno = values[colmap["CCBS_ACCOUNTNO"]].(string)
		TVSBNP.CCBSOUNo = values[colmap["CCBS_OUNO"]].(string)
		TVSBNP.CCBSSubNo = values[colmap["CCBS_SUBNO"]].(string)
		TVSBNP.CCBSURLSERVICE = values[colmap["CCBSURL_SERVICE"]].(string)
		TVSBNP.Refvalue = values[colmap["REFVALUE"]].(string)
		TVSBNP.Reftype = values[colmap["REFTYPE"]].(string)
		TVSBNP.ToCCBSCustomerno = values[colmap["TO_CCBS_CUSTOMERNO"]].(string)
		TVSBNP.OLDCCBSACCOUNTNO = values[colmap["OLD_CCBSACCOUNTNO"]].(string)
		TVSBNP.Oldccbssubno = values[colmap["OLD_CCBSSUBNO"]].(string)
		TVSBNP.AddSOCLevelOU = values[colmap["ADDSOC_LEVEL_OU"]].(string)
		TVSBNP.TVSBNOMXPropertyobj.Channel = values[colmap["CHANNEL"]].(string)
		TVSBNP.TVSBNOMXPropertyobj.DealerCode = values[colmap["DEALERCODE"]].(string)
		TVSBNP.TVSBNOMXPropertyobj.EffectiveDateSpecified = cm.StrToInt(values[colmap["EFFECTIVEDATESPECIFIED"]].(string))
		TVSBNP.TVSBNOMXPropertyobj.EffectiveDate = values[colmap["EFFECTIVEDATE"]].(string)

		TVSBNP.TVSBNCCBSOfferPropertylist = getccbsoffer(TVSBNP)
	}

	return TVSBNP
}
func getccbsoffer(TVSBNP st.TVSBNProperty) []st.TVSBNCCBSOfferProperty {
	var resultI driver.Rows
	var TVSBNCCBSOfferPropertylist []st.TVSBNCCBSOfferProperty
	var TVSBNCCBSOfferPropertyobj st.TVSBNCCBSOfferProperty
	var err error
	if TVSBNP.CCBSSubNo > 0 {
		tg.Getccbssubinfo(TVSBNP.CCBSSubNo)
	}
	cm.ExcutestoreDS("ICC", ` begin tvs_ccbsbn.Get_CCBSOFFER(:p_TRNSEQNO,:p_rs); end; `, TVSBNP.TRNSEQNO, sql.Out{Dest: &resultI})
	defer resultI.Close()
	values := make([]driver.Value, len(resultI.Columns()))
	colmap := createmapcol(resultI.Columns())
	for {
		err = resultI.Next(values)
		if err == nil {
			if err == io.EOF {
				break
			}
		} else {
			break
		}

		TVSBNCCBSOfferPropertyobj.Ccbsoffername = values[colmap["CCBS_OFFERNAME"]].(string)
		TVSBNCCBSOfferPropertyobj.Ccbssocid = values[colmap["CCBS_SOCID"]].(string)
		TVSBNCCBSOfferPropertyobj.Expirationdate = values[colmap["EXPIRATIONDATE"]].(time.Time)
		TVSBNCCBSOfferPropertyobj.TargetPayChannelID = values[colmap["TARGETPAYCHANNELID"]].(int64)
		if values[8].(int64) > 0 {
			TVSBNCCBSOfferPropertyobj.OverrideRCAmount = cm.StrTofloat64(values[colmap["OVERRIDE_RC_AMOUNT"]].(string))
			TVSBNCCBSOfferPropertyobj.OverrideRCDescription = values[colmap["OVERRIDE_RC_DESCRIPTION"]].(string)
		}
		TVSBNCCBSOfferPropertyobj.Processtype = values[colmap["PROCESSTYPE"]].(string)
		TVSBNCCBSOfferPropertyobj.EffectiveDateSpecified = values[colmap["EFFECTIVEDATESPECIFIED"]].(int64)
		if TVSBNCCBSOfferPropertyobj.EffectiveDateSpecified == 1 {
			TVSBNCCBSOfferPropertyobj.Effectivedate = values[colmap["EFFECTIVE_DATE"]].(time.Time)
		}

		TVSBNCCBSOfferPropertyobj.Newperiodind = values[colmap["NEW_PERIOD_IND"]].(string)
		TVSBNCCBSOfferPropertyobj.Action = values[colmap["ACTION"]].(string)
		TVSBNCCBSOfferPropertyobj.Extendedinfoname = values[colmap["EXTENDEDINFO_NAME"]].(string)
		TVSBNCCBSOfferPropertyobj.Extendedinfovalue = values[colmap["EXTENDEDINFO_VALUE"]].(string)
		TVSBNCCBSOfferPropertylist = append(TVSBNCCBSOfferPropertylist, TVSBNCCBSOfferPropertyobj)
	}

	return TVSBNCCBSOfferPropertylist
}

func changepackage(customerid int) string {
	var TVSBNP st.TVSBNProperty
	var omxRequest st.SubmitOrderOpRequest
	TVSBNP.CCBSORDERTYPEID = "128"
	TVSBNP.TVSCUSTOMERNO = customerid
	TVSBNP.CCBSFN = "CHANGEPACKAGE"
	TVSBNP.TRNSEQNO = preparejobdata(TVSBNP.TVSCUSTOMERNO, TVSBNP.CCBSFN)
	TVSBNP = getjobinfo(TVSBNP.TRNSEQNO)
	TVSBNP.TVSCUSTOMERNO = customerid
	omxRequest = mappingvalueomx(TVSBNP)
	print(omxRequest.Customer.CustomerId.Text)
	result, err := inboundtoomx(TVSBNP, omxRequest)
	if err == nil {
		print(result)
	}
	return TVSBNP.TRNSEQNO

}
func mappingvalueomx(TVSBNP st.TVSBNProperty) st.SubmitOrderOpRequest {
	var omxreq st.SubmitOrderOpRequest
	//
	mappingvalueomxorderinfo(TVSBNP, &omxreq)
	mappingvalueomxexistingsub(TVSBNP, &omxreq)
	mappingvalueomxoffer(TVSBNP, &omxreq)
	return omxreq
}
func mappingvalueomxorderinfo(TVSBNP st.TVSBNProperty, omxreq *st.SubmitOrderOpRequest) {
	omxreq.S = "http://services.omx.truecorp.co.th/SubmitOrder"
	omxreq.SE = "http://schemas.xmlsoap.org/soap/envelope/"
	omxreq.XSD = "http://www.w3.org/2001/XMLSchema"

	omxreq.Order.Channel.Text = TVSBNP.TVSBNOMXPropertyobj.Channel
	omxreq.Order.OrderId.Text = TVSBNP.CCBSorderno
	omxreq.Order.OrderType.Text = TVSBNP.OMXOrderType
	omxreq.Order.DealerCode.Text = TVSBNP.TVSBNOMXPropertyobj.DealerCode
	if TVSBNP.TVSBNOMXPropertyobj.EffectiveDateSpecified == 1 {
		omxreq.Order.EffectiveDate.Text = TVSBNP.TVSBNOMXPropertyobj.EffectiveDate //"2016-06-02T00:00:00"
	} else {
		omxreq.Order.EffectiveDate.Text = ""
	}
}
func mappingvalueomxexistingsub(TVSBNP st.TVSBNProperty, omxreq *st.SubmitOrderOpRequest) {

	omxreq.Customer.CustomerId.Text = TVSBNP.CCBSCustomerno
	omxreq.Customer.Account.AccountId = TVSBNP.CCBSAccountno
	omxreq.Customer.OU.OuId = TVSBNP.CCBSOUNo
	omxreq.Customer.OU.Subscriber.SubscriberId = TVSBNP.CCBSSubNo
	omxreq.Customer.OU.Subscriber.Status = "65"
	omxreq.Customer.OU.Subscriber.SubscriberNumber = cm.IntToStr(TVSBNP.TVSCUSTOMERNO)
	omxreq.Customer.OU.Subscriber.SubscriberType = "VC"
	omxreq.Customer.OU.Subscriber.PayChannelIdPrimary = TVSBNP.CCBSAccountno
	omxreq.Customer.OU.Subscriber.PayChannelIdSecondary = TVSBNP.CCBSAccountno
	omxreq.Customer.OU.Subscriber.ActivityInfo.ActivityReason = TVSBNP.CCBSACTIVITYREON
	omxreq.Customer.CustomerId.Text = TVSBNP.CCBSCustomerno
	//omxreq.Customer.OU.Subscriber.SubscriberAddress=nil

}
func mappingvalueomxoffer(TVSBNP st.TVSBNProperty, omxreq *st.SubmitOrderOpRequest) {
	var offer st.Omxccbsoffer
	//omxreq.Customer.OU.Subscriber.Offers := []Omxccbsoffer //[]Omxccbsoffer
	for i := 0; i < len(TVSBNP.TVSBNCCBSOfferPropertylist); i++ {
		offer.OfferName = TVSBNP.TVSBNCCBSOfferPropertylist[0].Ccbsoffername
		//append(omxreq.Customer.OU.Subscriber.Offers,offer)
	}

}
func inboundtoomx(TVSBNP st.TVSBNProperty, omxreq st.SubmitOrderOpRequest) (string, error) {
	url := TVSBNP.CCBSURLSERVICE
	client := &http.Client{}
	output, err := xml.MarshalIndent(omxreq, "  ", "    ")
	a := string(output)
	a = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sub="http://services.omx.truecorp.co.th/SubmitOrder">
	<soapenv:Header/>
	<soapenv:Body> ` + a + ` </soapenv:Body>
	</soapenv:Envelope> `
	//fmt.Println(a)
	//fmt.Println("********************************************************")
	requestContent := []byte(a)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("SOAPAction", `"/Services/SubmitOrderOp"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(errors.New("Error Respose " + resp.Status))
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	//m, _ := mxj.NewMapXml(contents, true)
	//fmt.Println(&m)
	return string(contents), nil
}
func getftid(ftid string) st.FinancialTransaction {
	var ftobj st.FinancialTransaction
	return ftobj
}
