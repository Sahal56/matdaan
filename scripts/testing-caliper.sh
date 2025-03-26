
npm install --only=prod @hyperledger/caliper-cli

npx caliper bind --caliper-bind-sut fabric:2.4

# Below command work for levelDB
npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/scenario/simple/config.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled