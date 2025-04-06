#!/bin/bash

# Fabric
export PROJ_PATH=~/Hyperledger/matdaan
export FABRIC_CFG_PATH=${PROJ_PATH}/hyperledger-fabric/fabric-samples/config
export FABRIC_CA_CLIENT_HOME=${PROJ_PATH}/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/


# Mongo DB
export MONGO_URI="mongodb://host.docker.interna:27017"
export DB_NAME="evotingDB"

echo "Environment variables set"
