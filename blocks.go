package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type block struct {
	command string
	rate    int
	Val     string
}

func (b *block) Exec(ch chan string) {
	for {
		val := cmd(b.command)
		ch <- val
		b.Val = val
		time.Sleep(time.Duration(b.rate) * time.Second)
	}
}

type blocks struct {
	Weather block
	Network block
	Cpu     block
	Temp    block
	Disk    block
	Mem     block
	Ads     block
	Date    block
}

func (bs blocks) Show() {
	ch := make(chan string)
	go bs.Weather.Exec(ch)
	go bs.Network.Exec(ch)
	go bs.Cpu.Exec(ch)
	go bs.Temp.Exec(ch)
	go bs.Disk.Exec(ch)
	go bs.Mem.Exec(ch)
	go bs.Ads.Exec(ch)
	go bs.Date.Exec(ch)

	for {
		select {
		case _, ok := <-ch:
			if ok {
				go updateBlocks(fmt.Sprintf("^c#282C34^^b#ABB2BF^ %s ^b#E06C75^ %s ^b#E5C07B^  %s ^b#C678DD^  %s ^b#61AFEF^  %s ^b#98C379^  %s ^b#56B6C2^  %s ^d^  %s",
					bs.Weather.Val,
					bs.Network.Val,
					bs.Cpu.Val,
					bs.Temp.Val,
					bs.Disk.Val,
					bs.Mem.Val,
					bs.Ads.Val,
					bs.Date.Val,
				))
			}
		}
	}
}

func cmd(command string) string {
	cmd := exec.Command(command)
	stdout, err := cmd.Output()
	resl := ""
	if err != nil {
		fmt.Println(err.Error())
	}
	resl = string(stdout)
	return resl
}

func updateBlocks(msg string) {
	cmd := exec.Command("xsetroot", "-name", msg)
	cmd.Run()
}

func main() {
	bs := blocks{
		Weather: block{
			command: "weather.sh",
			rate:    14400,
			Val:     "",
		},
		Network: block{
			command: "network",
			rate:    5,
			Val:     "",
		},
		Cpu: block{
			command: "cpu",
			rate:    5,
			Val:     "",
		},
		Temp: block{
			command: "temp.sh",
			rate:    5,
			Val:     "",
		},
		Disk: block{
			command: "disk.sh",
			rate:    900,
			Val:     "",
		},
		Mem: block{
			command: "mem.sh",
			rate:    10,
			Val:     "",
		},
		Ads: block{
			command: "pihole",
			rate:    1800,
			Val:     "",
		},
		Date: block{
			command: "date.sh",
			rate:    30,
			Val:     "",
		},
	}

	exit := make(chan string)
	bs.Show()
	for {
		select {
		case <-exit:
			os.Exit(0)
		}
	}

}
