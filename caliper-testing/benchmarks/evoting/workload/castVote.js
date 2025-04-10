'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class CastVoteWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
        this.candidateID = 'CAND0001';
        this.voterIDs = [];
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        console.log(`Initializing workload module for worker ${workerIndex} of ${totalWorkers}`);

        // Optional wait time to ensure registration completes | 20s delay
        await new Promise(resolve => setTimeout(resolve, 10000));

        const totalVoters = 100;
        const votersPerWorker = Math.floor(totalVoters / totalWorkers);
        const remainder = totalVoters % totalWorkers;

        // Account for uneven division
        let start = workerIndex * votersPerWorker + Math.min(workerIndex, remainder) + 1;
        let end = start + votersPerWorker - 1;
        if (workerIndex < remainder) {
            end += 1;
        }

        for (let i = start; i <= end; i++) {
            this.voterIDs.push(`VOT${String(i).padStart(4, '0')}`);
        }

        console.log(`[CastVote][Worker ${workerIndex}] Assigned ${this.voterIDs.length} voters: ${this.voterIDs.join(', ')}`);
    }

    async submitTransaction() {
        // await new Promise(resolve => setTimeout(resolve, 10));  // small backoff

        if (this.txIndex >= this.voterIDs.length) {
            return;
        }

        const voterID = this.voterIDs[this.txIndex];
        this.txIndex++;

        const request = {
            contractId: 'evoting',
            contractFunction: 'CastVote',
            invokerIdentity: 'User1',
            contractArguments: [voterID, this.candidateID],
            readOnly: false
        };

        try {
            await this.sutAdapter.sendRequests(request);
        } catch (error) {
            console.error(`Failed to cast vote for ${voterID}: ${error.message}`);
        }
    }
}

function createWorkloadModule() {
    return new CastVoteWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
