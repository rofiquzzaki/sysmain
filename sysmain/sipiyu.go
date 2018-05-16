package main

import (
	"time"
	"fmt"
	"io/ioutil"
	"strconv"
	"os"
	"encoding/json"
	"strings"

	"github.com/rofiquzzaki/sysmain/sysinfo"
	zmq "github.com/pebbe/zmq4"
)

type fn func(melbu string) (metu string)

func WrapThermal(melbu string) (metu string) {
	temp := sysinfo.Thermal()
	tempstr := strconv.Itoa(temp)
	metu = "{ \"temp\" : "+tempstr+" }"
	return
}


func WrapNetUsage(melbu string) (metu string) {
	rx, tx := sysinfo.NetUsage(melbu)
	//totbw := sysinfo.BwMon(melbu)
	//fmt.Println(totbw)
	rxo := strconv.Itoa(rx)
	txo := strconv.Itoa(tx)
	metu = "{ \"rx\" : "+rxo+", \"tx\" : "+txo+" }"
	return metu
}

func WrapDiskUsage(melbu string) (metu string) {
	disk := sysinfo.NewDiskUsage(melbu)
	diskusg := disk.Usage() * 100
	kestr := strconv.FormatFloat(diskusg, 'f', 2, 64)
	metu = "{ \"diskusg\" : "+kestr+" }"
	fmt.Printf("tipe metu : %T \n", metu)
	fmt.Println(metu)
	return metu
}

func WrapUptime(melbu string) (metu string) {
	uptime := sysinfo.Uptime()
	kestr := strconv.FormatFloat(uptime, 'f', 2, 64)
	metu = "{ \"uptime\" : "+kestr+" }"
	return metu
}

func WrapCpuLoad(melbu string) (metu string) {
	smin, lmin, lbmin := sysinfo.CpuLoad()
	sminstr := strconv.FormatFloat(smin, 'f', 2, 64)
	lminstr := strconv.FormatFloat(lmin, 'f', 2, 64)
	lbminstr := strconv.FormatFloat(lbmin, 'f', 2, 64)
	metu = "{ \"loadavg\" : "+sminstr+", \"loadavg1\" : "+lminstr+", \"loadavg2\" : "+lbminstr+" }"
	return metu
}

func WrapCpuUsage(melbu string) (metu string) {
	idle, total := sysinfo.CpuUsage()
	time.Sleep(3 * time.Second)
	idle1, total1 := sysinfo.CpuUsage()

	idleTik := float64(idle1 - idle)
	totalTik := float64(total1 -total)
	sipiyu := 100 * (totalTik - idleTik) / totalTik
	metu = "{ \"cpuusg\" : "+strconv.FormatFloat(sipiyu, 'f', 2, 64)+" }"
	return metu
}

func WrapMemInfo(melbu string) (metu string) {
	total, free, used := sysinfo.MemInfo()
	totalstr := strconv.Itoa(total)
	freestr := strconv.Itoa(free)
	usedstr := strconv.FormatFloat(used, 'f', 2, 64)
	metu = "{ \"total\" : "+totalstr+", \"free\" : "+freestr+", \"used\" : "+usedstr+" }"
	return metu
}

type Config struct {
	Intrf string `json:"intrf"`
	Partn string `json:"partn"`
	Sport string `json:"sport"`
}

type Sysinfo struct {
	Usage	float64	`json:"cpu"`
	Smin	float64	`json:"load0"`
	Lmin	float64	`json:"load1"`
	Lbmin	float64	`json:"load2"`
	Uptime	float64	`json:"uptime"`
	UsedMem	float64	`json:"usedmem"`
	DiskUsg	float64	`json:"diskusg"`
	NetRx	int		`json:"netrx"`
	NetTx	int		`json:"nettx"`
}

type InputMsg struct {
	Params	string
	//Idne	string
	Method	string
}

func LoadConfiguration(file string) Config {
    var config Config
    configFile, err := os.Open(file)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)
    return config
}

func LoadInput(inputan string) InputMsg {
	var inputMsg InputMsg
	jsonParser := json.NewDecoder(strings.NewReader(inputan))
	jsonParser.Decode(&inputMsg)
	return inputMsg
}

/*
func check(e error) {
	if e != nil {
		panic(e)
	}
}
*/

func StatMtime(name string) (mtime time.Time, err error) {
	fi, err := os.Stat(name)
	if err != nil {
		return
	}
	mtime = fi.ModTime()
	return
}

func bwstartcount(ether string) {
	fmt.Println("fungsi bwstartcount")
	var keisi string
	filepath := "/aino/bw_"+ether+".log"
	bacafilemon:
	isibwtot, err := ioutil.ReadFile(filepath)
    if err != nil {
        rusak := ioutil.WriteFile(filepath, []byte("0 0 0"), 0644)
        if rusak != nil {
            fmt.Println("Error nulis file : ", rusak)
        }
        goto bacafilemon
    }

	mtime, err := StatMtime(filepath)
	if err != nil {
		fmt.Println(err)
	}
	_, bulanMtime, _ := mtime.Date()
	_, bulanNow, _ := time.Now().Date()
	fmt.Println("bulanMtime : ", bulanMtime)
    if bulanMtime != bulanNow {
		namabaru := "/aino/bw_"+ether+"_"+bulanMtime.String()+".log"
		err := os.Rename(filepath, namabaru)
		if err != nil {
			fmt.Println(err)
		}
		goto bacafilemon
	}

    isisplit := strings.Split(string(isibwtot), "\n")
    isi := strings.Fields(isisplit[0])
    bwtot, err := strconv.Atoi(isi[0])
    if err != nil {
        fmt.Println("Error bwtot : ", isibwtot, err)
    }
    bwtemp, err := strconv.Atoi(isi[1])
    if err != nil {
        fmt.Println("Error bwtemp : ", isibwtot, err)
    }
	/*
    statstart, err := strconv.Atoi(isi[2])
    if err != nil {
        fmt.Println("Error statstart : ", isibwtot, err)
    }
	*/
	keisi = strconv.Itoa(bwtot)+" "+strconv.Itoa(bwtemp)+" "+strconv.Itoa(1)
	rusak := ioutil.WriteFile(filepath, []byte(keisi), 0644)
	if rusak != nil {
		fmt.Println("Gagal nulis file : ", rusak)
	}
}

func loopingcount() {
	for {
		fmt.Println("fungsi loopingcount")
		msgtime2, err := StatMtime("/aino/bw_monthenp3s0.log")
		if err != nil {
			fmt.Println(err)
		}
		_, bulanmsg8, _ := msgtime2.Date()
		fmt.Println("barlooping sebelum bwmon: ", bulanmsg8)

		sysinfo.BwMon("enp3s0")
		time.Sleep(time.Second * 10)
	}
}

func main() {
	//saat startup, status bwMon di ganti 1
	bwstartcount("enp3s0")

	//mapping wrap function
	m := map[string] fn {
		"net" : WrapNetUsage,
		"disk" : WrapDiskUsage,
		"uptime" : WrapUptime,
		"loadavg" : WrapCpuLoad,
		"cpuusg" : WrapCpuUsage,
		"memory" : WrapMemInfo,
		"thermal" : WrapThermal,
	}


	//wkwk := m["net"]("enp3s0")
	//fmt.Println(wkwk)

	//config := LoadConfiguration("konf.json")

	ruter, _ := zmq.NewSocket(zmq.ROUTER)
	defer ruter.Close()

	rutertcp, _ := zmq.NewSocket(zmq.ROUTER)
	defer rutertcp.Close()

	ruter.Bind("ipc:///tmp/ngawur")
	rutertcp.Bind("tcp://*:5671")

	poller := zmq.NewPoller()
	poller.Add(rutertcp, zmq.POLLIN)
	poller.Add(ruter, zmq.POLLIN)

	go loopingcount()

	for {
		fmt.Println("for poller")
		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket
			s {
			case ruter:
				idne, _ := s.Recv(0)
				isine, _ := s.Recv(0)
				fmt.Println(idne, isine)
				masukan := LoadInput(isine)
				if val, ok := m[masukan.Method]; ok {
					nyoh := val(masukan.Params)
					ruter.Send(idne, zmq.SNDMORE)
					ruter.Send(nyoh, 0)
				} else {
					ruter.Send(idne, zmq.SNDMORE)
					ruter.Send("salah", 0)
				}
			case rutertcp:
				idne, _ := s.Recv(0)
				isine, _ := s.Recv(0)
				fmt.Println(idne, isine)
				masukan := LoadInput(isine)
				if val, ok := m[masukan.Method]; ok {
					nyoh := val(masukan.Params)
					rutertcp.Send(idne, zmq.SNDMORE)
					rutertcp.Send(nyoh, 0)
				} else {
					rutertcp.Send(idne, zmq.SNDMORE)
					rutertcp.Send("salah", 0)
				}
			}
		}
	}

	/*
	rx, tx := sysinfo.NetUsage(config.Intrf)
	disk := sysinfo.NewDiskUsage(config.Partn)
	diskusg := disk.Usage()*100
	uptime := sysinfo.Uptime()

	totol, freo, usedmem := sysinfo.MemInfo()

	smin, lmin, lbmin := sysinfo.CpuLoad()
	idle, total := sysinfo.CpuUsage()
	time.Sleep(2 * time.Second)
	idle1, total1 := sysinfo.CpuUsage()

	idleTik := float64(idle1 - idle)
	totalTik := float64(total1 - total)
	sipiyu := 100 * (totalTik - idleTik) / totalTik

	jcpu := Sysinfo{sipiyu, smin, lmin, lbmin, uptime, usedmem, diskusg, rx, tx}
	b, _ := json.MarshalIndent(jcpu, "", "    ")
	ioutil.WriteFile("sistem.json", b, 0644)
	fmt.Printf("%+v \n", jcpu)
	fmt.Println("total ", totol, "free ", freo, "terpakai ", usedmem)
	*/
}
