package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	app        = "peer"
	arg0       = "chaincode"
	arg1       = "query"
	arg2       = "--tls"
	arg3       = "--cafile"
	arg4       = "/opt/home/managedblockchain-tls-chain.pem"
	arg5       = "--channelID"
	arg6       = "mychannel"
	arg7       = "--name"
	arg8       = "mycc"
	arg9       = "-c"
	bindingLen = 64
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
		query, err := composeQuery()
		handleError(err)
		cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, query)
		fmt.Println(cmd)
		_, err = f.WriteString(cmd.String() + "\n")
		handleError(err)
		start := time.Now()
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}
		handleError(err)
		fmt.Println(stdout.String())
		_, err = f.WriteString(stdout.String() + "\n")
		handleError(err)
		timeInterval := time.Since(start).Seconds()
		totalTime += timeInterval
		timeString := fmt.Sprintf("%f", timeInterval)
		_, err = f.WriteString(timeString + "\n")
		handleError(err)
	}
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
	var builder strings.Builder
	builder.WriteString(`'{"Args":["CreateTX", "`)
	builder.WriteString(strBinding)
	builder.WriteString(`", "`)
	builder.WriteString(strTimestamp)
	builder.WriteString(`"]}'`)
	return builder.String(), nil
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
