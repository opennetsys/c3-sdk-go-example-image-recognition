package main

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"time"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/common/txparamcoder"
	"github.com/c3systems/c3-go/core/chain/statechain"
	methodTypes "github.com/c3systems/c3-go/core/types/methods"
	"github.com/c3systems/c3-go/node"
	nodetypes "github.com/c3systems/c3-go/node/types"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func main() {
	imageHash := os.Getenv("IMAGEID")
	privPEM := "demo/priv2.pem"
	nodeURI := "/ip4/0.0.0.0/tcp/9004"
	peer := os.Getenv("PEERID")
	dataDir := "~/.c3"
	n := new(node.Service)
	ready := make(chan bool)

	go func() {
		go func() {
			err := node.Start(n, &nodetypes.Config{
				URI:     nodeURI,
				Peer:    peer,
				DataDir: dataDir,
				Keys: nodetypes.Keys{
					PEMFile:  privPEM,
					Password: "",
				},
				BlockDifficulty: 5,
			})

			if err != nil {
				log.Error(err)
			}
		}()

		time.Sleep(10 * time.Second)
		ready <- true
	}()

	<-ready

	priv, err := c3crypto.ReadPrivateKeyFromPem(privPEM, nil)
	if err != nil {
		log.Error(err)
	}

	pub, err := c3crypto.GetPublicKey(priv)
	if err != nil {
		log.Error(err)
	}

	encodedPub, err := c3crypto.EncodeAddress(pub)
	if err != nil {
		log.Error(err)
	}

	tx1 := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: imageHash,
		Method:    methodTypes.Deploy,
		Payload:   []byte(``),
		From:      encodedPub,
	})

	fileBytes, err := ioutil.ReadFile("images/cat/cat.jpg")
	if err != nil {
		log.Error(err)
	}

	payload := txparamcoder.ToJSONArray(
		txparamcoder.EncodeMethodName("processImage"),
		txparamcoder.EncodeParam(hex.EncodeToString(fileBytes)),
		txparamcoder.EncodeParam("jpg"),
	)

	tx2 := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: imageHash,
		Method:    methodTypes.InvokeMethod,
		Payload:   payload,
		From:      encodedPub,
	})

	tx := tx2
	if os.Getenv("METHOD") == "deploy" {
		tx = tx1
	}

	err = tx.SetHash()
	if err != nil {
		log.Error(err)
	}

	err = tx.SetSig(priv)
	if err != nil {
		log.Error(err)
	}

	resp, err := n.BroadcastTransaction(tx)
	if err != nil {
		log.Error(err)
	}

	if resp.TxHash == nil {
		log.Error("expected txhash")
	}

	spew.Dump(resp)
	time.Sleep(10 * time.Second) // needed
}
