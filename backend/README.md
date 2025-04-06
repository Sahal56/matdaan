# Run Backend
```
go run main.go
```

```sh
source ./env_vars.sh
```


.env:
```sh
PEER_IP="dns:///localhost:7051"
PEER_HOSTNAME="peer0.org1.example.com"
TLS_CERT_PATH="~/Hyperledger/matdaan/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
MONGO_URI="mongodb://host.docker.internal:27017"
PORT="5000"
CHAINCODE_NAME=evoting
CHANNEL_NAME=mychannel
```


> MonggoDB server selection error: Firewall issue:
> Windows + R
`wf.msc`
port: `27017`
Name: `MongoDB WSL Access`