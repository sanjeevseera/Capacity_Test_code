package main

import(
        "net/http"
        "io/ioutil"
        "log"
        "strings"
	"time"
	"net"
	"runtime/debug"
//        "io"
	//"fmt"
)

//var POST_DoorBell_body = string(`{"cloudProfileId": "Openstack-10.10.210.201-vnfm-vnfm","connectionPoints": [ {"cpInterfaceIpAddress": "XX.XX.XX.XX","interfaceName": "eth0"} ],"nsdId": "SDNS_NSD","nsdVersion": "v71537","uuid": "5158b400-0195-4f62-853b-952b1d52a188","vduId": "SDNS_SITE1_VDU", "vnfdId": "SDNS_SITE1_VNFD","vnfdVersion": "v71537","vnfType": "EAST-VNF","vnfVersion": "EAST-v35"}`)

var POST_VNFNotification_body = string(`{ "systemId":"uuuuiidd", "messageName":"vnfStateChange","event":"AAAAAAAAA"}`)

var POST_monitorHEALTH_body = string(`{ "measurementReport" :{ "num_reports":1,"Timestamp":"2018-08-16T10:23:19","monitoring_reports" :  [ {"param-name" : "health_monitoring","interface":"ve_vnfm_c","type":"guage_c","value_datatype":"STRING","value":0,"unit":"", "Timestamp":"2018-08-16T10:23:19"}]}}`)

var POST_monitorCPU_body = string(`{ "measurementReport" :{ "num_reports":1,"Timestamp":"2017-08-04T16:33:19","monitoring_reports" :  [ {"param-name" : "cpu_usage","interface":"ve_vnfm_c","type":"guage_c","value_datatype":"FLOAT","value":4,"unit":"PERCENTAGE", "Timestamp":"2017-08-04T16:33:19"}]}}`)

func POST_DoorBell(caddr string, POST_DoorBell_body string) (string){
///
	defer debug.FreeOSMemory()
	localAddr, err := net.ResolveIPAddr("ip", strings.Split(caddr,":")[0])
	if err != nil {
		log.Printf(" *ERR* %v error: %v\n", caddr,err)
		return "Spawn Failed"
	}

	localTCPAddr := net.TCPAddr{
	    IP: localAddr.IP,
	}

	cli := &http.Client{
	    Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	        DialContext: (&net.Dialer{
		    LocalAddr: &localTCPAddr,
	            Timeout:   30 * time.Second,
		    KeepAlive: 30 * time.Second,
	            DualStack: true,
		}).DialContext,
	        MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
	        TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	    },
	}
///
        buf := strings.Replace(POST_DoorBell_body, "XX.XX.XX.XX", strings.Split(caddr,":")[0], -1)
        req, err := http.NewRequest("POST", "http://" + raddr + "/vnflcm/api/v1.0/doorbell", strings.NewReader(string(buf)))
        if err != nil {
                log.Printf(" *ERR* %v http.NewRequest() error: %v\n", caddr,err)
                return "Spawn Failed"
        }
        req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "close")
        //cli := &http.Client{}
	resp, err := cli.Do(req)


        log.Printf("\n *DBG* %v Request: POST http://%v/vnflcm/api/v1.0/doorbell", caddr,raddr)
        log.Printf(" *DBG* %v Body: %v", caddr,string(buf))

        if err != nil {
                log.Printf(" *ERR* %v http.Do() error: %v\n", caddr,err)
                return "Spawn Failed"
        }
        defer resp.Body.Close()

        data, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Printf(" *ERR* %v ioutil.ReadAll() error: %v\n", caddr,err)
                return "Spawn Failed"
        }
        Status := resp.Status
	log.Printf(" *DBG* %v Status: %v", caddr, string(Status))
	log.Printf(" *DBG* %v Body: %v", caddr, string(data))
	if string(Status)[0:3] == "200"{
		return "Spawn Initiated"
	}else{
		log.Printf(" *DBG* %v StatusFailed: %v", caddr, string(Status))
		return "Spawn Failed"
	}
}

func POST_VNFNotification(caddr string, event string, uuid string)(string){
///
	defer debug.FreeOSMemory()
        localAddr, err := net.ResolveIPAddr("ip", strings.Split(caddr,":")[0])
        if err != nil {
		log.Printf(" *ERR* %v error: %v\n", caddr,err)
		return "Spawn Failed"
        }

        localTCPAddr := net.TCPAddr{
            IP: localAddr.IP,
        }

        cli := &http.Client{
            Transport: &http.Transport{
                Proxy: http.ProxyFromEnvironment,
                DialContext: (&net.Dialer{
                    LocalAddr: &localTCPAddr,
                    Timeout:   30 * time.Second,
                    KeepAlive: 30 * time.Second,
                    DualStack: true,    
                }).DialContext,
                MaxIdleConns:          100,
                IdleConnTimeout:       90 * time.Second,
                TLSHandshakeTimeout:   10 * time.Second,        
                ExpectContinueTimeout: 1 * time.Second,
            },
        }
///
        buf := strings.Replace(POST_VNFNotification_body, "AAAAAAAAA", event, -1)
	buf = strings.Replace(buf, "uuuuiidd", uuid, -1)
        req, err := http.NewRequest("POST", "http://" + raddr + "/vnflcm/api/v1.0/vnfNotification", strings.NewReader(string(buf)))
        if err != nil {
                log.Printf(" *ERR* %v http.NewRequest() error: %v\n", caddr,err)
                return "Spawn Failed"
        }
        req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "close")
        //cli := &http.Client{}
        resp, err := cli.Do(req)
        log.Printf("\n *DBG* %v Request: POST http://%v/vnflcm/api/v1.0/vnfNotification",caddr,raddr)
        log.Printf(" *DBG* %v Body: %v", caddr,string(buf))

        if err != nil {
                log.Printf(" *ERR* %v http.Do() error: %v\n", caddr,err)
                return "Spawn Failed"
        }
        defer resp.Body.Close()

        data, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Printf(" *ERR* %v ioutil.ReadAll() error: %v\n", caddr,err)
                return "Spawn Failed"
        }
        Status := resp.Status
        log.Printf(" *DBG* %v Status: %v", caddr, string(Status))
        log.Printf(" *DBG* %v Body: %v", caddr, string(data))
	if event == "StartSuccess"{return "Spawn Success"}
	return "Spawn In Progress"
}

func POST_monitorVnf(caddr string, param string, MonTime int, MonPeriod int){
///
    defer debug.FreeOSMemory()
    Max_time := MonTime
    now := time.Now()
    var buf,timestamp string
    var data []byte
    //var data int64
    var err error
    var localAddr *net.IPAddr
    var req *http.Request
    var localTCPAddr net.TCPAddr
    var cli *http.Client

    for true{
        localAddr, err = net.ResolveIPAddr("ip", strings.Split(caddr,":")[0])
        if err != nil {
		log.Printf(" *ERR* %v error: %v\n", caddr,err)
                return
        }

        localTCPAddr = net.TCPAddr{
            IP: localAddr.IP,
        }

        cli = &http.Client{
            Transport: &http.Transport{
                Proxy: http.ProxyFromEnvironment,
                DialContext: (&net.Dialer{
                    LocalAddr: &localTCPAddr,
                    Timeout:   30 * time.Second,
                    KeepAlive: 30 * time.Second,
                    DualStack: true,
                }).DialContext,
                MaxIdleConns:          100,
                IdleConnTimeout:       90 * time.Second,
                TLSHandshakeTimeout:   10 * time.Second,
                ExpectContinueTimeout: 1 * time.Second,
            },
        }
///
        //now := time.Now()
        now = time.Now()
        now = now.Add(-5 *time.Second)
        timestamp = strings.Replace(now.Format("2006-01-02 15:04:05"), " ", "T", -1)
        if param == "HEALTH" {
        buf = strings.Replace(POST_monitorHEALTH_body, string("2018-08-16T10:23:19"), string(timestamp), -1)
        } else if param == "CPU"{
        buf = strings.Replace(POST_monitorCPU_body, string("2018-08-16T10:23:19"), string(timestamp), -1)
        }

        req, err = http.NewRequest("POST", "http://" + raddr + "/vnfpm/rest/api/v1.0/monitorVnf", strings.NewReader(string(buf)))
        if err != nil {
                log.Printf(" *ERR* %v http.NewRequest() error: %v\n", caddr,err)
                //return
                if MonTime != 0{
                   if Max_time > 0{
                      Max_time = Max_time - MonPeriod
                   }else{break}
                time.Sleep(time.Duration(MonPeriod) * time.Millisecond)
                continue
                }
        }
        req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "close")
        resp, err := cli.Do(req)
        log.Printf("\n *DBG* %v Request: POST http://%v/vnfpm/rest/api/v1.0/monitorVnf", caddr,raddr)
        log.Printf(" *DBG* %v Body: ", caddr, string(buf))

        if err != nil {
                log.Printf(" *ERR* %v http.Do() error: %v\n",caddr, err)
                //return
                if MonTime != 0{
                   if Max_time > 0{
                      Max_time = Max_time - MonPeriod
                   }else{break}
                time.Sleep(time.Duration(MonPeriod) * time.Millisecond)
                continue
                }
        }else{
		defer resp.Body.Close()

	        data, err = ioutil.ReadAll(resp.Body)
                //_, err = io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
		log.Printf(" *ERR* %v ioutil.ReadAll() error: %v\n",caddr, err)
			//return
	                if MonTime != 0{
			   if Max_time > 0{
			      Max_time = Max_time - MonPeriod
	                   }else{break}
		        time.Sleep(time.Duration(MonPeriod) * time.Millisecond)
			continue
	                }
		}

                log.Printf(" *DBG* %v Status: %v", caddr, string(resp.Status))
                log.Printf(" *DBG* %v Body: %v", caddr, string(data))
                if MonTime != 0{
	           if Max_time > 0{
	              Max_time = Max_time - MonPeriod
	           }else{break}
                }
	}
        time.Sleep(time.Duration(MonPeriod) * time.Millisecond)
    }
}
