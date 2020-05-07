Start the `relay` in a network that is accessible by any peer. 
It also acts as the bootstrap peer.

Next start the first `peer`:
`peer -room 1234 -bootstrap <address of relay>`

This will create a room with id `1234`.

Next join the room in another network:
`peer -join 1234 -bootstrap <address of relay`