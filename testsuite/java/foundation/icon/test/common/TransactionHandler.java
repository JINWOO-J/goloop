/*
 * Copyright 2019 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package foundation.icon.test.common;

import foundation.icon.ee.tooling.deploy.OptimizedJarBuilder;
import foundation.icon.icx.Call;
import foundation.icon.icx.IconService;
import foundation.icon.icx.SignedTransaction;
import foundation.icon.icx.Transaction;
import foundation.icon.icx.TransactionBuilder;
import foundation.icon.icx.Wallet;
import foundation.icon.icx.data.Address;
import foundation.icon.icx.data.Bytes;
import foundation.icon.icx.data.ConfirmedTransaction;
import foundation.icon.icx.data.ScoreApi;
import foundation.icon.icx.data.TransactionResult;
import foundation.icon.icx.transport.jsonrpc.RpcError;
import foundation.icon.icx.transport.jsonrpc.RpcItem;
import foundation.icon.icx.transport.jsonrpc.RpcObject;
import foundation.icon.test.score.GovScore;
import foundation.icon.test.score.Score;
import org.aion.avm.utilities.JarBuilder;

import java.io.IOException;
import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import static foundation.icon.test.common.Env.LOG;

public class TransactionHandler {
    private final IconService iconService;
    private final Env.Chain chain;

    public TransactionHandler(IconService iconService, Env.Chain chain) {
        this.iconService = iconService;
        this.chain = chain;
    }

    public Score deploy(Wallet owner, String scorePath, RpcObject params)
            throws IOException, ResultTimeoutException, TransactionFailureException {
        return deploy(owner, scorePath, params, Constants.DEFAULT_INSTALL_STEPS);
    }

    public Score deploy(Wallet owner, byte[] content, RpcObject params)
            throws IOException, ResultTimeoutException, TransactionFailureException {
        return getScore(doDeploy(owner, content, params, Constants.CONTENT_TYPE_PYTHON));
    }

    public Score deploy(Wallet owner, String scorePath, RpcObject params, BigInteger steps)
            throws IOException, ResultTimeoutException, TransactionFailureException {
        return deploy(owner, scorePath, Constants.CHAINSCORE_ADDRESS, params, steps);
    }

    public Score deploy(Wallet owner, String scorePath, Address to, RpcObject params, BigInteger steps)
            throws IOException, ResultTimeoutException, TransactionFailureException {
        if (scorePath.endsWith(".jar")) {
            byte[] data = Files.readAllBytes(Path.of(scorePath));
            return getScore(doDeploy(owner, data, to, params, steps, Constants.CONTENT_TYPE_JAVA));
        } else {
            byte[] data = ZipFile.zipContent(scorePath);
            return getScore(doDeploy(owner, data, to, params, steps, Constants.CONTENT_TYPE_PYTHON));
        }
    }

    public Score deploy(Wallet owner, Class<?> mainClass, RpcObject params)
            throws IOException, ResultTimeoutException, TransactionFailureException {
        return deploy(owner, new Class<?>[]{mainClass}, params);
    }

    public Score deploy(Wallet owner, Class<?>[] classes, RpcObject params)
            throws IOException, ResultTimeoutException, TransactionFailureException {
        byte[] jar = makeJar(classes[0].getName(), classes);
        return getScore(doDeploy(owner, jar, params, Constants.CONTENT_TYPE_JAVA));
    }

    private byte[] makeJar(String name, Class<?>[] classes) {
        byte[] jarBytes = JarBuilder.buildJarForExplicitMainAndClasses(name, classes);
        return new OptimizedJarBuilder(false, jarBytes, true)
                .withUnreachableMethodRemover()
                .withRenamer()
                .getOptimizedBytes();
    }

    private Bytes doDeploy(Wallet owner, byte[] content, RpcObject params, String contentType)
            throws IOException {
        return doDeploy(owner, content, Constants.CHAINSCORE_ADDRESS, params, Constants.DEFAULT_INSTALL_STEPS, contentType);
    }

    private Bytes doDeploy(Wallet owner, byte[] content, Address to, RpcObject params, BigInteger steps, String contentType)
            throws IOException {
        Transaction transaction = TransactionBuilder.newBuilder()
                .nid(getNetworkId())
                .from(owner.getAddress())
                .to(to)
                .stepLimit(steps)
                .deploy(contentType, content)
                .params(params)
                .build();
        SignedTransaction signedTransaction = new SignedTransaction(transaction, owner);
        return iconService.sendTransaction(signedTransaction).execute();
    }

    public Score getScore(Bytes txHash)
            throws IOException, ResultTimeoutException, TransactionFailureException {
        TransactionResult result = getResult(txHash);
        if (!Constants.STATUS_SUCCESS.equals(result.getStatus())) {
            throw new TransactionFailureException(result.getFailure());
        }
        acceptScoreIfAuditEnabled(txHash);
        return new Score(this, new Address(result.getScoreAddress()));
    }

    public Bytes deployOnly(Wallet owner, String scorePath, RpcObject params) throws IOException {
        return deployOnly(owner, Constants.CHAINSCORE_ADDRESS, scorePath, params);
    }

    public Bytes deployOnly(Wallet owner, Address to, String scorePath, RpcObject params) throws IOException {
        byte[] data = ZipFile.zipContent(scorePath);
        return doDeploy(owner, data, to, params,
                Constants.DEFAULT_INSTALL_STEPS,
                Constants.CONTENT_TYPE_PYTHON);
    }

    public IconService getIconService() {
        return this.iconService;
    }

    public Env.Chain getChain() {
        return this.chain;
    }

    public BigInteger getNetworkId() {
        return BigInteger.valueOf(chain.networkId);
    }

    public BigInteger getBalance(Address address) throws IOException {
        return iconService.getBalance(address).execute();
    }

    public List<ScoreApi> getScoreApi(Address scoreAddress) throws IOException {
        return iconService.getScoreApi(scoreAddress).execute();
    }

    public RpcItem call(Call<RpcItem> call) throws IOException {
        return this.iconService.call(call).execute();
    }

    public Bytes invoke(Wallet wallet, Transaction tx) throws IOException {
        return this.iconService.sendTransaction(new SignedTransaction(tx, wallet)).execute();
    }

    public TransactionResult invokeAndWait(Wallet wallet, Transaction tx) throws IOException {
        return this.iconService.sendTransactionAndWait(new SignedTransaction(tx, wallet)).execute();
    }

    public TransactionResult waitResult(Bytes txHash) throws IOException {
        return this.iconService.waitTransactionResult(txHash).execute();
    }

    public TransactionResult getResult(Bytes txHash)
            throws IOException, ResultTimeoutException {
        return getResult(txHash, Constants.DEFAULT_WAITING_TIME);
    }

    public TransactionResult getResult(Bytes txHash, long waiting)
            throws IOException, ResultTimeoutException {
        long limitTime = System.currentTimeMillis() + waiting;
        while (true) {
            try {
                return iconService.getTransactionResult(txHash).execute();
            } catch (RpcError e) {
                if (e.getCode() == -31002 || e.getCode() == -31003) { // pending or executing
                    if (limitTime < System.currentTimeMillis()) {
                        throw new ResultTimeoutException(txHash);
                    }
                    try {
                        // wait until block confirmation
                        LOG.debug("RpcError: code(" + e.getCode() + ") message(" + e.getMessage() + "); Retry in 1 sec.");
                        Thread.sleep(1000);
                    } catch (InterruptedException ex) {
                        ex.printStackTrace();
                    }
                    continue;
                }
                throw e;
            }
        }
    }

    public Bytes transfer(Address to, BigInteger amount) throws IOException {
        return transfer(chain.godWallet, to, amount);
    }

    public Bytes transfer(Wallet owner, Address to, BigInteger amount) throws IOException {
        Transaction transaction = TransactionBuilder.newBuilder()
                .nid(getNetworkId())
                .from(owner.getAddress())
                .to(to)
                .value(amount)
                .stepLimit(Constants.DEFAULT_STEPS)
                .build();
        SignedTransaction signedTransaction = new SignedTransaction(transaction, owner);
        return iconService.sendTransaction(signedTransaction).execute();
    }

    public TransactionResult acceptScoreIfAuditEnabled(Bytes txHash)
            throws TransactionFailureException, IOException, ResultTimeoutException {
        GovScore govScore = new GovScore(this);
        if (govScore.isAuditEnabledOnly()) {
            LOG.infoEntering("invoke", "acceptScore");
            TransactionResult result = govScore.acceptScore(txHash);
            if (!Constants.STATUS_SUCCESS.equals(result.getStatus())) {
                LOG.infoExiting();
                throw new TransactionFailureException(result.getFailure());
            }
            LOG.infoExiting();
            return result;
        }
        return null;
    }

    public ConfirmedTransaction getTransaction(Bytes txHash) throws IOException {
        return iconService.getTransaction(txHash).execute();
    }
}
