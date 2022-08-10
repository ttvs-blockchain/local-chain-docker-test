package ledger

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const (
	queryBuilder          = "docker exec cli peer chaincode query --tls --cafile /opt/home/managedblockchain-tls-chain.pem --channelID mychannel --name mycc -c %s"
	createQueryArgBuilder = "'{\"Args\":[\"CreateTX\", %q, %q]}'"
	readQueryArg          = "'{\"Args\":[\"ReadTX\", %q]}'"
	bindingLen            = 64
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
	createQueryArg := fmt.Sprintf(createQueryArgBuilder, strBinding, strTimestamp)
	res := fmt.Sprintf(queryBuilder, createQueryArg)
	return res, nil
}

func ReadTX(id string) string {
	readQueryArg := fmt.Sprintf(readQueryArg, id)
	res := fmt.Sprintf(queryBuilder, readQueryArg)
	return res
}
