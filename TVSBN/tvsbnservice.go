package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
	"time"
	cm "tvsglobal/common"
	st "tvsglobal/tvsstructs"
	
	_ "gopkg.in/goracle.v2"
)
func getomxorderinfo(TVSBNP st.TVSBNProperty){
	var resultI driver.Rows
	cm.ExcutestoreDS("ICC", ` Begin  tvs_ccbsbn.get_orderinfo(:p_trnseqno,
	:p_suhistorno,:p_ccbsfnname,:p_ccbbssubfn,:p_rs); 
	end;`,TVSBNP.TRNSEQNO,0,TVSBNP.CCBSFN,TVSBNP.CCBSSUBFN, sql.Out{Dest: &resultI})

	
}
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
	for {
		err = resultI.Next(values)
		if err == nil {
			if err == io.EOF {
				break
			}
		} else {
			break
		}
		
		TVSBNCCBSOfferPropertyobj.Ccbsoffername = values[3].(string)
		TVSBNCCBSOfferPropertyobj.Ccbssocid = values[4].(string)
		TVSBNCCBSOfferPropertyobj.Expirationdate = values[5].(time.Time)
		TVSBNCCBSOfferPropertyobj.TargetPayChannelID = values[7].(int64) 
		if values[8].(int64)>0{
		TVSBNCCBSOfferPropertyobj.OverrideRCAmount = cm.StrTofloat64(values[8].(string))
		TVSBNCCBSOfferPropertyobj.OverrideRCDescription = values[9].(string)
		}
		TVSBNCCBSOfferPropertyobj.Processtype = values[10].(string)
		TVSBNCCBSOfferPropertyobj.EffectiveDateSpecified = values[11].(int64) 
		if TVSBNCCBSOfferPropertyobj.EffectiveDateSpecified == 1 {
		 	TVSBNCCBSOfferPropertyobj.Effectivedate = values[12].(time.Time)
		}
		//If IsDate(ds.Tables(0).Rows(i)("expirationdate")) Then
		
		//End If

		//TVSBNCCBSOfferPropertyobj.TargetPayChannelId =  values[7].(string)
		//TVSBNCCBSOfferPropertyobj.TargetPayChannelId =  values[8].(string)
		
		TVSBNCCBSOfferPropertyobj.Newperiodind = values[14].(string)
		TVSBNCCBSOfferPropertyobj.Action = values[15].(string)
		TVSBNCCBSOfferPropertyobj.Extendedinfoname = values[16].(string)
		TVSBNCCBSOfferPropertyobj.Extendedinfovalue = values[17].(string)
		TVSBNCCBSOfferPropertylist = append(TVSBNCCBSOfferPropertylist, TVSBNCCBSOfferPropertyobj)
		//TVSBNCCBSOfferPropertyobj.Ccbsservicetype =  values[14].(string)
	}

	return TVSBNCCBSOfferPropertylist
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
