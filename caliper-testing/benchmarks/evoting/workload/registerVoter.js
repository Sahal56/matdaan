'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RegisterVoterWorkload extends WorkloadModuleBase {

    /**
     * Initializes the workload module instance.
     */

    constructor() {
        super();
        this.txIndex = 0;
        this.constituency = 'Anand';
    }

    /**
     * Initialize the workload module with the given parameters.
     * @param {number} workerIndex The 0-based index of the worker instantiating the workload module.
     * @param {number} totalWorkers The total number of workers participating in the round.
     * @param {number} roundIndex The 0-based index of the currently executing round.
     * @param {Object} roundArguments The user-provided arguments for the round from the benchmark configuration file.
     * @param {BlockchainInterface} sutAdapter The adapter of the underlying SUT.
     * @param {Object} sutContext The custom context object provided by the SUT adapter.
     * @async
     */
    
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        // this.startingKey = 'Client' + this.workerIndex + '_CAR' + this.roundArguments.startKey;
        // this.endingKey = 'Client' + this.workerIndex + '_CAR' + this.roundArguments.endKey;
        console.log(`Initializing workload module for worker ${workerIndex}`);
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */

    async submitTransaction() {
        this.txIndex++;
        const voterID =  `VOT${String(this.txIndex).padStart(3, '0')}`;
        // const voterID = `VOT${Math.floor(Math.random() * 10000)}`;
        // const constituency = 'Anand';

        const request = {
            contractId: 'evoting',
            contractFunction: 'RegisterVoter',
            invokerIdentity: 'User1',
            contractArguments: [voterID, this.constituency],
            readOnly: false // Since this is a write operation
        }

        try {
            const result = await this.sutAdapter.sendRequests(request);
            console.log(`Transaction result: ${JSON.stringify(result)}`);
        } catch (error) {
            console.error(`Failed to submit transaction: ${error.message}`);
        }
    }
}

function createWorkloadModule() {
    return new RegisterVoterWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
