# Set Environment Variables
export PACKAGE_NAME="e_voting"
export PACKAGE_ZIP=${PACKAGE_NAME}.tar.gz
export PACKAGE_LANGUAGE="golang"
export CHANNEL_MAIN=mychannel
export HLF_PATH=~/Hyperledger

export PROJ_MATDAAN=${HLF_PATH}/Projects/matdaan

export MY_NETWORK=~/Hyperledger/fabric-samples/test-network
export FABRIC_CFG_PATH=~/Hyperledger/fabric-samples/config
# export FABRIC_CFG_PATH=/Users/sahal/Hyperledger/fabric-samples/test-network/organizations/cryptogen/

# --------------------------------- Building Go Module | Smart Contract ---------------------------------------------------------------------------------------------------
# Note: This is optional. Comment it out if already Go chaincode is build into e_voting.tar.gz
go build -C ${PROJ_MATDAAN}/chaincode/ -o ${PROJ_MATDAAN}/chaincode/

# --------------------------------- Start Test Network ---------------------------------------------------------------------------------------------------
${HLF_PATH}/fabric-samples/test-network/network.sh down
${HLF_PATH}/fabric-samples/test-network/network.sh up createChannel -c ${CHANNEL_MAIN} -ca 


# --------------------------------- Package the smart contract ---------------------------------------------------------------------------------------------------
peer lifecycle chaincode package ${PACKAGE_ZIP} --path ${PROJ_MATDAAN}/chaincode \
--lang ${PACKAGE_LANGUAGE} --label ${PACKAGE_NAME}

# --------------------------------- Install the chaincode package ------------------------------------------------------------------------------------------------
source ${PROJ_MATDAAN}/scripts/org1.sh
peer lifecycle chaincode install ${PACKAGE_ZIP}

source ${PROJ_MATDAAN}/scripts/org2.sh
peer lifecycle chaincode install ${PACKAGE_ZIP}


# --------------------------------- Get Package Chaincode's ID ---------------------------------------------------------------------------------------------------
peer lifecycle chaincode queryinstalled
export CC_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "Package ID:" | cut -d ' ' -f 3 | cut -d ':' -f 2 | tr -d ',')


# --------------------------------- Approve Chaincode ------------------------------------------------------------------------------------------------------------
# Approve for Org 1
source ${PROJ_MATDAAN}/scripts/org1.sh
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--channelID ${CHANNEL_MAIN} --name ${PACKAGE_NAME} --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 \
--tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# Approve for Org 2
source ${PROJ_MATDAAN}/scripts/org2.sh
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--channelID ${CHANNEL_MAIN} --name ${PACKAGE_NAME} --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 \
--tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# --------------------------------- Committing the chaincode definition to the channel -----------------------------------------------------------------------------
# Ready?
peer lifecycle chaincode checkcommitreadiness --channelID ${CHANNEL_MAIN} --name ${PACKAGE_NAME} --version 1.0 --sequence 1 \
--tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json

# Commit.
peer lifecycle chaincode commit \
-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--channelID ${CHANNEL_MAIN} --name ${PACKAGE_NAME} --version 1.0 --sequence 1 \
--tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
--peerAddresses localhost:7051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# Verify!
peer lifecycle chaincode querycommitted --channelID ${CHANNEL_MAIN} --name ${PACKAGE_NAME} \
--cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"



# Testing Invoke
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls true \
  --cafile ${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C ${CHANNEL_MAIN} -n ${PACKAGE_NAME} \
  --peerAddresses localhost:7051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"Args":["InitLedger"]}'
