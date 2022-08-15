package ledger

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const (
	cliBuilder       = "peer chaincode %s --tls --cafile /opt/home/managedblockchain-tls-chain.pem --channelID ourchannel --name myjointcc -c %s"
	cliInvoke        = "invoke"
	cliQuery         = "query"
	createArgBuilder = "'{\"Args\":[\"CreateTX\", %q]}'"
	readArgBuilder   = "'{\"Args\":[\"TXExists\", %q]}'"
	bindingLen       = 32
)

func DummyCreatTX() (string, error) {
	byteHash := make([]byte, bindingLen)
	_, err := rand.Read(byteHash)
	if err != nil {
		return "", err
	}
	strBinding := hex.EncodeToString(byteHash)
	createQueryArg := fmt.Sprintf(createArgBuilder, strBinding)
	res := fmt.Sprintf(cliBuilder, cliInvoke, createQueryArg)
	return res, nil
}

func ReadTX(id string) string {
	readQueryArg := fmt.Sprintf(readArgBuilder, id)
	res := fmt.Sprintf(cliBuilder, cliQuery, readQueryArg)
	return res
}
