#!/bin/bash

# Set Environment Variables
export PROJ_PATH=~/Hyperledger/matdaan
export MY_NETWORK=${PROJ_PATH}/hyperledger-fabric/fabric-samples/test-network
export FABRIC_CFG_PATH=${PROJ_PATH}/hyperledger-fabric/fabric-samples/config

export CHAINCODE_PATH=${PROJ_PATH}/chaincode
export CHAINCODE_PKG_NAME="evoting"
export CHAINCODE_PKG_ZIP=${CHAINCODE_PKG_NAME}.tar.gz
export CHAINCODE_LANGUAGE="go"
export CHANNEL_MAIN="mychannel"


# --------------------------------- Building Go Module | Smart Contract ---------------------------------------------------------------------------------------------------
# Note: This is optional. Comment it out if already Go chaincode is build into e_voting.tar.gz
# go build -C ${PROJ_PATH}/chaincode/ -o ${PROJ_PATH}/chaincode/

# --------------------------------- Start Test Network ---------------------------------------------------------------------------------------------------
${MY_NETWORK}/network.sh down
echo "1 : With Certificate Authority"
echo "2 : Without CA | CryptoGen?"

read -p "Input: " input

case $input in
    1)
        ${MY_NETWORK}/network.sh up createChannel -ca -c ${CHANNEL_MAIN}
        ;;
    2)
        ${MY_NETWORK}/network.sh up createChannel -c ${CHANNEL_MAIN}
        ;;
    *)
        echo "Invalid input. Please enter a number between 1 or 2."
        exit 1
        ;;
esac
#  To add Certificate Authority add above: -ca


#                              Automatically
# Install & Deploy Smart Contract
${MY_NETWORK}/network.sh deployCC -ccn ${CHAINCODE_PKG_NAME} -ccp ${CHAINCODE_PATH} -ccv 1.0 -ccl ${CHAINCODE_LANGUAGE}




#                      Manually (OLD)
# --------------------------------- Package the smart contract ---------------------------------------------------------------------------------------------------
# peer lifecycle chaincode package ${CHAINCODE_PKG_ZIP} --path  ${CHAINCODE_PATH} \
# --lang ${CHAINCODE_LANGUAGE} --label ${CHAINCODE_PKG_NAME}

# --------------------------------- Install the chaincode package ------------------------------------------------------------------------------------------------
# source ${PROJ_PATH}/scripts/org1.sh
# peer lifecycle chaincode install ${CHAINCODE_PKG_ZIP}

# source ${PROJ_PATH}/scripts/org2.sh
# peer lifecycle chaincode install ${CHAINCODE_PKG_ZIP}

# --------------------------------- Get Package Chaincode's ID ---------------------------------------------------------------------------------------------------
# peer lifecycle chaincode queryinstalled
# export CC_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "Package ID:" | cut -d ' ' -f 3 | cut -d ':' -f 2 | tr -d ',')


# --------------------------------- Approve Chaincode ------------------------------------------------------------------------------------------------------------
# Approve for Org 1
# source ${PROJ_PATH}/scripts/org1.sh
# peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
# --channelID ${CHANNEL_MAIN} --name ${CHAINCODE_PKG_NAME} --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 \
# --tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# Approve for Org 2
# source ${PROJ_PATH}/scripts/org2.sh
# peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
# --channelID ${CHANNEL_MAIN} --name ${CHAINCODE_PKG_NAME} --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 \
# --tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# --------------------------------- Committing the chaincode definition to the channel -----------------------------------------------------------------------------
# Ready?
# peer lifecycle chaincode checkcommitreadiness --channelID ${CHANNEL_MAIN} --name ${CHAINCODE_PKG_NAME} --version 1.0 --sequence 1 \
# --tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json

# Commit.
# peer lifecycle chaincode commit \
# -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
# --channelID ${CHANNEL_MAIN} --name ${CHAINCODE_PKG_NAME} --version 1.0 --sequence 1 \
# --tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
# --peerAddresses localhost:7051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
# --peerAddresses localhost:9051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# Verify!
# peer lifecycle chaincode querycommitted --channelID ${CHANNEL_MAIN} --name ${CHAINCODE_PKG_NAME} \
# --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

