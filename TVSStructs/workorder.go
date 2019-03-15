package tvsstructs

type WorkorderInfo struct {
	ID          int64
	ProblemDesc string
	WorkorderServiceDTlist []WorkorderServiceDTInfo  
}
 type WorkorderServiceDTInfo struct {
	ServiceId  int64
	ServiceDescription string

 }