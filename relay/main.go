package main

import (
	context2 "context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	relay "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"

	secio "github.com/libp2p/go-libp2p-secio"
)

type NullValidator struct {

}

func (v *NullValidator) Validate(key string, value []byte) error {
	return nil
}

func (v *NullValidator) Select(key string, values [][]byte) (int, error) {
	return 0, nil
}

//const Protocol = "kadbox"

func main() {
	listen := flag.String("listen", "/ip4/0.0.0.0/tcp/5000", "The listen address")
	flag.Parse()

	ctx, cancel := context2.WithCancel(context2.Background())
	defer cancel()

	//routing := libp2p.Routing(func(host host.Host) (routing.PeerRouting, error) {
	//	return dual.New(ctx, host)
	//})

	listenAddress := libp2p.ListenAddrStrings(*listen)

	enableRelay := libp2p.EnableRelay(relay.OptActive, relay.OptHop)

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}

	identity := libp2p.Identity(priv)

	security := libp2p.Security(secio.ID, secio.New)

	host, err := libp2p.New(
		ctx,
		//routing,
		listenAddress,
		enableRelay,
		identity,
		security,
	)

	if err != nil {
		panic(err)
	}

	for _, addr := range host.Addrs() {
		fmt.Printf("Addr: %s/p2p/%s\n", addr, host.ID().Pretty())
	}

	select {}
}
