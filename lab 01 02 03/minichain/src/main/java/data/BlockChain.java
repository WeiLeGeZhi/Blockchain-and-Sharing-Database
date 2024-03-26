package data;

import config.MiniChainConfig;
import utils.SecurityUtil;

import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.*;

/**
 * 区块链的类抽象，创建该对象时会自动生成创世纪块，加入区块链中
 */
public class BlockChain {

    private final LinkedList<Block> chain = new LinkedList<>();
    private final Account[] accounts;

    public BlockChain() {
        this.accounts = new Account[MiniChainConfig.ACCOUNT_NUM];
        for (int i = 0;i<accounts.length;++i){
            accounts[i] = new Account();
        }
        Transaction[] transactions = genesisTransactions(accounts);
        BlockHeader genesisBlockHeader = new BlockHeader(null, null,
                                                            Math.abs(new Random().nextLong()));
        BlockBody genesisBlockBody = new BlockBody(null,transactions);
        Block genesisBlock = new Block(genesisBlockHeader, genesisBlockBody);
        System.out.println("Create the genesis Block! ");
        System.out.println("And the hash of genesis Block is : " + SecurityUtil.sha256Digest(genesisBlock.toString()) +
                ", you will see the hash value in next Block's preBlockHash field.");
        System.out.println();
        chain.add(genesisBlock);
    }

    private Transaction[] genesisTransactions(Account[] accounts){
        UTXO[] outUtxos = new UTXO[accounts.length];
        for(int i = 0; i<accounts.length;++i){
            outUtxos[i] = new UTXO(accounts[i].getWalletAddress(),MiniChainConfig.INIT_AMOUNT,accounts[i].getPublicKey());
        }
        KeyPair dayDreamKeyPair = SecurityUtil.secp256k1Generate();
        PublicKey dayDreamPublicKey = dayDreamKeyPair.getPublic();
        PrivateKey dayDreamPrivateKey = dayDreamKeyPair.getPrivate();
        byte[] sign = SecurityUtil.signature("Everything in the dream!".getBytes(StandardCharsets.UTF_8),dayDreamPrivateKey);
        return new Transaction[]{new Transaction(new UTXO[]{},outUtxos,sign,dayDreamPublicKey,System.currentTimeMillis())};
    }

    public UTXO[] getTrueUtxos(String walletAddress) {
        Set<UTXO> trueUtxoSet = new HashSet<>();
        for (Block block : chain) {
            BlockBody blockBody = block.getBlockBody();
            Transaction[] transactions = blockBody.getTransactions();
            for (Transaction transaction : transactions) {
                UTXO[] inUtxos = transaction.getInUtxos();
                UTXO[] outUtxos = transaction.getOutUtxos();
                for (UTXO utxo : inUtxos) {
                    if (utxo.getWalletAddress().equals(walletAddress)) {
                        trueUtxoSet.remove(utxo);
                    }
                }
                for (UTXO utxo : outUtxos) {
                    if (utxo.getWalletAddress().equals(walletAddress)) {
                        trueUtxoSet.add(utxo);
                    }
                }
            }
        }
//        if (transactionPool != null) {
//            if(!transactionPool.isEmpty()){
//                Transaction[] notAddedYetTx = transactionPool.getCurTx();
//                for (Transaction transaction:notAddedYetTx) {
//                    UTXO[] inUtxos = transaction.getInUtxos();
//                    UTXO[] outUtxos = transaction.getOutUtxos();
//                    for (UTXO utxo : inUtxos) {
//                        if (utxo.getWalletAddress().equals(walletAddress)) {
//                            trueUtxoSet.remove(utxo);
//                        }
//                    }
//                    for (UTXO utxo : outUtxos) {
//                        if (utxo.getWalletAddress().equals(walletAddress)) {
//                            trueUtxoSet.add(utxo);
//                        }
//                    }
//                }
//            }
//        }
        UTXO[] trueUtxos = new UTXO[trueUtxoSet.size()];
        trueUtxoSet.toArray(trueUtxos);
        return trueUtxos;
    }

    public Account[] getAccounts() {
        return accounts;
    }

    public int getAllAccountAmount(){
        int sumAmount = 0;
        for (int i=0;i<accounts.length;++i){
            UTXO[] trueUtxo = getTrueUtxos(accounts[i].getWalletAddress());
            sumAmount += accounts[i].getAmount(trueUtxo);
        }
        return sumAmount;
    }

    /**
     * 向区块链中添加新的满足难度条件的区块
     *
     * @param block 新的满足难度条件的区块
     */
    public void addNewBlock(Block block) {
        chain.offer(block);
    }

    /**
     * 获取区块链的最后一个区块，矿工在组装新的区块时，需要获取上一个区块的哈希值，通过该方法获得
     *
     * @return 区块链的最后一个区块
     */
    public Block getLatestBlock() {
        return chain.peekLast();
    }

    public LinkedList<Block> getBlocks() {
        return chain;
    }
}
