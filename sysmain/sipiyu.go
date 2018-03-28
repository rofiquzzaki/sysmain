package main

import (
	"time"
	"fmt"
	"io/ioutil"
	"strconv"
	//"strings"
	"os"
	"encoding/json"

	"github.com/rofiquzzaki/sysmain/sysinfo"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Create("sistem.info")
	check(err)
	defer f.Close()

	uptime := sysinfo.Uptime()

	smin, lmin, lbmin := sysinfo.CpuLoad()
	fmt.Println("Load Average CPU : ", smin, lmin, lbmin)
	ss, _, _ := sysinfo.CpuLoad()
	fmt.Println("gur siji : ", ss)
	idle, total := sysinfo.CpuUsage()
	time.Sleep(1 * time.Second)
	idle1, total1 := sysinfo.CpuUsage()

	idleTik := float64(idle1 - idle)
	totalTik := float64(total1 - total)
	sipiyu := 100 * (totalTik - idleTik) / totalTik
	fmt.Printf("Penggunaan CPU : %f %%\n", sipiyu)
	fmt.Printf("sipiyu %T %+v\n", sipiyu, sipiyu)

	sipiye := strconv.FormatFloat(sipiyu, 'f', 3, 64)
	fmt.Printf("sipiye %T %+v\n", sipiye, sipiye)
	result := []byte(sipiye)

	n1, err := f.Write(result)
	check(err)
	fmt.Printf("wrote %d bytes\n", n1)

	f.Sync()

	jcpu := sysinfo.Sysinfo{sipiyu, smin, lmin, lbmin, uptime}
	b, _ := json.MarshalIndent(jcpu, "", "    ")
	err = ioutil.WriteFile("sistem.json", b, 0644)
	fmt.Printf("%+v", jcpu)

	//fmt.Printf("result:%+v\n", result)
}
