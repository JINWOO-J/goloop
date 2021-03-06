package i;

import a.ByteArray;
import p.score.Address;
import p.score.AnyDB;
import p.score.VarDB;
import s.java.lang.Class;
import s.java.lang.String;
import s.java.math.BigInteger;

/**
 * Interface to the blockchain runtime.
 */
public interface IBlockchainRuntime {
    //================
    // Transaction
    //================

    /**
     * Returns the transaction hash of the origin transaction.
     */
    ByteArray avm_getTransactionHash();

    /**
     * Returns the transaction index in a block.
     */
    int avm_getTransactionIndex();

    /**
     * Returns the timestamp of a transaction request.
     */
    long avm_getTransactionTimestamp();

    /**
     * Returns the nonce of a transaction request.
     */
    BigInteger avm_getTransactionNonce();

    /**
     * Returns the address of the currently-running SCORE.
     */
    Address avm_getAddress();

    /**
     * Returns the caller's address.
     */
    Address avm_getCaller();

    /**
     * Returns the originator's address.
     */
    Address avm_getOrigin();

    /**
     * Returns the address of the account who deployed the contract.
     */
    Address avm_getOwner();

    /**
     * Returns the value being transferred along the transaction.
     */
    BigInteger avm_getValue();

    //================
    // Block
    //================

    /**
     * Block timestamp.
     *
     * @return The time of the current block, as seconds since the Epoch.
     */
    long avm_getBlockTimestamp();

    /**
     * Block height.
     *
     * @return The height of the current block.
     */
    long avm_getBlockHeight();

    //================
    // Storage
    //================

    /**
     * Returns the balance of an account.
     *
     * @param address account address
     * @return the balance of the account
     */
    BigInteger avm_getBalance(Address address) throws IllegalArgumentException;

    //================
    // System
    //================

    /**
     * Calls the contract denoted by the targetAddress.  Returns the response of the contract.
     *
     * @param value         The value to transfer
     * @param targetAddress The address of the contract to call.
     * @param method        method
     * @param params        parameters
     * @return The response of executing the contract.
     */
    IObject avm_call(BigInteger value, Address targetAddress, String method,
                     IObjectArray params);

    /**
     * Stop the current execution, rollback any state changes, and refund the remaining energy to caller.
     */
    void avm_revert(int code, String message);

    void avm_revert(int code);

    /**
     * Requires that condition is true, otherwise triggers a revert.
     */
    void avm_require(boolean condition);

    /**
     * Prints a message to console for debugging purpose
     */
    void avm_println(String message);

    /**
     * Computes the SHA3-256 hash using the input data
     */
    ByteArray avm_sha3_256(ByteArray data);

    /**
     * Computes the SHA-256 hash using the input data
     */
    ByteArray avm_sha256(ByteArray data);

    /**
     * Recovers the public key from the message hash and the recoverable signature
     */
    ByteArray avm_recoverKey(ByteArray msgHash, ByteArray signature, boolean compressed);

    /**
     * Returns the address that is associated with the given public key
     */
    Address avm_getAddressFromKey(ByteArray publicKey);

    /**
     * Returns a new AnyDB instance
     */
    AnyDB avm_newAnyDB(String id, Class<?> vc);

    /**
     * Emits event logs
     */
    void avm_logEvent(IObjectArray indexed, IObjectArray data);
}
