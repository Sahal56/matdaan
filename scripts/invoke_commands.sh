
# --------------------------------- Testing -----------------------------------------------------------------------------
# Important: Network with CAuthority, any one Peer's Context (Org1.sh), --connTimeout fixed

# Always we have to use context of Peer 1 or 2
source ${PROJ_PATH}/scripts/org1.sh

# InitLedger
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls true \
  --connTimeout 10s \
  --cafile ${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C ${CHANNEL_MAIN} -n ${CHAINCODE_PKG_NAME} \
  --peerAddresses localhost:7051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"Args":["InitLedger"]}'


# RegisterVoter
...
-c '{"Args":["RegisterVoter","VOT001", "Anand"]}'

# CastVote
...
-c '{"Args":["CastVote", "VOT001", "CAND001"]}'

# GetResults
...
  -c '{"Args":["GetResults"]}'
