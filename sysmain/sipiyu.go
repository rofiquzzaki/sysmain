package main

import (
	"time"
	"fmt"
	"io/ioutil"
	//"strconv"
	//"os"
	"encoding/json"

	"github.com/rofiquzzaki/sysmain/sysinfo"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//niatnya mau buat goroutine
//func aptem() 

func main() {
	uptime := sysinfo.Uptime()

	totol, freo, usedmem := sysinfo.MemInfo()

	smin, lmin, lbmin := sysinfo.CpuLoad()
	fmt.Println("Load Average CPU : ", smin, lmin, lbmin)
	idle, total := sysinfo.CpuUsage()
	time.Sleep(2 * time.Second)
	idle1, total1 := sysinfo.CpuUsage()

	idleTik := float64(idle1 - idle)
	totalTik := float64(total1 - total)
	sipiyu := 100 * (totalTik - idleTik) / totalTik
	fmt.Printf("Penggunaan CPU : %f %%\n", sipiyu)

	jcpu := sysinfo.Sysinfo{sipiyu, smin, lmin, lbmin, uptime, usedmem}
	b, _ := json.MarshalIndent(jcpu, "", "    ")
	ioutil.WriteFile("sistem.json", b, 0644)
	fmt.Printf("%+v \n", jcpu)
	fmt.Println("total ", totol, "free ", freo, "terpakai ", usedmem)

	//fmt.Printf("result:%+v\n", result)
}
