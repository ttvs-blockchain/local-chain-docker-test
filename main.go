package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	shell "github.com/rfyiamcool/go-shell"
	"github.com/ttvs-blockchain/local-chain-docker-test/ledger"
)

const (
	createPrefix = "c"
	readPrefix   = "r"
)

func main() {
	tries := flag.Int("try", 100, "number of tries")
	read := flag.Bool("read", false, "if true then run ReadTX, otherwise run CreateTX")
	flag.Parse()
	var prefix string
	var idList []string
	if *read {
		f, err := os.Open("id_file/test_5_id.txt")
		handleError(err)
		defer func(f *os.File) {
			err := f.Close()
			handleError(err)
		}(f)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			id := scanner.Text()
			idList = append(idList, id)
		}
		if err := scanner.Err(); err != nil {
			handleError(err)
		}
		prefix = readPrefix
		fmt.Println("Read TX test")
	} else {
		prefix = createPrefix
		fmt.Println("Create TX test")
	}

	f, err := os.OpenFile(
		prefix+time.Now().Format("2006-01-02-15-04-05")+".log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
	)
	handleError(err)
	defer func(f *os.File) {
		err := f.Close()
		handleError(err)
	}(f)

	var totalTime float64
	for i := 0; i < *tries; i++ {
		fmt.Println("No.", i)
		_, err := f.WriteString(strconv.Itoa(i) + "\n")
		handleError(err)
		var query string
		if *read {
			randIdx := rand.Intn(len(idList))
			query = ledger.ReadTX(idList[randIdx])
		} else {
			query, err = ledger.DummyCreatTX()
			handleError(err)
		}
		cmd := shell.NewCommand(query)
		fmt.Println(query)
		_, err = f.WriteString(cmd.Bash + "\n")
		handleError(err)
		err = cmd.Run()
		handleError(err)
		status := cmd.Status
		if status.ExitCode != 0 {
			fmt.Println(status.Output)
			panic(status.Error)
		}
		_, err = f.WriteString(status.Output)
		handleError(err)
		costTimeSeconds := status.CostTime.Seconds()
		totalTime += costTimeSeconds
		_, err = f.WriteString(fmt.Sprintf("%f\n", costTimeSeconds))
		handleError(err)
		if i%1000 == 0 && i != 0 {
			fmt.Println("Idling...")
			time.Sleep(time.Second * 10)
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
	avgTime := totalTime / float64(*tries)
	fmt.Println("Average time: " + fmt.Sprintf("%f", avgTime))
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
