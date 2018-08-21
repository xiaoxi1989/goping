package main

import (
	"fmt"
	"time"
	"github.com/tatsushid/go-fastping"
	"os"
	"net"
)

func getPing(count int, ip string)  {
	rowcounter := 0
	errCount := 0
	res := time.Duration(0)
	avg := int64(0)
	p := fastping.NewPinger()
	fmt.Println(ip)
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		res = rtt
		avg = int64(rtt) + avg
	}
	p.OnIdle = func() {
		rowcounter++

		if res == time.Duration(0) {
			fmt.Println("error")
			errCount += 0
		}

	}
	p.MaxRTT = time.Second
	//p.Debug = true
	for {
		if err = p.Run();err != nil || rowcounter == count {
			res = time.Duration(0)
			//fmt.Println(err)
			break
		}
	}

	//fmt.Println(avg/10000000)

	fmt.Printf("ip: %v, avg: %vms, loss: %v \n",ip, avg/10000000,float64(errCount/10)*100)
}

func worker(id int, jobs <-chan string, results chan<- int) {
	for j := range jobs {
		//fmt.Println("worker", id, "started  job", j)
		getPing(10, j)
		//fmt.Println("worker", id, "finished job", j)
		results <- 2
	}
}



func main()  {
	jobs := make(chan string)
	results := make(chan int, 100)

	ipList := []string{"123.125.115.110", "220.181.57.216","192.168.3.1","123.125.115.110", "220.181.57.216","192.168.3.1","123.125.115.110", "220.181.57.216","192.168.3.1","123.125.115.110", "220.181.57.216","192.168.3.1","123.125.115.110", "220.181.57.216","192.168.3.1","123.125.115.110", "220.181.57.216","192.168.3.1"}
	alen := len(ipList)
	if alen > 100 {
		alen = 100
	}

	for w := 1; w <= alen; w++ {
		go worker(w, jobs, results)
	}

	for _,x := range ipList {
		jobs <- x
	}
	close(jobs)

	for a := 1; a <= 5; a++ {
		<-results
	}


}
