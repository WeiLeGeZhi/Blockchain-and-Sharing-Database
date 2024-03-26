package unit;

import consensus.MinerPeer;
import data.*;
import network.NetWork;
import utils.SecurityUtil;

import java.util.ArrayList;
import java.util.List;

/**
 * @ClassName: UtxoTest
 * @Description: TODO
 * @Author xtm
 * @Date: 2022/3/4 15:23
 * @Version 1.0
 */
public class UtxoTest {

    @org.junit.Test
    public void utxoTest() {
        // 相关类初始化
        BlockChain blockChain = new BlockChain();   // 初始一个链，内部会创建accounts
        TransactionPool transactionPool = new TransactionPool(1);   // 交易池有一个交易就会被矿工打包
        NetWork netWork = new NetWork();
        MinerPeer minerPeer = new MinerPeer(blockChain, netWork);

        // 生成一笔特殊交易
        Transaction transaction = getOneTransaction(blockChain);
        // 将该交易放入交易池中
        transactionPool.put(transaction);
        // 矿工工作
        //minerNode.run();
    }

    /**
     * 生成一笔特殊交易：accounts[1] 支付给 accounts[2] 1000元, accounts[1]使用自己的公钥对交易签名
     * @param blockChain
     * @return
     */
    Transaction getOneTransaction(BlockChain blockChain) {
        Transaction transaction = null;
        Account[] accounts = blockChain.getAccounts();

        Account aAccount = accounts[1];
        Account bAccount = accounts[2];
        String aWalletAddress = aAccount.getWalletAddress();
        String bWalletAddress = bAccount.getWalletAddress();

        UTXO[] aTrueUtxos = blockChain.getTrueUtxos(aWalletAddress);

        int txAmount = 1000;
        List<UTXO> inUtxoList = new ArrayList<>();
        List<UTXO> outUtxoList = new ArrayList<>();

        byte[] aUnlockSign = SecurityUtil.signature(aAccount.getPublicKey().getEncoded(), aAccount.getPrivateKey());

        int inAmount = 0;
        for(UTXO utxo:aTrueUtxos){
            if(utxo.unlockScript(aUnlockSign, aAccount.getPublicKey())){
                inAmount += utxo.getAmount();
                inUtxoList.add(utxo);
                if(inAmount>=txAmount){
                    break;
                }
            }
        }
        outUtxoList.add(new UTXO(bWalletAddress, txAmount, bAccount.getPublicKey()));
        if(inAmount>txAmount){
            outUtxoList.add(new UTXO(aWalletAddress, inAmount- txAmount, aAccount.getPublicKey()));
        }
        UTXO[] inUtxos = inUtxoList.toArray(new UTXO[0]);
        UTXO[] outUtxos = outUtxoList.toArray(new UTXO[0]);

        byte[] data = SecurityUtil.utxos2Bytes(inUtxos,outUtxos);
        byte[] sign = SecurityUtil.signature(data, aAccount.getPrivateKey());
        long timestamp = System.currentTimeMillis();
        transaction = new Transaction(inUtxos, outUtxos, sign, aAccount.getPublicKey(), timestamp);

        return transaction;
    }
}

