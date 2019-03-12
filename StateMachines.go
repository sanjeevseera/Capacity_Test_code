package main

import(
	//"log"
	//"fmt"
//	"github.com/qor/transition"
)

type states struct {
  NewEvent string
  MonTime int
  MonPeriod int
  DBellRetry int
  pos int
  State []string// {"DoorBell", "NotConfigure", "InitSuccess", "ConfigureSuccess", "StartSuccess", "MeasurementReport_Health", "MeasurementReport_CPU"}
  spawnStat string
}
//var VNFStates = []*states{}
func (main *states) SetState(st string){
	main.State = append(main.State, st)
}
func (main *states) GetState() string {
        return main.State[main.pos]
}
func (main *states) GetStates() []string {
	return main.State
}

func InitState(p *states){
	p.spawnStat = "Spawn Not Initiated"
	p.NewEvent = "Initial"
	p.pos = 0
	p.DBellRetry = DBLretry
	p.SetState("DoorBell")
	p.SetState("NotConfigure")
	p.SetState("InitSuccess")
	p.SetState("ConfigureSuccess")
	p.SetState("StartSuccess")
	p.SetState("MeasurementReport_Health")
	p.SetState("MeasurementReport_CPU")
}
