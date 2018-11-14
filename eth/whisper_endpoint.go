package eth

import (
	"context"
	"log"
	"fmt"
	"github.com/ShyftNetwork/go-empyrean/whisper/shhclient"
	//"github.com/ShyftNetwork/go-empyrean/whisper/whisperv6"

	"github.com/ShyftNetwork/go-empyrean/whisper/whisperv6"
	"github.com/ShyftNetwork/go-empyrean/common/hexutil"
)

func NewWhisperEndPoint() {
	client, err := shhclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}
	keyID, err := client.NewKeyPair(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(keyID) // 0ec5cfe4e215239756054992dbc2e10f011db1cdfc88b9ba6301e2f9ea1b58d2

	//func (sc *Client) PublicKey(ctx context.Context, id string) ([]byte, error) {

	pubKey, errf := client.PublicKey(context.Background(), keyID)
	if errf != nil {
		log.Fatal(errf)
	}
	fmt.Println("pub key")
	fmt.Println(hexutil.Encode(pubKey))


	_ = client // we'll be using this in the next section
	fmt.Println("we have a whisper connection")

	messages := make(chan *whisperv6.Message)

	criteria := whisperv6.Criteria{
		PrivateKeyID: keyID,
	}

	sub, err2 := client.SubscribeMessages(context.Background(), criteria, messages)
	//_, err := client.SubscribeMessages(context.Background(), criteria, messages)
	if err2 != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case message := <-messages:
			fmt.Printf(string(message.Payload)) // "Hello"
		}
	}

}
