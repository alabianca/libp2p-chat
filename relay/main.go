package main

import (
	secio "github.com/libp2p/go-libp2p-secio"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	//dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"

	//context2 "context"
	//"crypto/rand"
	autonat "github.com/libp2p/go-libp2p-autonat-svc"
	mrand "math/rand"
	//"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	//relay "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	//"github.com/libp2p/go-libp2p-core/host"
	//"github.com/libp2p/go-libp2p-core/routing"
	//dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"golang.org/x/net/context"

	//secio "github.com/libp2p/go-libp2p-secio"
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
	//listen := flag.String("listen", "/ip4/0.0.0.0/tcp/5000", "The listen address")

	ctx := context.Background()

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/5000")

	r := mrand.New(mrand.NewSource(int64(10)))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
		libp2p.EnableNATService(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("This node: ", host.ID().Pretty(), " ", host.Addrs())

	_, err = dual.New(ctx, host)
	if err != nil {
		panic(err)
	}
	//flag.Parse()
	//
	//ctx, cancel := context2.WithCancel(context2.Background())
	//defer cancel()
	//
	//var ddht *dual.DHT
	//routing := libp2p.Routing(func(host host.Host) (routing.PeerRouting, error) {
	//	var err error
	//	ddht, err = dual.New(ctx, host)
	//	return ddht,err
	//})
	//
	//listenAddress := libp2p.ListenAddrStrings(*listen)
	//
	//enableRelay := libp2p.EnableRelay(relay.OptActive, relay.OptHop)
	//
	//priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	//if err != nil {
	//	panic(err)
	//}
	//
	//identity := libp2p.Identity(priv)
	//
	//security := libp2p.Security(secio.ID, secio.New)
	//
	//host, err := libp2p.New(
	//	ctx,
	//	routing,
	//	listenAddress,
	//	enableRelay,
	//	identity,
	//	security,
	//)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, addr := range host.Addrs() {
	//	fmt.Printf("Addr: %s/p2p/%s\n", addr, host.ID().Pretty())
	//}
	//
	//if err := ddht.Bootstrap(ctx); err != nil {
	//	panic(err)
	//}

	select {}
}
