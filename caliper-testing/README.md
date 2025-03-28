# **HyperLedger Caliper - Benchmark/Test - HyperLedger Fabric**

### It is used to benchmark our fabric network

<details>
   <summary> <b> Prequisite </b> </summary>

Versions:
- node 18.19
- peer/hyperledger 2.5.x

```sh
nvm install 18.19 #use nvm for managing node/npm
```

> NOTE:
>  - If you are using fabric 2.4 or higher then use below to bind with the new fabric node sdk which uses the new peer gateway service.
>  - npx caliper bind --caliper-bind-sut fabric:2.4


</details>

> for gloabal use/insatllation (-g)
```sh
user@ubuntu:~$ npm install -g --only=prod @hyperledger/caliper-cli@0.6.0
user@ubuntu:~$ caliper bind --caliper-bind-sut fabric:2.4 --caliper-bind-args=-g
user@ubuntu:~$ caliper launch manager \
    --caliper-workspace ~/caliper-benchmarks \
    --caliper-benchconfig benchmarks/scenario/simple/config.yaml \
    --caliper-networkconfig networks/fabric/test-network.yaml
```

```sh
# Clone Offical Caliper Benchmark repo for referrence
# I have kept that in .gitignore. so it wont appear in directory.
# But one should download it !!!
# git clone https://github.com/hyperledger/caliper-benchmarks

cd caliper-benchmarks

# we are installing locally | for global: use -g
npm install --only=prod @hyperledger/caliper-cli --omit=dev
npm install --only=prod @hyperledger/caliper-core --omit=dev
npx caliper bind --caliper-bind-sut fabric:2.4
npm ls # list local dir installed packages

source ./setup_network.sh

# Before Benchmarking. check & confirm, modify if needed
#   - benchmarks/evoting/*
#   - networks/test-network.yaml

# Start Caliper
npx caliper launch manager \
--caliper-workspace ./ \
--caliper-networkconfig networks/fabric/test-network.yaml \
--caliper-benchconfig benchmarks/evoting/config.yaml \
--caliper-flow-only-test \
--caliper-fabric-gateway-enabled

```