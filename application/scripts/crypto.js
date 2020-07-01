const crypto = require('crypto')
const fs = require('fs')
const p = 'ab59a697ad185fe2dc924430d2782c6ed9b7338445e6a1296320e4ad6232373d5c866975e42e7037488f096d23b3a72f1ef168740698fbdbaadcf50b30ed7c58596c2641d5aabd6bb6173ce07e742c9cd3d910494117a215bbb1d36ebf9cc7a2b736ebff32eb873e4ff7820ecc21f00b44e8294b11a73cdc960e3e9c2a965833'
const g = '02'
var dh = crypto.createDiffieHellman(p, g);
const algorithm = 'aes-256-cbc'

function generateKeys(x) {
    publicKey = dh.generateKeys()
    privateKey = dh.getPrivateKey()
    fs.writeFileSync('./priv' + x, privateKey);

    return publicKey
}

function decrypt(cryptogram, publicKey, x) {
    privateKey = fs.readFileSync('./priv' + x,'base64')
    dh.setPrivateKey(privateKey, 'base64')
    sharedKey = dh.computeSecret(publicKey)
    sharedKeyHash = crypto.createHash('sha256').update(sharedKey).digest()
    const iv = Buffer.alloc(16, 0); 
    decipher = crypto.createDecipheriv(algorithm, sharedKeyHash, iv);
    let plaintext = decipher.update(cryptogram,'binary','utf8');
    plaintext += decipher.final('utf8');

    return plaintext
}

function encrypt(plaintext, publicKey, x) {
    privateKey = fs.readFileSync('./priv' + x,'base64')
    dh.setPrivateKey(privateKey, 'base64')
    sharedKey = dh.computeSecret(publicKey)
    sharedKeyHash = crypto.createHash('sha256').update(sharedKey).digest()
    const iv = Buffer.alloc(16, 0); 
    cipher = crypto.createCipheriv(algorithm, sharedKeyHash, iv);
    let criptogram = cipher.update(plaintext, 'utf8', 'binary');
    criptogram += cipher.final('binary');

    return criptogram
}

module.exports.generateKeys = generateKeys
module.exports.decrypt = decrypt
module.exports.encrypt = encrypt


function main(){
    pk1 = generateKeys('1')
    dh = crypto.createDiffieHellman(p, g);
    pk2 = generateKeys('2')

    c = encrypt("ola", pk2, '1')
    console.log(decrypt(c, pk1, '2'))
}

main();