export MY_NETWORK=~/Hyperledger/fabric-samples/test-network

# Let's install the chaincode on the Org1 peer first.
# Set the following environment variables to operate the peer CLI as the Org1 admin user.
# The CORE_PEER_ADDRESS will be set to point to the Org1 peer, peer0.org1.example.com.

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
