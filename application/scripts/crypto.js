const crypto = require('crypto')
const fs = require('fs')
const p = '4G7o9nf+amsT3pwhHipZxDiZ9n3/pm60aQEcgLELJNUetyst7+A8QkJV/x+9pql1FyccE/fE5eE0g215fZxJib6VvAPKjrYNJqjFijkbGkweCyb/C3diZU6GL7Pjgq+qF9MR7591hl4yhjwnQWe6DLchOHKnjGJ9sMbxY1xRHb08s0epre1xdFwWD2rusNxvTJYIQBAGDHI2JMg5Z7lza/IXava63KWYxJX5SoZGCvi6qPe15FMBF5kA7Wsfcnmi6LadL8PaVb+B6xL8gShf4axdQQNhuRxJ8hesrv9+jwy0BRcggFNRBMKHPYYzxEDcvcI6y7G+jQpUbP8dQQWyOw=='
const g = 'Ag=='
const dh = crypto.createDiffieHellman(p, 'base64', g, 'base64');
const algorithm = 'aes-256-cbc'

function generateKeys() {
    publicKey = dh.generateKeys('base64')
    privateKey = dh.getPrivateKey('base64')
    fs.writeFileSync('./priv', privateKey);

    return publicKey
}

function decrypt(cryptogram, publicKey) {
    privateKey = fs.readFileSync('./priv','base64')
    dh.setPrivateKey(privateKey, 'base64')
    sharedKey = dh.computeSecret(publicKey,'base64','utf8')
    sharedKeyHash = crypto.createHash('sha256').update(sharedKey).digest()
    const iv = Buffer.alloc(16, 0); 
    decipher = crypto.createDecipheriv(algorithm, sharedKeyHash, iv);
    let plaintext = decipher.update(cryptogram,'binary','utf8');
    plaintext += decipher.final('utf8');

    return plaintext
}

function encrypt(plaintext, publicKey) {
    privateKey = fs.readFileSync('./priv','base64')
    dh.setPrivateKey(privateKey, 'base64')
    sharedKey = dh.computeSecret(publicKey,'base64','utf8')
    sharedKeyHash = crypto.createHash('sha256').update(sharedKey).digest()
    const iv = Buffer.alloc(16, 0); 
    cipher = crypto.createCipheriv(algorithm, sharedKeyHash, iv);
    let criptogram = cipher.update(plaintext, 'utf8', 'binary');
    criptogram += cipher.final('binary');

    return criptogram
}

function main() {
    pk = generateKeys()
    let c = encrypt(fs.readFileSync('./invoke.js','utf8'), pk)
    console.log(c)
    let p = decrypt(c, pk)
    console.log(p)
}


main()
module.exports.generateKeys = generateKeys
module.exports.decrypt = decrypt
module.exports.encrypt = encrypt