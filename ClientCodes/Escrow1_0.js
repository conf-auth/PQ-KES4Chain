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
const { buildCCPOrg0, buildWallet } = require('/home/user/go/src/github.com/hyperledger/fabric/scripts/fabric-samples/test-application/javascript/AppUtil.js');



const channelName = 'mychannel';
const chaincodeName = 'Escrow';
const mspOrg0 = 'Org0MSP';

const walletPath = path.join(__dirname, 'wallet0');
const userId = 'admin';

const kem_algorithm_name = "FrodoKEM-640-AES";
/*
Enabled KEMs:
[DEFAULT BIKE1-L1-CPA BIKE1-L3-CPA BIKE1-L1-FO BIKE1-L3-FO 
Classic-McEliece-348864 Classic-McEliece-348864f Classic-McEliece-460896 
Classic-McEliece-460896f Classic-McEliece-6688128 Classic-McEliece-6688128f 
Classic-McEliece-6960119 Classic-McEliece-6960119f Classic-McEliece-8192128 Classic-McEliece-8192128f 
HQC-128 HQC-192 HQC-256 
Kyber512 Kyber768 Kyber1024 Kyber512-90s Kyber768-90s Kyber1024-90s 
NTRU-HPS-2048-509 NTRU-HPS-2048-677 NTRU-HPS-4096-821 NTRU-HRSS-701 
ntrulpr653 ntrulpr761 ntrulpr857 sntrup653 sntrup761 sntrup857 
LightSaber-KEM Saber-KEM FireSaber-KEM 
FrodoKEM-640-AES FrodoKEM-640-SHAKE FrodoKEM-976-AES FrodoKEM-976-SHAKE FrodoKEM-1344-AES FrodoKEM-1344-SHAKE 
SIDH-p434 SIDH-p503 SIDH-p610 SIDH-p751 SIDH-p434-compressed SIDH-p503-compressed SIDH-p610-compressed SIDH-p751-compressed 
SIKE-p434 SIKE-p503 SIKE-p610 SIKE-p751 SIKE-p434-compressed SIKE-p503-compressed SIKE-p610-compressed SIKE-p751-compressed]
*/

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
		
		const ccp = buildCCPOrg0();
		const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org0.example.com');
		const wallet = await buildWallet(Wallets, walletPath);	
		await enrollAdmin(caClient, wallet, mspOrg0);
		await registerAndEnrollUser(caClient, wallet, mspOrg0, userId, 'org2.department1');

		const gateway = new Gateway();
		try {
			
			await gateway.connect(ccp, {
				wallet,
				identity: userId,
				discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
			});
			
			console.log('\n---------------------');
			const network = await gateway.getNetwork(channelName);
			const contract = network.getContract(chaincodeName);
			console.log('\n---------------------');
		
			console.log("***********************************************************************")
			console.log('\n--> Submit Transaction: genKey');
			let prikey1 = await contract.submitTransaction("Gen_EA_KeyPair","EsPublicKey01", kem_algorithm_name);

			console.log(prikey1.toString());
			fs.writeFileSync("prikey1.txt",prikey1);
			console.log('*** prikey1: committed');

			
		} finally {
			gateway.disconnect();
		}
	} catch (error) {
		console.error(`******** FAILED to run the application: ${error}`);
	}

	console.log('*** application ending');

}

main();
