# **MATDAAN** : Blockchain based E-Voting System

# TASK 1 Week
Digital Identity
Voter Authentication

## Solution
- SAML based Auth like Google, Facebook
- Biometric Face Auth
- MFA


---

<!-- Fabric v2.5.x is the current long-term support (LTS) release. -->

## **Download Hyperledger Fabric (Script)**
```sh
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
```

## **Download Components**
```sh
# b : binary | d : docker | s : samples
./install-fabric.sh # will download all (recommended)
./install-fabric.sh -h # usage help
```

# New Project
```sh
mkdir matdaan && cd matdaan
```

# Go Commands
```sh
go mod init e-voting # Initialize a new Go Module
go get -u github.com/hyperledger/fabric-contract-api-go # Install Package

# After adding or changing dependencies, it's always a good idea to Tidy Your Modules
# - Remove any unused dependencies from go.mod.
# - Download any missing dependencies.
# - Update go.sum.
go mod tidy

go build # Build an executable
```

# Deploy
```sh
network.sh deployCC -ccn e_voting -ccp ~/Hyperledger/matdaan/chaincode -ccv 1 -ccl go
```

## **Hyperledger Blockchain Explorer**
referrences :
    - [official](https://github.com/hyperledger-labs/blockchain-explorer?tab=readme-ov-file)
    - [medium](https://abhibvp003.medium.com/hyperledger-explorer-setup-with-hyperledger-fabric-c65f99749a03)

```sh
# prequisite node
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.2/install.sh | bash
nvm install 14.1
source ~/.profile

git clone https://github.com/hyperledger-labs/blockchain-explorer.git

# Install PostgreSQL
sudo apt-get install postgresql postgresql-contrib

# Configure & set username & password : postgres
cd blockchain-explorer/app
nano explorerconfig.json

# better use below
export DATABASE_HOST=127.0.0.1 \
export DATABASE_PORT=5432 \
export DATABASE_DATABASE=fabricexplorer \
export DATABASE_USERNAME=postgres \
export DATABASE_PASSWD=postgres

cd ../../
cd blockchain-explorer/app/persistence/fabric/postgreSQL/
chmod -R775 db/

cd db/
./createdb.sh



cd ~/blockchain-explorer/app/platform/fabric/connection-profile/
nano test-network.json
nano test-network-ca.json
# path > /home/ubuntu/Hyperledger/matdaan/hyperledger-fabric/fabric-samples/*

cd blockchain-explorer/app/persistence/fabric/postgreSQL/db
sudo -u postgres ./createdb.sh

sudo -u postgres psql -c '\l'
sudo -u postgres psql $DATABASE_DATABASE -c '\d'

# Do this below step
# sudo -u postgres psql
# [sudo] password for ubuntu: 
# psql (16.8 (Ubuntu 16.8-0ubuntu0.24.04.1))
# Type "help" for help.

# postgres=# ALTER USER postgres WITH PASSWORD 'postgres';
# ALTER ROLE
# postgres=# \q
# sudo service postgresql restart


# Required for Ubuntu
sudo apt-get install g++ build-essential

# root of the repository : explorer
./main.sh clean
./main.sh install
# might result into error ! Chill and relax

# Bootup Mode
nano ~/app/platform/fabric/config.json
    # "bootMode": "ALL", OR CUSTOM
    # "noOfBlocks": 0    OR 5 i.e. show latest 5 block

cd ~/blockchain-explorer/
./start.sh

# Go at http://localhost:8080/#/
username : exploreradmin
password : exploreradminpw

```


---
**Stuff**
```md
The fabric-contract-api provides the contract interface, a high level API for application developers to implement Smart Contracts. Within Hyperledger Fabric, Smart Contracts are also known as Chaincode. Working with this API provides a high level entry point to writing business logic.
For [Go](https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi)
Note that when using the contract api, each chaincode function that is called is passed a transaction context “ctx”, from which you can get the chaincode stub (GetStub() ), which has functions to access the ledger (e.g. GetState() ) and make requests to update the ledger (e.g. PutState() ).
```


```sh
# Commands

#   cd ~/Hyperledger/fabric-samples/test-network/ && ./network.sh up
#   peer lifecycle chaincode package voting.tar.gz --path ~/Learning/Projects/matdaan/chaincode --lang golang --label e_voting_1.0
#   peer lifecycle chaincode install voting.tar.gz

# ---------------------------------
# Start test network
cd ~/Hyperledger/fabric-samples/test-network/ && ./network.sh down && ./network.sh up createChannel

# ---------------------------------
# Package the smart contract
peer lifecycle chaincode package voting.tar.gz --path ~/Learning/Projects/matdaan/chaincode --lang golang --label e_voting_1.0

# ---------------------------------
# Install the chaincode package
source ./preqs/org1.sh
peer lifecycle chaincode install voting.tar.gz

source ./preqs/org2.sh
peer lifecycle chaincode install voting.tar.gz


# ---------------------------------
# find the package ID of a chaincode
peer lifecycle chaincode queryinstalled
export CC_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "Package ID:" | cut -d ' ' -f 3 | cut -d ':' -f 2 | tr -d ',')


# ---------------------------------
# Approve a chaincode definition¶

# Approve for Org 1
source ./preqs/org1.sh

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name e_voting --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# Approve for Org 2
source ./preqs/org2.sh
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name e_voting --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"



# ---------------------------------
# Committing the chaincode definition to the channel

# Ready?
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name e_voting --version 1.0 --sequence 1 --tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json

# Commit.
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name e_voting --version 1.0 --sequence 1 --tls --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# Verify!
peer lifecycle chaincode querycommitted --channelID mychannel --name e_voting --cafile "${MY_NETWORK}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"


# ---------------------------------
# ---------------------------------
# Testing ::
# Initialize Ledger:
peer chaincode invoke -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com \
-C ${CHANNEL_MAIN} -n e_voting --peerAddresses localhost:7051 \
--tlsRootCertFiles ${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
-c '{"function":"InitLedger","Args":[]}'



peer chaincode invoke -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com -C mychannel -n e_voting --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles ${MY_NETWORK}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

# Register Vote
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -C mychannel -n e_voting --peerAddresses localhost:7051 --tlsRootCertFiles ${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c '{"function":"RegisterVoter","Args":["VOTER001"]}'

# Cast Vote
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -C mychannel -n e_voting --peerAddresses localhost:7051 --tlsRootCertFiles ${MY_NETWORK}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c '{"function":"CastVote","Args":["VOTER001","CAND001"]}'

# Results
peer chaincode query -C mychannel -n e_voting -c '{"function":"GetResults","Args":[]}'
```