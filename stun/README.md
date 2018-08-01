


## get public address and nat type

run TestPubAddr


## connect to peer using udp

peer A and peer B, Connecting to each other using the other's public address.

run TestUdpConn

## connect to peer using tcp


peer A and peer B, Connecting to each other using the other's public address.


run TestTcpConn




little tips:
you can send udp or tcp packet by nc.

udp:
> nc -p <localPort> <remoteAddr> <remotePort> -u -w 5

tcp:

> nc -p <localPort> <remoteAddr> <remotePort> -w 5

-u: use udp
-w: connection wait time, unit is sec

more see: man nc

