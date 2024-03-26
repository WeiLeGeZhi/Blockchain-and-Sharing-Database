package data;

import utils.Base58Util;
import utils.SecurityUtil;

import java.security.KeyPair;
import java.security.PublicKey;
import java.security.PrivateKey;

public class Account {

    private final PublicKey publicKey;
    private final PrivateKey privateKey;

    public Account(){
        KeyPair keyPair = SecurityUtil.secp256k1Generate();
        this.privateKey = keyPair.getPrivate();
        this.publicKey = keyPair.getPublic();
    }

    public String getWalletAddress(){
        byte[] publicKeyHash = SecurityUtil.ripemd160Digest(SecurityUtil.sha256Digest(publicKey.getEncoded()));
        byte[] data = new byte[1+ publicKeyHash.length];
        data[0] = (byte) 0;
        for (int i = 0;i< publicKeyHash.length;++i){
            data[1+i] = publicKeyHash[i];
        }

        byte[] doubleHash = SecurityUtil.sha256Digest(SecurityUtil.sha256Digest(data));
        byte[] walletEncoded = new byte[1+ publicKeyHash.length+4];
        walletEncoded[0] = (byte)0;
        for (int i = 0;i< publicKeyHash.length;++i){
            walletEncoded[i+1] = publicKeyHash[i];
        }
        for (int i=0;i<4;++i){
            walletEncoded[i+ publicKeyHash.length+1]=doubleHash[i];
        }
        String walletAddress = Base58Util.encode(walletEncoded);
        return walletAddress;
    }

    public int getAmount (UTXO[] trueUtxos){
        int amount = 0;
        for (int i = 0;i<trueUtxos.length;++i){
            amount+=trueUtxos[i].getAmount();
        }
        return amount;
    }

    public PublicKey getPublicKey() {
        return publicKey;
    }

    public PrivateKey getPrivateKey() {
        return privateKey;
    }

    @Override
    public String toString() {
        return "Account{" +
                "publicKey=" + SecurityUtil.bytes2HexString(publicKey.getEncoded()) +
                ", privateKey=" + SecurityUtil.bytes2HexString(privateKey.getEncoded()) +
                '}';
    }
}
