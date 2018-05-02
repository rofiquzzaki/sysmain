package sysinfo

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"path/filepath"
)

//func BwUsage(ether string) (rx int, tx int) {
//}

func Thermal() (temp int) {
	files, err := filepath.Glob("/sys/class/thermal/thermal_zone*")
	if err != nil {
		return
	}
	jumlah := len(files)
	path := files[jumlah-1]
	//var path string
	//for _, x := range files {
	//	path = x
	//}
	isi, err := ioutil.ReadFile(path+"/temp")
	if err !=nil {
		return
	}
	isisplit := strings.Split(string(isi), "\n")
	temp, err = strconv.Atoi(isisplit[0])
	return
}

//func Thermal() (temp int) {
//	isi, err :=

//Network Interface belum dinamis
//Percobaan pakai argument ether
func NetUsage(ether string) (rx int, tx int) {
	isi, err := ioutil.ReadFile("/sys/class/net/"+ether+"/statistics/rx_bytes")
	isi1, err := ioutil.ReadFile("/sys/class/net/"+ether+"/statistics/tx_bytes")
	if err != nil {
		return
	}
	isisplit := strings.Split(string(isi), "\n")
	isi1split := strings.Split(string(isi1), "\n")
	rx, err = strconv.Atoi(isisplit[0])
	if err != nil {
		fmt.Println("Error: ", isi, err)
	}
	tx, err = strconv.Atoi(isi1split[0])
	if err != nil {
		fmt.Println("Error: ", isi1, err)
	}
	return
}

//DiskUsage ntah kenapa di struct
type DiskUsage struct {
	stat *syscall.Statfs_t
}

//Return volumePath, harus valid path
func NewDiskUsage(volumePath string) *DiskUsage {
	var stat syscall.Statfs_t
	syscall.Statfs(volumePath, &stat)
	return &DiskUsage{&stat}
}

//Free bytes on file system
func (this *DiskUsage) Free() uint64 {
	return this.stat.Bfree * uint64(this.stat.Bsize)
}

//Available bytes on file system to an unprivileged user
func (this *DiskUsage) Available() uint64 {
	return this.stat.Bavail * uint64(this.stat.Bsize)
}

//Total size of the file system
func (this *DiskUsage) Size() uint64 {
	return this.stat.Blocks * uint64(this.stat.Bsize)
}

//Total bytes used of the file system
func (this *DiskUsage) Used() uint64 {
	return this.Size() - this.Free()
}

//Percentage of use of the file system
func (this *DiskUsage) Usage() float64 {
	return float64(this.Used()) / float64(this.Size())
}

func Uptime() (uptime float64) {
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

func CpuUsage() (idle, total uint64) {
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

func CpuLoad() (smin float64, lmin float64, lbmin float64) {
	isi, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return
	}
	lines := strings.Split(string(isi), " ")
	for i := 0; i < 3; i++ {
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

func MemInfo() (total int, free int, used float64) {
    isi, err := ioutil.ReadFile("/proc/meminfo")
    if err != nil {
        return
    }

    lines := strings.Split(string(isi), "\n")
    for _, line := range(lines) {
        fields := strings.Fields(line)
		//soalnya lines terakhir fields[] kosong
		if len(fields) <= 0 {
			break
		}
        if fields[0] == "MemTotal:" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                if i == 1 {
                    val, err := strconv.Atoi(fields[i])
                    if err != nil {
                        fmt.Println("Error: ", i, fields[i], err)
                    }
                    total = val
					fmt.Println("total : ", total)
                }
            }
        } else if fields[0] == "MemAvailable:" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                if i == 1 {
                    val, err := strconv.Atoi(fields[i])
                    if err != nil {
                        fmt.Println("Error: ", i, fields[i], err)
                    }
                    free = val
					fmt.Println("free avail :", free)
                }
            }
			used = ((float64(total) - float64(free)) / float64(total)) * 100
			fmt.Println("used dari avail :", used)
		} else if fields[0] == "MemFree:" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				if i == 1 {
					val, err := strconv.Atoi(fields[i])
					if err != nil {
						fmt.Println("Error: ", i, fields[i], err)
					}
					free = val
					fmt.Println("free free :", free)
				}
			}
			used = ((float64(total) - float64(free)) / float64(total)) * 100
			fmt.Println("used dari free :", used)
		}
    }
    return
}
