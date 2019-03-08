package main

import (
	"database/sql"
	"database/sql/driver"

	//"encoding/json"
	//"log"
	cm "tvsglobal/common"
	st "tvsglobal/tvsstructs"

	//"github.com/streadway/amqp"
	_ "gopkg.in/goracle.v2"
)

func generatetasklist(Trackingno string, TVSOrdReq  st.TVSSubmitOrderToQueue) (string, st.TVSQueueSubmitOrderReponse) {

	var TVSOrdRes st.TVSQueueSubmitOrderReponse
	var resultI driver.Rows
	var err error
	cm.ExcutestoreDS("ICC", `begin tvs_servorder.generatetasklist(:p_trackingno,:p_ordertype,:p_rs );  end;`,
		Trackingno, TVSOrdReq.TVSOrdReq.OrderType, sql.Out{Dest: &resultI})
	 	defer resultI.Close()
		values := make([]driver.Value, len(resultI.Columns()))
		colmap := cm.Createmapcol(resultI.Columns())
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
		return "", TVSOrdRes
}
