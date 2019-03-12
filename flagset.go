package main

import(
	"flag"
	"log"
	"fmt"
)

func flagInit(){
//flag.StringVar(&ipRange, "cip", "127.0.0.1", "Address range for VNFs\n*****Examples*****\nFor Single VNF:  -cip \"127.0.0.1\" \nFor More VNFs(if 10):  -cip \"127.0.0.1-10\"\n")
flag.IntVar(&port, "vnfPort", 8080, "VNF Port")
flag.StringVar(&raddr, "vnfmAddr", "0.0.0.0:8080", "VNFM Port")
flag.StringVar(&csvFile, "vnfFile", "", "CSV file, which contains VNF(s) data")
flag.IntVar(&rTime, "rTime", 0, "VNF run type:\n--> To run parallel -rTime 0\nTo run serial -rTime <seconds>")
flag.IntVar(&DBLretry, "DBLretry", 0, "DoorBell re-try timer")
flag.Parse()

if csvFile == ""{
fmt.Println(" *ERR* input CSV file is not provide")
log.Fatal(" *ERR* input CSV file is not provide")
}
}
