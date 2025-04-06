#!/bin/bash

# --------------------------------- Testing -----------------------------------------------------------------------------
# Important: Network with CAuthority, any one Peer's Context (Org1.sh), --connTimeout fixed

# Always we have to use context of Peer 1 or 2
export PROJ_PATH=~/Hyperledger/matdaan
export FABRIC_CFG_PATH=${PROJ_PATH}/hyperledger-fabric/fabric-samples/config
export CHAINCODE_PKG_NAME="evoting"
export CHANNEL_MAIN="mychannel"

export CA_FILE=${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER_1_TLS_FILES="${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
export PEER_2_TLS_FILES="${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"



echo "1 : InitLedger"
echo "2 : RegisterVoter"
echo "3 : CastVote"
echo "4 : CastVoteByCandidate"
echo "5 : GetResults"
echo "6 : GetCandidatesByVoter"

read -p "Input: " input
# votID="VOT0001"

case $input in
    1)
        QUERY='{"Args":["InitLedger"]}'
        ;;
    2)
        read -p "VoterID: " votID 
        QUERY='{"Args":["RegisterVoter","'"$votID"'","Anand"]}'
        ;;
    3)
        read -p "VoterID: " votID
        QUERY='{"Args":["CastVote","'"$votID"'","CAND0001"]}'
        ;;
    4)
        QUERY='{"Args":["CastVoteByCandidate","CAND0001","CAND0001"]}'
        ;;
    5)
        QUERY='{"Args":["GetResults"]}'
        ;;
    6)
        QUERY='{"Args":["GetResultByVoterID", "VOT0001"]}'
        ;;
    *)
        echo "Invalid input. Please enter a number between 1 and 5."
        exit 1
        ;;
esac

# InitLedger
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls true \
  --connTimeout 10s \
  --cafile ${CA_FILE} \
  -C ${CHANNEL_MAIN} \
  -n ${CHAINCODE_PKG_NAME} \
  --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER_1_TLS_FILES} \
  --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER_2_TLS_FILES} \
  -c "$QUERY"
