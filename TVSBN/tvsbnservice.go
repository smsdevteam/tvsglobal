package main

import (
	"database/sql"
	"database/sql/driver"

	"fmt"
	"io"
	"strconv"
	cm "tvsglobal/common"
	st "tvsglobal/tvsstructs"

	_ "gopkg.in/goracle.v2"
)

func preparejobdata(customerid int, ccbsfn string) string {
	var trnseqno string
	cm.ExcutestoreDS("ICC", ` begin :result := tvs_goccbsbn.create_omxjob(:p_custno,
	 		 :p_ccbsfn,:p_shno,:p_sheventno,:p_processflag,
	 		 :p_allowsendomx,:p_ccbs_activityreason); end;`, sql.Out{Dest: &trnseqno}, customerid, ccbsfn, 0, 0, 1, 1, "CREQ")
	return trnseqno
}
func getjobinfo(trnseqno string) st.TVSBNProperty {
	var TVSBNP st.TVSBNProperty
	var resultI driver.Rows
	var err error
	cm.ExcutestoreDS("ICC", ` begin  tvs_goccbsbn.getjobinfo(:p_transeq,
		:p_rs); end;`, trnseqno, sql.Out{Dest: &resultI})
	defer resultI.Close()
	values := make([]driver.Value, len(resultI.Columns()))
	for {
		err = resultI.Next(values)
		if err == nil {
			if err == io.EOF {
				break
			}
		} else {
			break
		}

		TVSBNP.TRNSEQNO = values[0].(string)
		TVSBNP.CCBSorderno = values[1].(string)
		custno, err := strconv.Atoi(values[2].(string)) // strconv.ParseInt( values[2].(string), 10, 64)
		if err != nil {
			TVSBNP.TVSCUSTOMERNO = custno
		}

		TVSBNP.TVSAccountno = values[3].(string)
		TVSBNP.SHNO = values[4].(string)
		TVSBNP.CCBSFN = values[5].(string)
		TVSBNP.RECALLOFFER = values[6].(string)
		TVSBNP.HAVEOCCHARE = values[7].(string)
		TVSBNP.OMXOrderType = values[8].(string)
		TVSBNP.TVSCustomerType = values[9].(string)
		TVSBNP.CCBSACTIVITYREON = values[10].(string)
		TVSBNP.CCBSUSERTEXT = values[11].(string)
		TVSBNP.CCBSCustomerno = values[12].(string)
		TVSBNP.CCBSAccountno = values[13].(string)
		TVSBNP.CCBSOUNo = values[14].(string)
		TVSBNP.CCBSSubNo = values[15].(string)
		TVSBNP.CCBSURLSERVICE = values[16].(string)
		TVSBNP.Refvalue = values[17].(string)
		TVSBNP.Reftype = values[18].(string)
		TVSBNP.ToCCBSCustomerno = values[19].(string)
		TVSBNP.OLDCCBSACCOUNTNO = values[20].(string)
		TVSBNP.Oldccbssubno = values[21].(string)
		TVSBNP.AddSOCLevelOU = values[22].(string)
		fmt.Println(TVSBNP)
		TVSBNP.TVSBNCCBSOfferPropertylist= getccbsoffer(TVSBNP.TRNSEQNO)
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
	for {
		err = resultI.Next(values)
		if err == nil {
			if err == io.EOF {
				break
			}
		} else {
			break
		}
	 	TVSBNCCBSOfferPropertyobj.ccbs_offername = values[0].(string)
		TVSBNCCBSOfferPropertyobj.offerInstanceId = values[1].(string)
		TVSBNCCBSOfferPropertyobj.processtype =values[2].(string)
		TVSBNCCBSOfferPropertyobj.action = values[30].(string)
		TVSBNCCBSOfferPropertyobj.effectiveDateSpecified = values[4].(string)
		If TVSBNCCBSOfferPropertyobj.effectiveDateSpecified = 1 {
		TVSBNCCBSOfferPropertyobj.effectiveDateSpecified = values[5].(string)
			TVSBNCCBSOfferPropertyobj.effective_date =  values[6].(string)
		}

		//If IsDate(ds.Tables(0).Rows(i)("expirationdate")) Then
		TVSBNCCBSOfferPropertyobj.expirationdate =  values[0].(string)
		//End If

		TVSBNCCBSOfferPropertyobj.targetPayChannelId =  values[7].(string)
		TVSBNCCBSOfferPropertyobj.targetPayChannelId =  values[8].(string)
		TVSBNCCBSOfferPropertyobj.Override_RC_Amount =  values[9].(string)
		TVSBNCCBSOfferPropertyobj.Override_RC_Description = values[10].(string)
		TVSBNCCBSOfferPropertyobj.New_period_ind =  values[11].(string)
		TVSBNCCBSOfferPropertyobj.extendedinfo_name =  values[12].(string)
		TVSBNCCBSOfferPropertyobj.extendedinfo_value =  values[13].(string)
		TVSBNCCBSOfferPropertyobj.ccbs_servicetype =  values[14].(string)
	}

	return nil
}
func changepackage(customerid int) string {
	var TVSBNP st.TVSBNProperty
	TVSBNP.CCBSORDERTYPEID = "128"
	TVSBNP.TVSCUSTOMERNO = customerid
	TVSBNP.CCBSFN = "CHANGEPACKAGE"
	TVSBNP.TRNSEQNO = preparejobdata(TVSBNP.TVSCUSTOMERNO, TVSBNP.CCBSFN)
	TVSBNP = getjobinfo(TVSBNP.TRNSEQNO)
	return TVSBNP.TRNSEQNO
}
func getftid(ftid string) st.FinancialTransaction {
	var ftobj st.FinancialTransaction
	return ftobj
}
