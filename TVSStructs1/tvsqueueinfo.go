package tvsstructs

type Tvsqueueinfo struct {
	Queuename string
	Address   string
}
type tvsorderinfo struct {
	OrderType string
	Otdername string
	tasklist  []taskinfo
}
type taskinfo struct {
	step     int
	Taskid   int
	Taskdesc string
}
