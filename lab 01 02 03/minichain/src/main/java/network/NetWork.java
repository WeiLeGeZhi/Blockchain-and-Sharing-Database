package network;

import config.MiniChainConfig;
import consensus.MinerPeer;
import consensus.TransactionProducer;
import data.*;
import spv.SpvPeer;
import utils.SecurityUtil;

import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;

/**
 * 该类模拟一个网络环境，在该网络中主要有区块链和矿工，另外地，出于工程实现的角度，还有一个交易池和一个生成随机交易的线程
 *
 */
public class NetWork {

    private final Account[] accounts;
    private final MinerPeer minerPeer;
    private final BlockChain blockChain = new BlockChain();
    private final TransactionPool transactionPool;
    private final TransactionProducer transactionProducer;
    private final SpvPeer[] spvPeers;

    /**
     * 系统中几个主要成员的初始化
     */
    public NetWork() {
        System.out.println("\naccounts and spvPeers config...");
        accounts = new Account[MiniChainConfig.ACCOUNT_NUM];
        spvPeers = new SpvPeer[MiniChainConfig.ACCOUNT_NUM];
        for (int i = 0;i<MiniChainConfig.ACCOUNT_NUM;++i){
            accounts[i] = new Account();
            System.out.println("network register new account: "+accounts[i]);
            spvPeers[i] = new SpvPeer(accounts[i], this);
        }
        System.out.println("\ntransactionPool config...");
        transactionPool = new TransactionPool(MiniChainConfig.MAX_TRANSACTION_COUNT);

        System.out.println("\ntransactionProducer config...");
        transactionProducer = new TransactionProducer(this);

        System.out.println("\nminerPeer config...");
        minerPeer = new MinerPeer(blockChain, this);

        System.out.println("\nnetwork start!\n");

        minerPeer.broadcast(blockChain.getLatestBlock());

        theyHaveADayDream();
    }

    public void theyHaveADayDream() {

        // 在创世区块中为每个账户分配一定金额的 utxo，便于后面交易的进行
        UTXO[] outUtxos = new UTXO[accounts.length];
        for (int i = 0;  i < accounts.length; ++i) {
            outUtxos[i] = new UTXO(accounts[i].getWalletAddress(), MiniChainConfig.INIT_AMOUNT, accounts[i].getPublicKey());
        }
        // 神秘的公私钥
        KeyPair dayDreamKeyPair = SecurityUtil.secp256k1Generate();
        PublicKey dayDreamPublicKey = dayDreamKeyPair.getPublic();
        PrivateKey dayDreamPrivateKey = dayDreamKeyPair.getPrivate();
        // 神秘的签名内容
        byte[] sign = SecurityUtil.signature("Everything in the dream!".getBytes(StandardCharsets.UTF_8), dayDreamPrivateKey);
        // 构造交易
        Transaction transaction = new Transaction(new UTXO[]{}, outUtxos, sign, dayDreamPublicKey, System.currentTimeMillis());
        // 交易数组只有这一个交易
        Transaction[] transactions = { transaction };
        // 前一个区块的哈希
        String preBlockHash = SecurityUtil.sha256Digest(blockChain.getLatestBlock().toString());
        // 因为本区块只有一个交易，所以merkle根哈希即为该交易的哈希
        String merkleRootHash = SecurityUtil.sha256Digest(transaction.toString());
        // 构建区块
        BlockHeader blockHeader = new BlockHeader(preBlockHash, merkleRootHash, Math.abs(new Random().nextLong()));
        BlockBody blockBody = new BlockBody(merkleRootHash, transactions);
        Block block = new Block(blockHeader, blockBody);
        // 添加到链中
        blockChain.addNewBlock(block);

        minerPeer.broadcast(block);

    }

    public List<Transaction> getTransactionInLatestBlock(String walletAddress){
        List<Transaction> list = new ArrayList<>();
        Block block = blockChain.getLatestBlock();

        for (Transaction transaction: block.getBlockBody().getTransactions()){
            boolean have = false;
            for (UTXO utxo: transaction.getInUtxos()){
                if (utxo.getWalletAddress().equals(walletAddress)){
                    have = true;
                    break;
                }
            }
            if (have){
                continue;
            }
            for (UTXO utxo: transaction.getOutUtxos()){
                if (utxo.getWalletAddress().equals(walletAddress)){
                    list.add(transaction);
                    break;
                }
            }
        }
        return list;
    }

    /**
     * 启动挖矿线程和生成随机交易的线程
     */
    public void start() {
        transactionProducer.start();
        minerPeer.start();
    }

    public Account[] getAccounts() {
        return accounts;
    }

    public TransactionPool getTransactionPool() {
        return transactionPool;
    }

    public BlockChain getBlockChain() {
        return blockChain;
    }

    public MinerPeer getMinerPeer() {
        return minerPeer;
    }

    public SpvPeer[] getSpvPeers() {
        return spvPeers;
    }
}
