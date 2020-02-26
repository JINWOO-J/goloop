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

package testcases;

import score.Address;
import score.Context;
import foundation.icon.ee.tooling.abi.External;
import foundation.icon.ee.tooling.abi.EventLog;
import foundation.icon.ee.tooling.abi.Optional;
import foundation.icon.ee.tooling.abi.Payable;

import java.math.BigInteger;

public class APITest
{
    public APITest() {
    }

    @EventLog
    public void EmitEvent(byte[] data) {}

    //================================
    // Address
    //================================

    @External
    public void getAddress(Address addr) {
        Context.require(Context.getAddress().equals(addr));
    }

    @External(readonly=true)
    public Address getAddressQuery() {
        return Context.getAddress();
    }

    @External
    public void getCaller(Address caller) {
        Context.require(Context.getCaller().equals(caller));
    }

    @External(readonly=true)
    public Address getCallerQuery() {
        return Context.getCaller();
    }

    @External
    public void getOrigin(Address origin) {
        Context.require(Context.getOrigin().equals(origin));
    }

    @External(readonly=true)
    public Address getOriginQuery() {
        return Context.getOrigin();
    }

    @External
    public void getOwner(Address owner) {
        Context.require(Context.getOwner().equals(owner));
    }

    @External(readonly=true)
    public Address getOwnerQuery() {
        return Context.getOwner();
    }

    //================================
    // Block
    //================================

    @External
    public void getBlockTimestamp() {
        Context.require(Context.getBlockTimestamp() > 0L);
        EmitEvent(BigInteger.valueOf(Context.getBlockTimestamp()).toByteArray());
    }

    @External(readonly=true)
    public long getBlockTimestampQuery() {
        return Context.getBlockTimestamp();
    }

    @External
    public void getBlockHeight() {
        Context.require(Context.getBlockHeight() > 0L);
        EmitEvent(BigInteger.valueOf(Context.getBlockHeight()).toByteArray());
    }

    @External(readonly=true)
    public long getBlockHeightQuery() {
        return Context.getBlockHeight();
    }

    //================================
    // Transaction
    //================================

    @External
    public void getTransactionHash() {
        Context.require(Context.getTransactionHash() != null);
        EmitEvent(Context.getTransactionHash());
    }

    @External(readonly=true)
    public byte[] getTransactionHashQuery() {
        return Context.getTransactionHash();
    }

    @External
    public void getTransactionIndex() {
        Context.require(Context.getTransactionIndex() >= 0);
        EmitEvent(BigInteger.valueOf(Context.getTransactionIndex()).toByteArray());
    }

    @External(readonly=true)
    public int getTransactionIndexQuery() {
        return Context.getTransactionIndex();
    }

    @External
    public void getTransactionTimestamp() {
        Context.require(Context.getTransactionTimestamp() > 0L);
        EmitEvent(BigInteger.valueOf(Context.getTransactionTimestamp()).toByteArray());
    }

    @External(readonly=true)
    public long getTransactionTimestampQuery() {
        return Context.getTransactionTimestamp();
    }

    @External
    public void getTransactionNonce() {
        EmitEvent(Context.getTransactionNonce().toByteArray());
    }

    @External(readonly=true)
    public BigInteger getTransactionNonceQuery() {
        return Context.getTransactionNonce();
    }

    //================================
    // ICX coin
    //================================

    @External
    @Payable
    public void getValue() {
        EmitEvent(Context.getValue().toByteArray());
    }

    @External(readonly=true)
    public BigInteger getValueQuery() {
        return Context.getValue();
    }

    @External
    public void getBalance(@Optional Address address) {
        if (address == null) {
            address = Context.getAddress();
        }
        EmitEvent(Context.getBalance(address).toByteArray());
    }

    @External(readonly=true)
    public BigInteger getBalanceQuery(@Optional Address address) {
        if (address == null) {
            address = Context.getAddress();
        }
        return Context.getBalance(address);
    }

    //================================
    // Crypto
    //================================

    @External
    public void computeHash(byte[] data) {
        EmitEvent(Context.sha3_256(data));
    }

    @External(readonly=true)
    public byte[] computeHashQuery(byte[] data) {
        return Context.sha3_256(data);
    }
}