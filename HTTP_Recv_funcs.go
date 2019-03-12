package main

import(
	"net/http"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"regexp"
//	"fmt"
	"strconv"
	"runtime/debug"
//	"io"
)

func Platform(w http.ResponseWriter, r *http.Request){
        w.WriteHeader(http.StatusCreated)
        defer r.Body.Close()
	defer debug.FreeOSMemory()
	host := r.Host
	index := indexOf(strings.Split(host,":")[0],vnfips)

	loopCount := 1
        for loopCount<5{
                if VNFStates[index].NewEvent == "Recv"{ break }
                time.Sleep(1 * time.Second)
		loopCount += 1
        }

        data, _ := ioutil.ReadAll(r.Body)
	log.Printf("\n *DBG* %v URL: %v", string(host),r.URL.Path)
        log.Printf(" *DBG* %v Body: %v",string(host),string(data))
        if strings.Contains(string(data), "\"param_name\":\"health_monitoring\""){
		VNFStates[index].NewEvent = "Send"
		VNFStates[index].pos = 5
		re := regexp.MustCompile("{\"subscription.*lifetime\":(.*),\"periodicity\":(.*),\"param_name\":*")
		VNFStates[index].MonTime, _ = strconv.Atoi(re.FindStringSubmatch(string(data))[1])
		VNFStates[index].MonPeriod, _ = strconv.Atoi(re.FindStringSubmatch(string(data))[2])
		log.Printf(" *INFO* %v MeasurementReport_Health lifetime:%v, periodicity:%v",string(host),VNFStates[index].MonTime, VNFStates[index].MonPeriod)
		log.Printf(" *INFO* %v index of MeasurementReport_Health is: %v", string(host),VNFStates[index].pos)
        //MeasurementReport_health
        }else if strings.Contains(string(data), "\"param_name\":\"cpu_usage\""){
		VNFStates[index].NewEvent = "Send"
		VNFStates[index].pos = 6
		re := regexp.MustCompile("{\"subscription.*lifetime\":(.*),\"periodicity\":(.*),\"param_name\":*")
                VNFStates[index].MonTime, _ = strconv.Atoi(re.FindStringSubmatch(string(data))[1])
                VNFStates[index].MonPeriod, _ = strconv.Atoi(re.FindStringSubmatch(string(data))[2])
                log.Printf(" *INFO* %v MeasurementReport_CPU lifetime:%v, periodicity:%v",string(host),VNFStates[index].MonTime, VNFStates[index].MonPeriod)
		log.Printf(" *INFO* %v index of MeasurementReport_CPU is: %v", string(host),VNFStates[index].pos)
		wg_http.Done()
        //MeasurementReport_CPU
        }else if strings.Contains(string(data), "netconf"){
		VNFStates[index].NewEvent = "Send"
		VNFStates[index].pos = 1
		log.Printf(" *INFO* %v index of NotConfigure is: %v", string(host),VNFStates[index].pos)

		confre := regexp.MustCompile("{\"vnf\":{\"vimType\".*netconfPort\":\"(.*)\",\"netconfIPAddress\":\"(.*)\",\"vimId\"*")
                netconfPort := confre.FindStringSubmatch(string(data))[1]
                netconfIPAddress := confre.FindStringSubmatch(string(data))[2]
		log.Printf(" *INFO* %v netConf Addr: [%v:%v]",string(host),netconfIPAddress,netconfPort)

		//wg_Listen.Add(1)
		go Listen(netconfIPAddress+":"+netconfPort)
///////////////////////////////////////
        //NotConfigure
        }else{VNFStates[index].NewEvent = "Recv"}
}

func Operations_start(w http.ResponseWriter, r *http.Request){
	debug.FreeOSMemory()
        w.WriteHeader(http.StatusOK)
	host := r.Host
	index := indexOf(strings.Split(host,":")[0],vnfips)

	loopCount := 1
        for loopCount<5{
                if VNFStates[index].NewEvent == "Recv"{ break }
                time.Sleep(1 * time.Second)
                loopCount += 1
        }

        data, _ := ioutil.ReadAll(r.Body)
	log.Printf("\n *DBG* %v URL: %v", string(host),r.URL.Path)
        log.Printf("\n *DBG* %v Body: %v",string(host),string(data))

        w.Header().Set("content-type", "application/vnd.yang.operation+json")
        w.Write([]byte("{\"output\": {\"result\": \"success\"}}"))

	VNFStates[index].NewEvent = "Send"
	VNFStates[index].pos = 4
        log.Printf(" *INFO* %v index of StartSuccess is: %v", string(host),VNFStates[index].pos)
	VNFStates[index].spawnStat = "Spawn Success"
        //StartSuccess
}
