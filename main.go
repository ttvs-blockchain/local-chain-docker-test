package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	shell "github.com/rfyiamcool/go-shell"
)

const (
	// app          = "peer"
	// arg0         = "chaincode"
	// arg1         = "query"
	// arg2         = "--tls"
	// arg3         = "--cafile"
	// arg4         = "/opt/home/managedblockchain-tls-chain.pem"
	// arg5         = "--channelID"
	// arg6         = "mychannel"
	// arg7         = "--name"
	// arg8         = "mycc"
	queryBuilder = "peer chaincode query --tls --cafile /opt/home/managedblockchain-tls-chain.pem --channelID mychannel --name mycc -c '{\"Args\":[\"CreateTX\", \"%s\", \"%s\"]}'"
	bindingLen   = 64
)

func main() {
	tries := flag.Int("try", 100, "number of tries")
	flag.Parse()
	f, err := os.OpenFile("test_"+time.Now().Local().String()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
		query, err := composeQuery()
		handleError(err)
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

func composeQuery() (string, error) {
	byteBinding := make([]byte, bindingLen)
	_, err := rand.Read(byteBinding)
	if err != nil {
		return "", err
	}
	strBinding := hex.EncodeToString(byteBinding)
	timestamp := time.Now().UnixNano()
	strTimestamp := strconv.Itoa(int(timestamp))
	res := fmt.Sprintf(queryBuilder, strBinding, strTimestamp)
	return res, nil
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
