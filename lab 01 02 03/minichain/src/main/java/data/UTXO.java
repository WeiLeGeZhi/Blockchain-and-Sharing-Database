package data;

import utils.SecurityUtil;

import java.security.PublicKey;
import java.util.Arrays;
import java.util.Stack;

public class UTXO {

    private final String walletAddress;
    private final int amount;
    private final byte[] publicHash;

    public UTXO(String walletAddress, int amount, PublicKey publicKey){
        this.walletAddress = walletAddress;
        this.amount = amount;
        publicHash = SecurityUtil.ripemd160Digest(SecurityUtil.sha256Digest(publicKey.getEncoded()));
    }

    public boolean unlockScript(byte[] sign, PublicKey publicKey){
        Stack<byte[]> stack = new Stack<>();
        stack.push(sign);
        stack.push(publicKey.getEncoded());
        stack.push(stack.peek());
        byte[] data = stack.pop();
        stack.push(SecurityUtil.ripemd160Digest(SecurityUtil.sha256Digest(data)));
        stack.push(publicHash);
        byte[] publicKeyHash1 = stack.pop();
        byte[] publicKeyHash2 = stack.pop();
        if(!Arrays.equals(publicKeyHash1, publicKeyHash2)){
            return false;
        }
        byte[] publicKeyEncoded = stack.pop();
        byte[] sign1 = stack.pop();
        return SecurityUtil.verify(publicKey.getEncoded(), sign1, publicKey);
    }


    public String getWalletAddress() {
        return walletAddress;
    }

    public int getAmount() {
        return amount;
    }

    public byte[] getPublicHash() {
        return publicHash;
    }

    @Override
    public String toString() {
        return "\n\tUTXO{" +
                "walletAddress='" + walletAddress + '\'' +
                ", amount=" + amount +
                ", publicHash=" + SecurityUtil.bytes2HexString(publicHash) +
                '}';
    }
}
