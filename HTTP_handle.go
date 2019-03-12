package main

import(
	"log"
	"strings"
	"net/http"
	//"fmt"
	"time"
	"sync"
	"runtime/debug"
)
var httpcount int = 0
var VNFStates = []*states{}
var wg_http sync.WaitGroup
//var wg_Listen sync.WaitGroup
//var wg_listen sync.WaitGroup

func NewVnf(caddr string, uuid string, body string) {
	var OrderStateMachine = new(states)
	var index int
	var running bool = true
	//var NewEvent string = "Initial"
	log.Println(" *INFO* "+"VNF Address is:",caddr)
	log.Println(" *INFO* "+"VNFM Address is:",raddr)
	log.Println(" *INFO* "+"VNF Address IP:",strings.Split(caddr,":")[0],"PORT", strings.Split(caddr,":")[1])

// HTTP SERVER *START*
	if httpcount == 0{
	http.HandleFunc("/api/running/sys/platform", Platform)
	http.HandleFunc("/api/running/sysops/lcm/_operations/start", Operations_start)
	httpcount = 1
	}
	wg_http.Add(1)
	go http.ListenAndServe(caddr, nil)
	log.Printf(" *INFO* %v HTTP server connection Created: ", caddr)


// HTTP SERVER *END*

	InitState(OrderStateMachine)
	VNFStates=append(VNFStates,OrderStateMachine)
	index = indexOf(strings.Split(caddr,":")[0],vnfips)

    for running{
	//log.Printf(" *INFO* vnfips: %v", vnfips)
	//log.Printf(" *INFO* VNFStates: %v", VNFStates)
	log.Printf(" *INFO* %v vnf Current State index: %v", caddr, index)
        log.Printf(" *INFO* %v vnf Current State is: %v", caddr, VNFStates[index].GetState())
        switch VNFStates[index].NewEvent{
        case "Initial":
		VNFStates[index].spawnStat = POST_DoorBell(caddr, body)
		if VNFStates[index].spawnStat == "Spawn Failed"{
			running = false
		}else{
			VNFStates[index].NewEvent = "Recv"
		}
        case "Send":
                if VNFStates[index].GetState() == "NotConfigure"{
			VNFStates[index].spawnStat = POST_VNFNotification(caddr, string(VNFStates[index].GetState()), uuid)
			if VNFStates[index].spawnStat == "Spawn Failed"{
				running = false
			} else{
                        VNFStates[index].NewEvent = "Send"
			VNFStates[index].pos += 1
			}
                }else if VNFStates[index].GetState() == "InitSuccess"{
			VNFStates[index].spawnStat = POST_VNFNotification(caddr, string(VNFStates[index].GetState()), uuid)
                        VNFStates[index].NewEvent = "Recv"
                        if VNFStates[index].GetState() == "ConfigureSuccess"{
                            VNFStates[index].NewEvent = "Send"
                        }
//			VNFStates[index].NewEvent = "Send"
//			VNFStates[index].pos += 1
                }else if VNFStates[index].GetState() == "ConfigureSuccess"{
			time.Sleep(2 * time.Second)
			VNFStates[index].spawnStat = POST_VNFNotification(caddr, string(VNFStates[index].GetState()), uuid)
                        VNFStates[index].NewEvent = "Recv"
			//VNFStates[index].pos += 1
                }else if VNFStates[index].GetState() == "StartSuccess"{
			VNFStates[index].spawnStat = POST_VNFNotification(caddr, string(VNFStates[index].GetState()), uuid)
                        VNFStates[index].NewEvent = "Recv"
			//VNFStates[index].pos += 1
                }else if VNFStates[index].GetState() == "MeasurementReport_Health"{
			go POST_monitorVnf(caddr, "HEALTH", VNFStates[index].MonTime,VNFStates[index].MonPeriod)
                        VNFStates[index].NewEvent = "Recv"
                        //VNFStates[index].pos += 1
                }else if VNFStates[index].GetState() == "MeasurementReport_CPU"{
                        go POST_monitorVnf(caddr, "CPU",VNFStates[index].MonTime,VNFStates[index].MonPeriod)
			//VNFStates[index].NewEvent = "Recv"
			VNFStates[index].NewEvent = "Exit"
			//log.Printf(" *MSG * VNF: %v , All Messages SEND/RECV DONE", vnfips[index])
			//running = false
                }
        case "Recv":
		loopCount := 1
		for loopCount <= 120{
			if VNFStates[index].NewEvent == "Send"{ break }
			time.Sleep(1 * time.Second)
			loopCount += 1
		}
		if loopCount > 120 && VNFStates[index].GetState() == "DoorBell"{
			if VNFStates[index].DBellRetry != 0{
				VNFStates[index].DBellRetry = VNFStates[index].DBellRetry - 1
				VNFStates[index].NewEvent = "Initial"
			}else{
				log.Printf(" *ERR * %v Probelem in receiveing HTTP message, STOP sending other HTTP messages", caddr)
				VNFStates[index].spawnStat = "Spawn Failed"
				running = false
			}
		}else if loopCount > 120 && "MeasurementReport" != strings.Split(VNFStates[index].GetState(),"_")[0]{
			log.Printf(" *ERR * %v Probelem in receiveing HTTP message, STOP sending other HTTP messages", caddr)
			VNFStates[index].spawnStat = "Spawn Failed"
			running = false
		}else{
			log.Printf(" *INFO* %v changed Event form Recv to Send", caddr)
		}
	case "Exit":
		running = false
        default:
		log.Printf(" *WRG * %v Unexpected Exit", caddr)
		running = false
        }
}
wg_http.Wait()
//wg_Listen.Wait()
//wg_listen.Wait()
debug.FreeOSMemory()
wg_vnf.Done()
log.Printf(" *INFO * %v received Wait groups", caddr)
} //NewVnf END
