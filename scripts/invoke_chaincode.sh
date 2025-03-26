#!/bin/bash
# --------------------------------- Testing -----------------------------------------------------------------------------
# Important: Network with CAuthority, any one Peer's Context (Org1.sh), --connTimeout fixed

# Always we have to use context of Peer 1 or 2
export PROJ_PATH=~/Hyperledger/matdaan
source ${PROJ_PATH}/scripts/org1.sh
export CHAINCODE_PKG_NAME="evoting"

export CA_FILE=${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER_1_TLS_FILES="${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
export PEER_2_TLS_FILES="${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"



echo "1 : InitLedger"
echo "2 : RegisterVoter"
echo "3 : CastVote"
echo "4 : CastVoteByCandidate"
echo "5 : GetResults"

read -p "Input: " input

case $input in
    1)
        QUERY='{"Args":["InitLedger"]}'
        ;;
    2)
        QUERY='{"Args":["RegisterVoter","VOT001","Anand"]}'
        ;;
    3)
        QUERY='{"Args":["CastVote","VOT001","CAND001"]}'
        ;;
    4)
        QUERY='{"Args":["CastVoteByCandidate","CAND001","CAND001"]}'
        ;;
    5)
        QUERY='{"Args":["GetResults"]}'
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
  -c ${QUERY}


# # RegisterVoter
# ...
#   -c '{"Args":["RegisterVoter","VOT001", "Anand"]}'

# # CastVoteByCandidate
# ...
#   -c '{"Args":["CastVoteByCandidate", "CAND001", "CAND001"]}'

# # CastVote
# ...
#   -c '{"Args":["CastVote", "VOT001", "CAND001"]}'



# # GetResults
# ...
#   -c '{"Args":["GetResults"]}'
