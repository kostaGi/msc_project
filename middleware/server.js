const express = require('express');
const app = express();
const port = 2800;

"use strict";

var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};

Object.defineProperty(exports, "__esModule", { value: true });
const grpc = __importStar(require("@grpc/grpc-js"));
const fabric_gateway_1 = require("@hyperledger/fabric-gateway");
const crypto = __importStar(require("crypto"));
const fs_1 = require("fs");
const path = __importStar(require("path"));
const util_1 = require("util");
const utf8Decoder = new util_1.TextDecoder();
const channelName = envOrDefault('CHANNEL_NAME', 'mychannel');
const chaincodeName = envOrDefault('CHAINCODE_NAME', 'basic');
const mspId = envOrDefault('MSP_ID', 'Org1MSP');

// Path to crypto materials.
const cryptoPath = envOrDefault('CRYPTO_PATH', path.join('/home', 'kosix', 'Desktop', 'hyperledge_fabric' ,'fabric-samples', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com'));
// Path to user private key directory.
const keyDirectoryPath = envOrDefault('KEY_DIRECTORY_PATH', path.join(cryptoPath, 'users', 'User1@org1.example.com', 'msp', 'keystore'));
// Path to user certificate directory.
const certDirectoryPath = envOrDefault('CERT_DIRECTORY_PATH', path.join(cryptoPath, 'users', 'User1@org1.example.com', 'msp', 'signcerts'));
// Path to peer tls certificate.
const tlsCertPath = envOrDefault('TLS_CERT_PATH', path.join(cryptoPath, 'peers', 'peer0.org1.example.com', 'tls', 'ca.crt'));
// Gateway peer endpoint.
const peerEndpoint = envOrDefault('PEER_ENDPOINT', 'localhost:7051');
// Gateway peer SSL host name override.
const peerHostAlias = envOrDefault('PEER_HOST_ALIAS', 'peer0.org1.example.com');

//helper functions

async function newGrpcConnection() {
    const tlsRootCert = await fs_1.promises.readFile(tlsCertPath);
    const tlsCredentials = grpc.credentials.createSsl(tlsRootCert);
    return new grpc.Client(peerEndpoint, tlsCredentials, {
        'grpc.ssl_target_name_override': peerHostAlias,
    });
}
async function newIdentity() {
    const certPath = await getFirstDirFileName(certDirectoryPath);
    const credentials = await fs_1.promises.readFile(certPath);
    return { mspId, credentials };
}
async function getFirstDirFileName(dirPath) {
    const files = await fs_1.promises.readdir(dirPath);
    return path.join(dirPath, files[0]);
}
async function newSigner() {
    const keyPath = await getFirstDirFileName(keyDirectoryPath);
    const privateKeyPem = await fs_1.promises.readFile(keyPath);
    const privateKey = crypto.createPrivateKey(privateKeyPem);
    return fabric_gateway_1.signers.newPrivateKeySigner(privateKey);
}

function envOrDefault(key, defaultValue) {
    return process.env[key] || defaultValue;
}

async function displayInputParameters() {
    console.log(`channelName:       ${channelName}`);
    console.log(`chaincodeName:     ${chaincodeName}`);
    console.log(`mspId:             ${mspId}`);
    console.log(`cryptoPath:        ${cryptoPath}`);
    console.log(`keyDirectoryPath:  ${keyDirectoryPath}`);
    console.log(`certDirectoryPath: ${certDirectoryPath}`);
    console.log(`tlsCertPath:       ${tlsCertPath}`);
    console.log(`peerEndpoint:      ${peerEndpoint}`);
    console.log(`peerHostAlias:     ${peerHostAlias}`);
}

//verifyProof function
async function VerifyProof(contract,proof) {
    //console.log('*** Proof:', proof);
    console.log('\n--> Evaluate Transaction: verifyProof, function returns checking the proof');
    const resultBytes = await contract.evaluateTransaction('Verify', proof["public_key"], proof["proof"], proof["signature"]);
    const resultJson = utf8Decoder.decode(resultBytes);
    console.log('*** resultJson:', resultJson);
    const result = JSON.parse(resultJson);
    console.log('*** Result:', result);
    return result;
}

//verifyProof function
async function StorePublicKeys(contract,data) {
    console.log('\n--> StoreOwnerPublicKeyMiddleware');
    const resultBytes = await contract.submitTransaction('StoreOwnerPublicKeyMiddleware', data["owner"], data["public_key_D"], data["public_key_UE"]);
    const resultJson = utf8Decoder.decode(resultBytes);
    console.log('*** resultJson:', resultJson);
    const result = JSON.parse(resultJson);
    console.log('*** Result:', result);
    return result;
}

async function GetPublicKeyUE(contract,data) {
    console.log('\n--> ReadOwnerPublicKeyUE');
    const resultBytes = await contract.evaluateTransaction('ReadOwnerPublicKeyUE', data["id"], data["owner"]);
    const resultJson = utf8Decoder.decode(resultBytes);
    console.log('*** resultJson:', resultJson);
    const result = JSON.parse(resultJson);
    console.log('*** Result:', result);
    return result;
}

async function Init(contract,data) {
    console.log('\n--> InitLedger');
    const resultBytes = await contract.submitTransaction('InitLedger');
    const resultJson = utf8Decoder.decode(resultBytes);
    console.log('*** resultJson:', resultJson);
    const result = JSON.parse(resultJson);
    console.log('*** Result:', result);
    return result;
}


//main function
async function main(inputData) {
    await displayInputParameters();
    // The gRPC client connection should be shared by all Gateway connections to this endpoint.
    const client = await newGrpcConnection();
    const gateway = (0, fabric_gateway_1.connect)({
        client,
        identity: await newIdentity(),
        signer: await newSigner(),
        // Default timeouts for different gRPC calls
        evaluateOptions: () => {
            return { deadline: Date.now() + 5000 }; // 5 seconds
        },
        endorseOptions: () => {
            return { deadline: Date.now() + 15000 }; // 15 seconds
        },
        submitOptions: () => {
            return { deadline: Date.now() + 5000 }; // 5 seconds
        },
        commitStatusOptions: () => {
            return { deadline: Date.now() + 60000 }; // 1 minute
        },
    });
    try {
        // Get a network instance representing the channel where the smart contract is deployed.
        const network = gateway.getNetwork(channelName);
        // Get the smart contract from the network.
        const contract = network.getContract(chaincodeName);
        
        console.log('\n--> RequestType');
        console.log(inputData["RequestType"]);

        switch(String(inputData["RequestType"])){

            case "0":
                contract.evaluateTransaction('main');
                gateway.close();
                client.close();
                return 0
            case "1":
                response = await VerifyProof(contract,inputData);
                gateway.close();
                client.close();
                return response;
            
            case "2":
                response = await StorePublicKeys(contract,inputData);
                gateway.close();
                client.close();
                return response;
            
            case "3":
                response = await GetPublicKeyUE(contract,inputData);
                gateway.close();
                client.close();
                return response;
            
            case "4":
                response = await Init(contract,inputData);
                gateway.close();
                client.close();
                return response;


            default:
                gateway.close();
                client.close();
                console.log('\n--> Could not proccess request');

        }
        return null;
    } catch (error) {

        console.log ("error in try catch",error);
    }

}











// Middleware to parse JSON bodies
//app.use(express.json());
app.use(express.json({ limit: '10mb' }));

// Define the checkProof endpoint
app.post('/checkProof', async(req, res) => {

    console.log("api called");
    const data = req.body;

    const responseOnCheck =  await main(data)

    // Simple response to return the received data
    res.json({
        response:responseOnCheck
    });
});

// Start the server
app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});
