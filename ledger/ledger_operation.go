package ledger

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const (
	cliBuilder       = "peer chaincode %s --channelID mychannel --name mycc -c %s"
	cliInvoke        = "invoke"
	cliQuery         = "query"
	createArgBuilder = "'{\"Args\":[\"CreateTX\", %q, %q]}'"
	readArgBuilder   = "'{\"Args\":[\"ReadTX\", %q]}'"
	bindingLen       = 64
)

func DummyCreatTX() (string, error) {
	byteBinding := make([]byte, bindingLen)
	_, err := rand.Read(byteBinding)
	if err != nil {
		return "", err
	}
	strBinding := hex.EncodeToString(byteBinding)
	timestamp := time.Now().UnixNano()
	strTimestamp := strconv.Itoa(int(timestamp))
	createQueryArg := fmt.Sprintf(createArgBuilder, strBinding, strTimestamp)
	res := fmt.Sprintf(cliBuilder, cliInvoke, createQueryArg)
	return res, nil
}

func ReadTX(id string) string {
	readQueryArg := fmt.Sprintf(readArgBuilder, id)
	res := fmt.Sprintf(cliBuilder, cliQuery, readQueryArg)
	return res
}
