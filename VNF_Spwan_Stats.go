package main

import(
	"fmt"
	"time"
	"os"
	"runtime/debug"
)

func statsPrint(){
   for true{
	var spawnNotInit, spawnInit, spawnInPro, spawnSucc, spawnFail int = 0,0,0,0,0
	time.Sleep(1000 * time.Millisecond)
	for index, _ := range vnfips{
		//fmt.Printf("\tVNF: [%v]>[%v]\n",index,len(VNFStates))
		if index >= len(VNFStates){
			spawnNotInit += 1
		}else if VNFStates[index].spawnStat == "Spawn Initiated"{
			spawnInit += 1
		}else if VNFStates[index].spawnStat == "Spawn In Progress"{
			spawnInPro += 1
		}else if VNFStates[index].spawnStat == "Spawn Success"{
			spawnSucc += 1
		}else if VNFStates[index].spawnStat == "Spawn Failed"{
			spawnFail += 1
		}
	}
   os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
   fmt.Printf("+-+-+-+\t\t+-+-+-+-+-+\t\t+-+-+-+\n\n")
   fmt.Printf("\tSpawns Not Initiated\t--> (%v/%v)\n",spawnNotInit,len(vnfips))
   fmt.Printf("\tSpawns Initiated\t--> (%v/%v)\n",spawnInit,len(vnfips))
   fmt.Printf("\tSpawns In Progress\t--> (%v/%v)\n",spawnInPro,len(vnfips))
   fmt.Printf("\tSpawns Success\t\t--> (%v/%v)\n",spawnSucc,len(vnfips))
   fmt.Printf("\tSpawns Failed\t\t--> (%v/%v)\n",spawnFail,len(vnfips))
   fmt.Printf("\n+-+-+-+\t\t+-+-+-+-+-+\t\t+-+-+-+\n")
   debug.FreeOSMemory()
}
}
