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

func main() {
	tries := flag.Int("try", 100, "number of tries")
	mode := flag.Bool("mode", true, "if true then run CreateTX, otherwise run ReadTX")
	flag.Parse()
	f, err := os.OpenFile(time.Now().Format("2006-01-02-15-04-05")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	handleError(err)
	defer func(f *os.File) {
		err := f.Close()
		handleError(err)
	}(f)

	var idList []string
	if *mode {
		f, err := os.Open("id_file/test_2_2_id.txt")
		handleError(err)
		defer func(f *os.File) {
			err := f.Close()
			handleError(err)
		}(f)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			idList = append(idList, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			handleError(err)
		}
		fmt.Println("Read TX test")
	} else {
		fmt.Println("Create TX test")
	}

	var totalTime float64
	for i := 0; i < *tries; i++ {
		fmt.Println("No.", i)
		_, err := f.WriteString(strconv.Itoa(i) + "\n")
		handleError(err)
		var query string
		if *mode {
			query, err = ledger.DummyCreatTX()
			handleError(err)
		} else {
			randIdx := rand.Intn(len(idList))
			query = ledger.ReadTX(idList[randIdx])
		}
		cmd := shell.NewCommand(query)
		fmt.Println(query)
		_, err = f.WriteString(cmd.Bash + "\n")
		handleError(err)
		cmd.Run()
		status := cmd.Status
		if status.ExitCode != 0 {
			panic(status.Error)
		}
		_, err = f.WriteString(status.Output)
		handleError(err)
		costTimeSeconds := status.CostTime.Seconds()
		totalTime += costTimeSeconds
		_, err = f.WriteString(fmt.Sprintf("%f\n", costTimeSeconds))
		handleError(err)
	}
	avgTime := totalTime / float64(*tries)
	fmt.Println("Average time: " + fmt.Sprintf("%f", avgTime))
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
