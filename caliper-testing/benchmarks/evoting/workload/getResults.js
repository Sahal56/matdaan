'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class GetResultsWorkload extends WorkloadModuleBase {

    /**
     * Initializes the workload module instance.
     */
    constructor() {
        super();
        // this.startingKey = '';
        // this.endingKey = '';
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
        const request = {
            contractId: 'evoting',
            contractVersion: 'v1',
            contractFunction: 'GetResults',
            invokerIdentity: 'User1',
            contractArguments: [], // No arguments required for this function
            readOnly: true // Ensures this is treated as a query
        };

        try {
            const result = await this.sutAdapter.sendRequests(request);
            // console.log(`Transaction result: ${JSON.stringify(result)}`);
        } catch (error) {
            console.error(`Failed to submit transaction: ${error.message}`);
        }
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */

function createWorkloadModule() {
    return new GetResultsWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
