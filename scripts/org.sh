#!/bin/bash
# --------------------------------- Testing -----------------------------------------------------------------------------
# Important: Network with CAuthority, any one Peer's Context (Org1.sh), --connTimeout fixed

# Always we have to use context of Peer 1 or 2
export PROJ_PATH=~/Hyperledger/matdaan
export MY_NETWORK=${PROJ_PATH}/hyperledger-fabric/fabric-samples/test-network
export FABRIC_CFG_PATH=${PROJ_PATH}/hyperledger-fabric/fabric-samples/config
source ${PROJ_PATH}/scripts/org1.sh
export CHAINCODE_PKG_NAME="evoting"


echo "1 : Org1 || 2 : Org2"

read -r -p "Input: " org_choice
case $org_choice in
    1)
        echo "Switching to Org1..."
        export CORE_PEER_TLS_ENABLED=true
        export CORE_PEER_LOCALMSPID="Org1MSP"
        export CORE_PEER_TLS_ROOTCERT_FILE=${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
        export CORE_PEER_MSPCONFIGPATH=${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
        export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
        ;;
    2)
        echo "Switching to Org2..."
        export CORE_PEER_TLS_ENABLED=true
        export CORE_PEER_LOCALMSPID="Org2MSP"
        export CORE_PEER_TLS_ROOTCERT_FILE=${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
        export CORE_PEER_MSPCONFIGPATH=${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
        export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
        ;;
    *)
        echo "Invalid input. Please select 1 (Org1) or 2 (Org2)."
        exit 1
        ;;
esac