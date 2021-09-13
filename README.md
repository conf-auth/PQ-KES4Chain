# PQ-KES4Chain

About The Function of Chaincodes
----------
The main implementation of PQ-KES4Chain consists of three chaincodes, namely Sender, Receiver, and Escrow.<br>
<br>
The `Sender chaincode` provides API to generate/upload encrypted secret data and escrowed shared secret (i.e., sender data) and API to retrieve the sender data for other chaincodes. <br>
<br>
The `Receiver chaincode` only provides API to decrypt the encrypted secret data.<br>
<br>
The `Escrow chaincode` can be invoked to generate the escrow agent's public/private key pair, get the related public key, use the private key to decapsulate/recover the shared secret, and utilize the shared secret to decrypt the secret data.<br>

About The Function of Client Codes
----------
To further help developers to create their post-quantum supervised secret data sharing applications, we also provide client codes showing how to invoke the chaincodes. <br>
<br>
The `Sender.js` invokes  `Sender chaincode APIs` to generate and upload of encrypted secret data and escrowed shared secret.<br>
<br>
The `Receiver.js` invokes  `Receiver chaincode APIs` to download and decrypt the secret data.<br>
<br>
The `Escrow1_0.js` and `Escrow2_0.js` invoke  `Escrow chaincode APIs` to generate the post-quantum public/private key pair for escrow agent 1/2.<br>
<br>
The `Escrow1_1.js` and `Escrow2_1.js` invoke  `Escrow chaincode APIs` to decapsulate (recover) the shared secret escrowed to escrow agent 1/2.<br>
<br>
The `Supervisor.js` invokes `Escrow chaincode APIs` to decrypt the secret data using the decapsulated shared secret.<br>
<br>
The `CAUtil.js` and `AppUtil` are invoked by the other client codes to get the Fabric CA certificates and setup information of our system.<br>
<br>
`Notes:` In `Sender.js`,`Escrow1_0.js`, `Escrow1_1.js`,`Escrow2_0.js` and `Escrow2_1.js` codes, we fix the kem algorithm name as shown below. Developers can choose one of the following `Enabled KEMs` for PQ-KES4Chain operation.
```
const kem_algorithm_name = "FrodoKEM-640-AES";
```
```
Enabled KEMs:
DEFAULT BIKE1-L1-CPA BIKE1-L3-CPA BIKE1-L1-FO BIKE1-L3-FO 
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
SIKE-p434 SIKE-p503 SIKE-p610 SIKE-p751 SIKE-p434-compressed SIKE-p503-compressed SIKE-p610-compressed SIKE-p751-compressed

```

About The Function of Command Lines
----------
If developers do not want to use client codes to create post-quantum supervised secret data sharing applications, they can also use the command lines to complete the post-quantum supervised secret data sharing operations.


About The Use of Docker Image
----------
In PQ-KES4Chain, we utilize the liboqs 0.4.0 library [1] together with its Go wrapper [2] to generate public/private key pair, encapsulate and decapsulate the shared secret for all the related post-quantum KEM algorithms. Since the liboqs library is incompatible with the native Hyperledger Fabric docker image for chaincode execution, we build a new docker image integrated with liboqs based on Ubuntu 18.04 and use it as the execution environment of our chaincodes. The versions of docker and docker-compose we use are 20.10.2 and 1.26.2. The tutorial of how to build and use the new docker image is given as follows. <br>
<br>
1.In ccenv dockerfile, use the `Ubuntu18.04` as the docker image environment, download and install `liboqs library` required for the experiment, install and configure the `Go env(go1.15.7.linux-amd64.tar.gz)`. <br>
<br>
2.Build the dockerfile.
```
docker built -t ccenv:latest ./ 
```
3.Modify core.yaml.<br>
  Set the `chaincode.builder` and `chaincode.golang.runtime` as `ccenv:latest`, and set the `chaincode.golang.dynamicLink` as `true`. The part of the code that needs to be modified in core.yaml is set as follows:

```
###############################################################################
#
#    Chaincode section
#
###############################################################################
chaincode:

    # The id is used by the Chaincode stub to register the executing Chaincode
    # ID with the Peer and is generally supplied through ENV variables
    # the `path` form of ID is provided when installing the chaincode.
    # The `name` is used for all other requests and can be any string.
    id:
        path:
        name:

    # Generic builder environment, suitable for most chaincode types
    builder: ccenv:latest

    # Enables/disables force pulling of the base docker images (listed below)
    # during user chaincode instantiation.
    # Useful when using moving image tags (such as :latest)
    pull: false

    golang:
        # golang will never need more than baseos
        runtime: ccenv:latest
        # whether or not golang chaincode should be linked dynamically
        dynamicLink: true
```
4.Start the fabric network and deploy chaincodes.

About The Experiment Data
----------
Firstly, we test the execution time of every step in PQ-KES4Chain (except the second step because different developer may store their secret data off-chain and read the data in different ways) based on all the different post-quantum KEM algorithms in the NIST call. To provide enough security level, the length of pre-negotiated AES session key SK between sender/receiver peers is set to 256 bits, whereas the length of the secret data M is set to the input block length (i.e., 16 bytes) of AES algorithm.

We record all the execution time available in the experiment data, and highlight the consumed time of one key pair generation, one encapsulation operation and one decapsulation operation of each post-quantum KEM algorithm in boldface.

A full list of the on-chain storage space sizes can also be found in the experiment data, and we highlight the needed on-chain space of two ciphertexts, two public keys and two shared secrets of each post-quantum KEM algorithm in boldface.

References
----------
[1] Liboqs, https://github.com/open-quantum-safe/liboqs

[2] Liboqs-go, https://github.com/open-quantum-safe/liboqs-go


