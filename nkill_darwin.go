// Kills all processes listening on the given TCP ports.

package nkill

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/process"
)

type Process struct {
	Name  string
	Pid   string
	State string
	Port  int64
}

func (p *Process) Kill() error {
	pid, _ := strconv.Atoi(p.Pid)
	proc, _ := os.FindProcess(pid)
	return proc.Kill()
}

func netstat(portToKill int64) []Process {
	tcpStats := statTCP(portToKill)
	tcp6Stats := statTCP(portToKill)
	return append(tcpStats, tcp6Stats...)
}

// To get pid of all network process running on system, you must run this script
// as superuser
func statTCP(portToKill int64) []Process {
	var processes []Process

	cmd := exec.Command("lsof", "-i", ":"+strconv.Itoa(portToKill))
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	p := regexp.MustCompile("\\s+")

	reader := bufio.NewReader(bytes.NewReader(out))
	for {
		// 读取一行数据，交给后台处理
		line, isPrefix, err := reader.ReadLine()
		if len(line) > 0 {
			// fmt.Println(string(line))
			if !bytes.HasPrefix(line, []byte("COMMAND")) {
				// fmt.Println(string(line))
				s := p.Split(string(line), -1)
				if len(s) >= 2 {
					pid, err := strconv.Atoi(s[1])
					if err == nil && pid > 0 {
						proc := process.NewProcess(int32(pid))
						// fmt.Println(os.FindProcess(pid))
						exe, _ := proc.Exe()
						state, _ := proc.Status()

						p := Process{Name: exe, Pid: s[1], State: state, Port: portToKill}
						processes = append(processes, p)
					}
				}

			}
			if isPrefix {
				break
			}
		}
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}
	return processes
}

func KillPort(portToKill int64) {
	killed := false
	for {
		for _, conn := range netstat(portToKill) {
			if err := conn.Kill(); err != nil {
				log.Printf("Kill %s (pid: %s) listening on port %d failed: %s", conn.Name, conn.Pid, conn.Port, err)
			} else {
				log.Printf("Killed %s (pid: %s) listening on port %d", conn.Name, conn.Pid, conn.Port)
				killed = true
			}
		}
		if len(netstat(portToKill)) == 0 {
			break
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
	if !killed {
		log.Printf("No process found listening on port %d\n", portToKill)
	}
}

func init() {
	log.SetFlags(0)
}
