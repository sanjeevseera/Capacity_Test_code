package main

import(
        "strconv"
        "fmt"
	"log"
        "time"
	"io/ioutil"
        "os/signal"
        "os"
	"strings"
	"encoding/csv"
	"io"
	"sync"
	_ "net/http/pprof"
)
//var ipRange string
var raddr, t, csvFile string
var vnfips []string
var port, rTime, DBLretry int
var wg_vnf sync.WaitGroup

func main(){

// Interrupt Handle *START*
    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel, os.Interrupt)
    go func() {
        sig := <-signalChannel
        switch sig {
        case os.Interrupt:
             os.Exit(0)
        }
    }()

argsWithProg := os.Args
if len(argsWithProg) == 1{
fmt.Println("|****\t****\t****\t****\t****\t****\t****\t****\t****\t****\t****\t****\t****|\n|")
fmt.Println("ERROR: No Arguments provided\nScript Example::\n./vnf_simulator -vnfPort 8008 -vnfmAddr "+"0.0.0.0"+":"+"8080 -vnfFile \"VNF_ConnectionPoint_Info.db\" -rTime <int> -DBLretry <int>")
fmt.Println("|\n|****\t****\t****\t****\t****\t****\t****\t****\t****\t****\t****\t****\t****|\n")
os.Exit(0)
}

go statsPrint()
// Interrupt Handle *END*
// Logger *START*
t := strings.Split(time.Now().String()," ")[1]
t = strings.Replace(t, ":", "-", -1)
t = strings.Replace(t, ".", "-", -1)
f, err := os.OpenFile("/tmp/vnf_"+t+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
os.Remove("/tmp/latest.log")
os.Symlink("/tmp/vnf_"+t+".log", "/tmp/latest.log")
if err != nil {
     log.Fatal(err)
}
//defer f.Close()
defer f.Close()
log.SetOutput(f)
// Logger *END*

flagInit()  // Parse command line Arguments
log.Println(" *INFO* VNF port: ", port)
log.Println(" *INFO* VNFM address: ", raddr)
log.Println(" *INFO* CSV file: ", csvFile)
if rTime == 0{
log.Println(" *INFO* Start VNFs in parallel")
}else{
log.Printf(" *INFO* Start VNFs in serial with the dutarion of %vsec(s)", rTime)
}



csvdata, err := ioutil.ReadFile(csvFile)
if err != nil {
	log.Fatal(" *ERR* Read CSV File: ", err)
}
r := csv.NewReader(strings.NewReader(string(csvdata)))
r.Comma = ';'
r.LazyQuotes = true // If LazyQuotes is true, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field.
r.FieldsPerRecord = -1

for{
        record, err := r.Read()
        if err == io.EOF {
                break
        }
        if err != nil {
                log.Fatal(" *ERR* Read CSV Data: ",err)
        }
        vnfips=append(vnfips,string(record[2]))
}

rs := csv.NewReader(strings.NewReader(string(csvdata)))
rs.Comma = ';'
rs.LazyQuotes = true // If LazyQuotes is true, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field.
rs.FieldsPerRecord = -1
for{
	records, err := rs.Read()
	if err == io.EOF {
                break
        }
	if err != nil {
                log.Fatal(" *ERR* Read CSV Data: ",err)
        }
	log.Printf(" *INFO* VNF IP: [%v], uuid: [%v], DoorBell Body: [%v]\n", records[2], records[1], records[10])
	//vnfips=append(vnfips,string(record[2]))
	caddr := string(records[2]) + ":" + strconv.Itoa(port)
	wg_vnf.Add(1)
	go NewVnf(caddr, string(records[1]), string(records[10]))
	if rTime != 0{
	time.Sleep(time.Duration(rTime) * time.Second)
	}else{time.Sleep(10 *  time.Millisecond)}
}
wg_vnf.Wait()
log.Printf(" *INFO * received VNF Wait group")
select{ }
//time.Sleep(60 * time.Second)
}
