# **Blockchain Explorer**

## It will be used to visualize Blockchain transactions, network, GUI, graphs, etc...

## **Connections**
> Go at http://localhost:8080/#/
username : `exploreradmin`
password : `exploreradminpw`

```sh
# Bring up the containers
docker-compose up
docker-compose up -d # detached terminal
```

---

<details>
   <summary> <b> Installation (Docker way) </b> </summary>

```sh
# Download below files | If you don't have wget => use curl
wget https://github.com/hyperledger-labs/blockchain-explorer/blob/main/.env
wget https://raw.githubusercontent.com/hyperledger/blockchain-explorer/main/examples/net1/config.json
wget https://raw.githubusercontent.com/hyperledger/blockchain-explorer/main/examples/net1/connection-profile/test-network.json -P connection-profile
wget https://raw.githubusercontent.com/hyperledger/blockchain-explorer/main/docker-compose.yaml

# Make the network up and running & Copy crypto files from our test network
export PROJ_PATH=${HOME}/Hyperledger/matdaan
source ${PROJ_PATH}/scripts/chaincode_install.sh
source ${PROJ_PATH}/scripts/org.sh
mkdir organizations
cp -r ${PROJ_PATH}/hyperledger-fabric/fabric-samples/test-network/organizations/ ${PROJ_PATH}/explorer/organizations

# Now open ${PROJ_PATH}/explorer/connection-profile/test-network.json
# >>> Replace : User@org1 => Admin@org1
```json
"organizations": {
		"Org1MSP": {
			"mspid": "Org1MSP",
			"adminPrivateKey": {
				"path": "/tmp/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
			},
			"peers": ["peer0.org1.example.com"],
			"signedCert": {
				"path": "/tmp/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
			}
		}
}
```

<details>

