package main

import (
	"context"
	dht "github.com/libp2p/go-libp2p-kad-dht"

	//relay "github.com/libp2p/go-libp2p-circuit"
	//"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	"time"

	//"crypto/rand"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	//"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	//dht "github.com/libp2p/go-libp2p-kad-dht"

	//dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"

	//"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	//"github.com/libp2p/go-libp2p-core/routing"
	//"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	//dht "github.com/libp2p/go-libp2p-kad-dht"
	//"github.com/libp2p/go-libp2p-kad-dht/dual"
	secio "github.com/libp2p/go-libp2p-secio"
	"github.com/multiformats/go-multiaddr"
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
	room := flag.String("room", "", "the room you want to create")
	joinRoom := flag.String("join", "", "the room you want to join")
	bootstrap := flag.String("bootstrap", "/ip4/134.209.171.195/tcp/5000/p2p/QmWpBxWhq8G9G9m2yxc314Hfmd39PiHuWC5EJv3xZz9KxZ", "the relay")
	flag.Parse()

	ctx := context.Background()
	//defer cancel()

	if *joinRoom == "" && *room == "" {
		panic("At least one flag must be provided")
	}

	var ddht *dual.DHT
	var routingDiscovery *discovery.RoutingDiscovery
	routing := libp2p.Routing(func(host host.Host) (routing.PeerRouting, error) {
		var err error
		ddht, err = dual.New(ctx, host, dht.Mode(dht.ModeServer))
		routingDiscovery = discovery.NewRoutingDiscovery(ddht)

		return ddht, err
	})

	listenAddress := libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0")
	//myAddr, err := circuitRelay(*bootstrap, "/ip4/0.0.0.0/tcp/0")
	//if err != nil {
	//	panic(err)
	//}
	//listenAddress := libp2p.ListenAddrs(myAddr)

	enableRelay := libp2p.EnableRelay()

	security := libp2p.Security(secio.ID, secio.New)

	host, err := libp2p.New(ctx, listenAddress, security, routing, enableRelay, libp2p.NATPortMap())
	if err != nil {
		panic(err)
	}

	fmt.Println("Host Created")
	for _, addr := range host.Addrs() {
		fmt.Printf("Address: %s\n", addr)
	}

	fmt.Printf("My ID %s\n", host.ID().Pretty())

	// connect to the bootstrap peers
	ma, err := multiaddr.NewMultiaddr(*bootstrap)
	if err != nil {
		panic(err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		panic(err)
	}

	if err := host.Connect(ctx, *peerInfo); err != nil {
		panic(err)
	}
	fmt.Println("we are connected to the bootstrap peer")

	if err := ddht.Bootstrap(ctx); err != nil {
		panic(err)
	}

	fmt.Println("DHT in a bootstrapped state")

	time.Sleep(time.Second * 5)

	// now do chat specific stuff
	if *room != "" {
		host.SetStreamHandler(protocolKey(*room), handleStream)
		// create a room and wait for connections in it
		fmt.Printf("Advertising %s\n", protocolKey(*room))

		discovery.Advertise(ctx, routingDiscovery, string(protocolKey(*room)))

		fmt.Printf("Successfully advertised room: %s\n", *room)

		// wait forever
		select {}
	}

	if *joinRoom != "" {
		fmt.Printf("Joining room %s\n", protocolKey(*joinRoom))
		pctx, _ := context.WithTimeout(ctx, time.Second*10)
		peers, err := discovery.FindPeers(pctx, routingDiscovery, string(protocolKey(*joinRoom)))
		//peerChan, err := routingDiscovery.FindPeers(pctx, protocolKey(*joinRoom))
		if err != nil {
			panic(err)
		}

		if len(peers) == 0 {
			fmt.Println("No Peers Found")
			return
		}

		fmt.Printf("Found %d peers\n", len(peers))

		for _, p := range peers {
			fmt.Printf("Trying peer %s\n", p.ID.Pretty())
			fmt.Println("Addresses")
			for _, addr := range p.Addrs {
				ma, err := circuitRelay(addr.String(), p.ID.String())
				if err != nil {
					fmt.Printf("Error circuit Relay address %s\n", err)
					continue
				}
				info, err := peer.AddrInfoFromP2pAddr(ma)
				if err != nil {
					fmt.Printf("Error converting to addr info %s\n", err)
					continue
				}

				//fmt.Printf("Adding  the info to the peerstore")
				//host.Peerstore().AddAddrs(p.ID, p.Addrs, peerstore.TempAddrTTL)

				fmt.Printf("Going to try and connect via relay (%s)\n", ma)
				if err := host.Connect(ctx, *info); err == nil {
					fmt.Printf("We Have a connection try to create a stream now\n")
					stream, err := host.NewStream(ctx, p.ID, protocolKey(*joinRoom))
					if err != nil {
						fmt.Printf("Error dialing %s <%s>\n", p.ID.Pretty(), err)
					} else {
						go handleStream(stream)
						break
					}

				}

			}

		}

		select {}
	}

}

func protocolKey(key string) protocol.ID {
	return protocol.ID("/chat/" + key)
}

func handleStream(stream network.Stream) {
	fmt.Println("Handling Stream")
	defer stream.Close()
}

func circuitRelay(target string, id string) (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr("/p2p-circuit" + target + "/p2p/" + id)
}
