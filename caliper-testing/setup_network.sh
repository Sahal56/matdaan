# Follow this https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html

# Set Environment Variables# Set Environment Variables
export PROJ_PATH=~/Hyperledger/matdaan
export MY_NETWORK=${PROJ_PATH}/hyperledger-fabric/fabric-samples/test-network
export FABRIC_CFG_PATH=${PROJ_PATH}/hyperledger-fabric/fabric-samples/config
export CALIPER_PATH=${PROJ_PATH}/caliper-testing

export CHANNEL_MAIN=mychannel
export CHAINCODE_PATH=${PROJ_PATH}/chaincode
export CHAINCODE_PKG_NAME=evoting
export CHAINCODE_LANGUAGE=go

export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
export CORE_PEER_MSPCONFIGPATH=${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export ORDERER_CA_FILE=${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER_1_TLS_FILES="${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
export PEER_2_TLS_FILES="${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"



# Bring Up the Network without Certicate Authority (no -ca)
cd ${MY_NETWORK}
./network.sh down && ./network.sh up createChannel -c ${CHANNEL_MAIN}
# docker ps # ensure 3 peers are running and network is up

# Deploy Smart Contract before Testing !!!
${MY_NETWORK}/network.sh deployCC \
-ccn ${CHAINCODE_PKG_NAME} \
-ccp ${CHAINCODE_PATH} \
-ccv 1.0 \
-ccl ${CHAINCODE_LANGUAGE}

cd ${CALIPER_PATH}

# Invoking InitLedger to populate all candidates
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls true \
  --connTimeout 10s \
  --cafile ${ORDERER_CA_FILE} \
  -C ${CHANNEL_MAIN} \
  -n ${CHAINCODE_PKG_NAME} \
  --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER_1_TLS_FILES} \
  --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER_2_TLS_FILES} \
  -c '{"Args":["InitLedger"]}'
#   --tls \
