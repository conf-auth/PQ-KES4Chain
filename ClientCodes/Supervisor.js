/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */		

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const fs = require('fs');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('/home/user/go/src/github.com/hyperledger/fabric/scripts/fabric-samples/test-application/javascript/CAUtil.js');
const { buildCCPOrg3, buildWallet } = require('/home/user/go/src/github.com/hyperledger/fabric/scripts/fabric-samples/test-application/javascript/AppUtil.js');


const channelName = 'mychannel';
const chaincodeName = 'Escrow';
const mspOrg3 = 'Org3MSP';


const walletPath = path.join(__dirname, 'wallet3');
const userId = 'admin';


function prettyJSONString(inputString) {
	return JSON.stringify(JSON.parse(inputString), null, 2);
}


async function main() {
	let skipInit = false;
	if (process.argv.length > 2) {
		if (process.argv[2] === 'skipInit') {
			skipInit = true;
		}
	}


	try {
		
		const ccp = buildCCPOrg3();
		const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org3.example.com');
		const wallet = await buildWallet(Wallets, walletPath);	
		await enrollAdmin(caClient, wallet, mspOrg3);
		await registerAndEnrollUser(caClient, wallet, mspOrg3, userId, 'org0.department1');

		const gateway = new Gateway();
		try {
			await gateway.connect(ccp, {
				wallet,
				identity: userId,
				discovery: { enabled: true, asLocalhost: false } // using asLocalhost as this gateway is using a fabric network deployed locally
			});

			const network = await gateway.getNetwork(channelName);
			const contract = network.getContract(chaincodeName);
			console.log('\n--> Evaluate Transaction: Dec_Sec_Data');
			let result = await contract.evaluateTransaction('Dec_Sec_Data',  "send001", "EsSharedSecret01", "EsSharedSecret02", "PrivateKey1Collection", "PrivateKey2Collection");
			console.log(result.toString());

			
		} finally {
			
			gateway.disconnect();
		}
	} catch (error) {
		console.error(`******** FAILED to run the application: ${error}`);
	}

	console.log('*** application ending');

}

main();
