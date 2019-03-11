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

	tg "github.com/smsdevteam/tvsglobal/TVSCCGLOBAL"
	st "github.com/smsdevteam/tvsglobal/TVSStructs"
	cm "github.com/smsdevteam/tvsglobal/common"

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
func saveinboundomxlog(TVSBNP st.TVSBNProperty, requstd string, responsed string, startdate time.Time, stopdate time.Time) string {
	var trnseqno string
	cm.ExcutestoreDS("ICC", `  "Begin  tvs_ccbsbn.save_inboundlog(:p_tvs_custno,:p_trnseqno,:p_ccbs_fn,:p_ccbs_subfn,:p_url_service,
		:p_processflag,:p_error_code,:p_error_desc, 
	 :p_omxorderno,:p_omxtrackingid,:p_omx_repcode,:p_omx_repmsg,:p_requestdata,:p_responsedata,
	 :p_omxurl,:p_startdare,:p_stopdate,:p_ExternalRefno); end; 
	`, TVSBNP.TVSCUSTOMERNO, TVSBNP.TRNSEQNO, TVSBNP.CCBSFN, TVSBNP.CCBSSUBFN, TVSBNP.CCBSarURLSERVICE,
		1, " ", "  ",
		TVSBNP.CCBSarURLSERVICE, startdate, stopdate, " ")
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
	if TVSBNP.CCBSSubNo != "0" {
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
		//TVSBNCCBSOfferPropertyobj.Ecbsservicetype
		TVSBNCCBSOfferPropertyobj.Ccbsoffername = values[colmap["CCBS_OFFERNAME"]].(string)
		TVSBNCCBSOfferPropertyobj.Ccbsservicetype = values[colmap["CCBS_SERVICETYPE"]].(string)

		TVSBNCCBSOfferPropertyobj.Ccbssocid = values[colmap["CCBS_SOCID"]].(string)
		TVSBNCCBSOfferPropertyobj.Expirationdate = cm.TimeToStr(values[colmap["EXPIRATIONDATE"]].(time.Time))
		TVSBNCCBSOfferPropertyobj.TargetPayChannelID = values[colmap["TARGETPAYCHANNELID"]].(int64)
		if values[8].(int64) > 0 {
			TVSBNCCBSOfferPropertyobj.OverrideRCAmount = cm.StrTofloat64(values[colmap["OVERRIDE_RC_AMOUNT"]].(string))
			TVSBNCCBSOfferPropertyobj.OverrideRCDescription = values[colmap["OVERRIDE_RC_DESCRIPTION"]].(string)
		}
		TVSBNCCBSOfferPropertyobj.Processtype = values[colmap["PROCESSTYPE"]].(string)
		TVSBNCCBSOfferPropertyobj.EffectiveDateSpecified = values[colmap["EFFECTIVEDATESPECIFIED"]].(int64)
		if TVSBNCCBSOfferPropertyobj.EffectiveDateSpecified == 1 {
			TVSBNCCBSOfferPropertyobj.Effectivedate = cm.TimeToStr(values[colmap["EFFECTIVE_DATE"]].(time.Time))
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
func changepackagep(trackingno string, customerid int) st.ResponseResult {
	var TVSBNP st.TVSBNProperty
	var omxRequest st.SubmitOrderOpRequest
	var tags []string
	var applog cm.Applog
	var res st.ResponseResult
	tags = append(tags, "env7")
	tags = append(tags, "TVSNote")
	tags = append(tags, "applogs")
	applog   = cm.NewApploginfo(trackingno,"TVSBN","changepackagep" ,tags)
	applog.Request=cm.IntToStr(customerid)
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
		res.ErrorCode= 0 
		res.ErrorDesc=" "
	}else{
		res.ErrorCode= 999
		res.ErrorDesc=err.Error()
	}
	return res

}
func suspendsub(customerid int) string {
	var TVSBNP st.TVSBNProperty
	var omxRequest st.SubmitOrderOpRequest
	TVSBNP.CCBSORDERTYPEID = "124"
	TVSBNP.TVSCUSTOMERNO = customerid
	TVSBNP.CCBSFN = "SUSPENDSUB"
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
func restoresub(customerid int) string {
	var TVSBNP st.TVSBNProperty
	var omxRequest st.SubmitOrderOpRequest
	TVSBNP.CCBSORDERTYPEID = "125"
	TVSBNP.TVSCUSTOMERNO = customerid
	TVSBNP.CCBSFN = "RESTORESUB"
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
	omxreq.Customer.OU.Subscriber.ResourceInfo.ResourceCategory = "R"
	omxreq.Customer.OU.Subscriber.ResourceInfo.ResourceCategory = "R"
	omxreq.Customer.OU.Subscriber.ResourceInfo.ResourceName = "Television ID"
	omxreq.Customer.OU.Subscriber.ResourceInfo.ValuesArray = cm.IntToStr(TVSBNP.TVSCUSTOMERNO)
	//omxreq.Customer.OU.Subscriber.SubscriberAddress=nil
}
func mappingvalueomxoffer(TVSBNP st.TVSBNProperty, omxreq *st.SubmitOrderOpRequest) {
	var offer st.Omxccbsoffer
	var offerpara st.Omxccbsofferpara
	//omxreq.Customer.OU.Subscriber.Offers := []Omxccbsoffer //[]Omxccbsoffer
	for i := 0; i < len(TVSBNP.TVSBNCCBSOfferPropertylist); i++ {
		offer.OfferName = TVSBNP.TVSBNCCBSOfferPropertylist[i].Ccbsoffername
		offer.Action = TVSBNP.TVSBNCCBSOfferPropertylist[i].Action
		//offer.EffectiveDate = TVSBNP.TVSBNCCBSOfferPropertylist[0].Effectivedate
		//	offer.ExpirationDate = TVSBNP.TVSBNCCBSOfferPropertylist[0].Effectivedate //"2019-02-11T00:00:01"
		offer.ServiceType = TVSBNP.TVSBNCCBSOfferPropertylist[i].Ccbsservicetype
		//offer.TargetPayChannelId = nil
		if TVSBNP.TVSBNCCBSOfferPropertylist[i].Ccbssocid != "0" {
			offer.OfferInstanceId = TVSBNP.TVSBNCCBSOfferPropertylist[i].Ccbssocid
		}
		if TVSBNP.TVSBNCCBSOfferPropertylist[i].EffectiveDateSpecified == 1 {
			offer.EffectiveDate = TVSBNP.TVSBNCCBSOfferPropertylist[i].Effectivedate
		}
		if TVSBNP.TVSBNCCBSOfferPropertylist[i].Action == "REMOVE" {
			offer.ExpirationDate = TVSBNP.TVSBNCCBSOfferPropertylist[i].Effectivedate
		}

		//check new period
		if TVSBNP.TVSBNCCBSOfferPropertylist[i].Newperiodind != " " {
			offerpara.ParamName = "New period ind"
			offerpara.ValuesArray = TVSBNP.TVSBNCCBSOfferPropertylist[i].Newperiodind
			offer.Offerparas = append(offer.Offerparas, offerpara)
		}
		if TVSBNP.TVSBNCCBSOfferPropertylist[i].OverrideRCAmount != 0 {
			offerpara.ParamName = "Override RC description Thai"
			offerpara.ValuesArray = TVSBNP.TVSBNCCBSOfferPropertylist[i].OverrideRCDescription
			offer.Offerparas = append(offer.Offerparas, offerpara)

			offerpara.ParamName = "Override RC description Eng"
			offerpara.ValuesArray = TVSBNP.TVSBNCCBSOfferPropertylist[i].OverrideRCDescriptionEng
			offer.Offerparas = append(offer.Offerparas, offerpara)

			offerpara.ParamName = "Override RC amount"
			//offerpara.ValuesArray =cm.Int64ToStr( TVSBNP.TVSBNCCBSOfferPropertylist[0].OverrideRCAmount)
			offer.Offerparas = append(offer.Offerparas, offerpara)
		}
		if TVSBNP.TVSBNCCBSOfferPropertylist[0].OverrideOCAmount != 0 {
			offerpara.ParamName = "Override OC description Thai"
			offerpara.ValuesArray = TVSBNP.TVSBNCCBSOfferPropertylist[0].OverrideOCDescription
			offer.Offerparas = append(offer.Offerparas, offerpara)

			offerpara.ParamName = "Override OC description Eng"
			offerpara.ValuesArray = TVSBNP.TVSBNCCBSOfferPropertylist[0].OverrideOCDescriptionEng
			offer.Offerparas = append(offer.Offerparas, offerpara)

			offerpara.ParamName = "Override OC amount"
			//offerpara.ValuesArray = TVSBNP.TVSBNCCBSOfferPropertylist[0].OverrideOCAmount
			offer.Offerparas = append(offer.Offerparas, offerpara)
		}
		if TVSBNP.TVSBNCCBSOfferPropertylist[0].Extendedinfoname == "RC_END_DATE" {
			offerpara.ParamName = TVSBNP.TVSBNCCBSOfferPropertylist[0].Extendedinfoname
			offerpara.ValuesArray = TVSBNP.TVSBNCCBSOfferPropertylist[0].Extendedinfovalue
			if TVSBNP.TVSBNCCBSOfferPropertylist[0].Extendedinfoname == "RC_END_DATE" {
				offerpara.EffectiveDate = TVSBNP.TVSBNCCBSOfferPropertylist[0].Effectivedate
			}
			offer.Offerparas = append(offer.Offerparas, offerpara)
		}
		// check add offer
		omxreq.Customer.OU.Subscriber.Offers = append(omxreq.Customer.OU.Subscriber.Offers, offer)
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
	fmt.Println("********************************************************")

	fmt.Println(a)
	fmt.Println("********************************************************")
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
