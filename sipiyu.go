package main

import (
	"time"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"os"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}


type Sysinfo struct {
	Usage	float64	`json:"cpu"`
	Smin	float64	`json:"load0"`
	Lmin	float64	`json:"load1"`
	Lbmin	float64	`json:"load2"`
	Uptime	float64 `json:"uptime"`
}

func uptime() (uptime float64) {
	isi, err := ioutil.ReadFile("/proc/uptime")
	if err !=nil {
		return
	}
	lines := strings.Split(string(isi), " ")
	for i :=0; i < 1; i++ {
		val, err := strconv.ParseFloat(lines[i], 64)
		if err != nil {
			fmt.Println("Error: ", i, lines[i], err)
		}
		if i == 0 {
			uptime = val
		}
	}
	return
}

func cpuUsage() (idle, total uint64) {
	isi, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(isi), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
			return
		}
	}
	return
}

func cpuLoad() (smin float64, lmin float64, lbmin float64) {
	isi, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return
	}
	lines := strings.Split(string(isi), " ")
	for i :=0; i < 3; i++ {
		val, err := strconv.ParseFloat(lines[i], 64)
		if err != nil {
			fmt.Println("Error: ", i, lines[i], err)
		}
			switch i {
				case 0:
					smin = val
				case 1:
					lmin = val
				case 2:
					lbmin = val
				}
	}
	return
}

func main() {
	f, err := os.Create("sistem.info")
	check(err)
	defer f.Close()

	uptime := uptime()

	smin, lmin, lbmin := cpuLoad()
	fmt.Println("Load Average CPU : ", smin, lmin, lbmin)
	ss, _, _ := cpuLoad()
	fmt.Println("gur siji : ", ss)
	idle, total := cpuUsage()
	time.Sleep(1 * time.Second)
	idle1, total1 := cpuUsage()

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

	jcpu := Sysinfo{sipiyu, smin, lmin, lbmin, uptime}
	b, _ := json.MarshalIndent(jcpu, "", "    ")
	err = ioutil.WriteFile("sistem.json", b, 0644)
	fmt.Printf("%+v", jcpu)

	//fmt.Printf("result:%+v\n", result)
}
