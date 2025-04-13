#!/bin/bash

# --------------------------------- Stop Test Network ---------------------------------------------------------------------------------------------------
export PROJ_PATH=~/Hyperledger/matdaan
export MY_NETWORK=${PROJ_PATH}/hyperledger-fabric/fabric-samples/test-network
${MY_NETWORK}/network.sh down
