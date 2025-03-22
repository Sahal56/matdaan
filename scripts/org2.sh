export MY_NETWORK=~/Hyperledger/fabric-samples/test-network

# We can now install the chaincode on the Org2 peer.
# Set the following environment variables to operate as the Org2 admin 
# and target the Org2 peer, peer0.org2.example.com.

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=peer0.org2.example.com:9051