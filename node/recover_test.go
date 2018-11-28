package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"testing"

	"github.com/ShyftNetwork/go-empyrean/common/hexutil"
	"github.com/ShyftNetwork/go-empyrean/crypto"
	"github.com/ShyftNetwork/go-empyrean/signer/core"
)

func TestEC(t *testing.T) {
	// ecrecover should return
	//7da99df96259305ee38c9fa9e9d551118b12ec3b
	sig := "0x75ea60ce26d00a43d89833bbee60a1e87a9cded792a76b1f35c4bccca20288eb53c8651d92da37c030b15d64bb49e855f2ac9eecbbd4b1d8e589409efd69427b1c"
	sigByteArray, _ := hexutil.Decode(sig)
	var sighex = hexutil.Bytes(sigByteArray)
	sighex[64] -= 27
	msgValue := "tneiftw"
	msgHash, _ := core.SignHash(hexutil.Bytes(msgValue))
	rpk, err := crypto.Ecrecover(hexutil.Bytes(msgHash), sighex)

	if err != nil {
		fmt.Println("Error in EcRecover FOO", err)
	}

	fmt.Println("rpk", rpk)

	x, y := elliptic.Unmarshal(crypto.S256(), rpk)
	foo := &ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}
	//pubKey, err := crypto.DecompressPubkey(rpk)
	//if err != nil {
	//	fmt.Println("Error in Decompress", err)
	//}
	recoveredAddr := crypto.PubkeyToAddress(*foo)
	fmt.Println("Client connected with address :", recoveredAddr.Hex())

}
