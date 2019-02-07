package main

import (
	"database/sql"
	"database/sql/driver"

	//"fmt"
	"io"
	"strconv"
	"time"
	cm "tvsglobal/common"
	st "tvsglobal/tvsstructs"

	//"github.com/jmoiron/sqlx"
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
		TVSBNP.TRNSEQNO = values[colmap["TRNSEQNO"]].(string)
		TVSBNP.CCBSorderno = values[colmap["CCBS_ORDERNO"]].(string)
		custno, err := strconv.Atoi(values[colmap["TVS_CUSTOMERNO"]].(string)) // strconv.ParseInt( values[2].(string), 10, 64)
		if err != nil {
			TVSBNP.TVSCUSTOMERNO = custno
		}

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

		TVSBNP.TVSBNCCBSOfferPropertylist = getccbsoffer(TVSBNP.TRNSEQNO)
	}

	return TVSBNP
}
func getccbsoffer(trnseqno string) []st.TVSBNCCBSOfferProperty {
	var resultI driver.Rows
	var TVSBNCCBSOfferPropertylist []st.TVSBNCCBSOfferProperty
	var TVSBNCCBSOfferPropertyobj st.TVSBNCCBSOfferProperty
	var err error
	cm.ExcutestoreDS("ICC", ` begin tvs_ccbsbn.Get_CCBSOFFER(:p_TRNSEQNO,:p_rs); end; `, trnseqno, sql.Out{Dest: &resultI})
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
	omxRequest = mappingvalueomx(TVSBNP)
	print(omxRequest.Text)
	return TVSBNP.TRNSEQNO
}
func mappingvalueomx(TVSBNP st.TVSBNProperty) st.SubmitOrderOpRequest {
	var omxreq st.SubmitOrderOpRequest
	//
	mappingvalueomxexistingsub(TVSBNP, &omxreq)
	mappingvalueomxoffer(TVSBNP, &omxreq)
	return omxreq
}
func mappingvalueomxexistingsub(TVSBNP st.TVSBNProperty, omxreq *st.SubmitOrderOpRequest) {

	omxreq.Customer.CustomerId.Text = TVSBNP.CCBSCustomerno
	omxreq.Customer.Account.AccountId = TVSBNP.CCBSAccountno
	omxreq.Customer.OU.OuId = TVSBNP.CCBSOUNo
	omxreq.Customer.OU.Subscriber.SubscriberId = TVSBNP.CCBSSubNo
	omxreq.Customer.CustomerId.Text = TVSBNP.CCBSCustomerno

}
func mappingvalueomxoffer(TVSBNP st.TVSBNProperty, omxreq *st.SubmitOrderOpRequest) {
	var offer st.Omxccbsoffer
	//omxreq.Customer.OU.Subscriber.Offers := []Omxccbsoffer //[]Omxccbsoffer
	for i := 0; i < len(TVSBNP.TVSBNCCBSOfferPropertylist); i++ {
		offer.OfferName = TVSBNP.TVSBNCCBSOfferPropertylist[0].Ccbsoffername
		//append(omxreq.Customer.OU.Subscriber.Offers,offer)
	}

}
func getftid(ftid string) st.FinancialTransaction {
	var ftobj st.FinancialTransaction
	return ftobj
}
