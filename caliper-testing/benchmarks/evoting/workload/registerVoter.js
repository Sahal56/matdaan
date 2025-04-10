'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RegisterVoterWorkload extends WorkloadModuleBase {

    constructor() {
        super();
        this.txIndex = 0;
        this.constituency = 'Anand';
        this.voterIDs = [];
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        // Precompute voter IDs and assign range per worker
        const totalVoters = 100;
        const votersPerWorker = Math.ceil(totalVoters / totalWorkers);
        const start = workerIndex * votersPerWorker + 1;
        const end = Math.min((workerIndex + 1) * votersPerWorker, totalVoters);

        for (let i = start; i <= end; i++) {
            this.voterIDs.push(`VOT${String(i).padStart(4, '0')}`);
        }

        console.log(`[RegisterVoter][Worker ${workerIndex}] Assigned voter IDs from VOT${String(start).padStart(4, '0')} to VOT${String(end).padStart(4, '0')}`);


        // Optional: shuffle to simulate random ordering
        // this.voterIDs = this.voterIDs.sort(() => Math.random() - 0.5);
        // await new Promise(resolve => setTimeout(resolve, 2000));
    }

    async submitTransaction() {
        if (this.txIndex >= this.voterIDs.length) {
            return;
        }

        const voterID = this.voterIDs[this.txIndex];
        this.txIndex++;

        const request = {
            contractId: 'evoting',
            contractFunction: 'RegisterVoter',
            invokerIdentity: 'User1',
            contractArguments: [voterID, this.constituency],
            readOnly: false
        };

        try {
            await this.sutAdapter.sendRequests(request);
            // Commented out for speed: console.log(`Registered ${voterID}`);
        } catch (error) {
            console.error(`Failed to register ${voterID}: ${error.message}`);
        }
    }
}

function createWorkloadModule() {
    return new RegisterVoterWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
