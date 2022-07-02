<!--
order: 1
-->

# Concepts

## Key generation

`secp256k1` refers to the parameters of the elliptic curve used in generating cryptographic public keys. Like Bitcoin, IBC compatible chains like the Cosmos chain use `secp256k1` for public key generation.

Some chains use different elliptic curves for generating public keys. An example is the`eth_secp256k1`used by Ethereum and canto chain for generating public keys.

```go
// Generate new random ethsecp256k1 private key and address

ethPrivKey, err := ethsecp256k1.GenerateKey()
ethsecpAddr := sdk.AccAddress(ethPrivKey.PubKey().Address())

// Bech32 "canto" address
ethsecpAddrcanto := sdk.AccAddress(ethPk.PubKey().Address()).String()

// We can also change the HRP to use "cosmos"
ethsecpAddrCosmos := sdk.MustBech32ifyAddressBytes(sdk.Bech32MainPrefix, ethsecpAddr)
```

The above example code demonstrates a simple user account creation on canto.
On the second line, a private key is generated using the `eth_secp256k1` curve, which is used to create a human readable `PubKey` string.
For more detailed info on accounts, please check the [accounts section](https://canto.dev/technical_concepts/accounts.html#canto-accounts) in the official canto documentation.

## Stuck funds

The primary use case of the `x/recovery` module is to enable the recovery of tokens, that were sent to unsupported canto addresses. These tokens are termed “stuck”, as the account’s owner cannot sign transactions that transfer the tokens to other accounts. The owner only holds the private key to sign transactions for its `eth_secp256k1` public keys on canto, not other unsupported keys (i.e. `secp256k1` keys) .They are unable to transfer the tokens using the keys of the accounts through which they were sent due to the incompatibility of their elliptic curves.

## Recovery

After the initial canto launch (`v1.1.2`), tokens got stuck from accounts with and without claims records (airdrop allocation):

1. Osmosis/Cosmos Hub account without claims record sent IBC transfer to canto `secp256k1` receiver address

    **Consequences**

    - IBC vouchers from IBC transfer got stuck in the receiver’s balance

    **Recovery procedure**

    - The receiver can send an IBC transfer from their Osmosis / Cosmos Hub  account (i.e `osmo1...` or `cosmos1...`) to its same canto account (`canto1...`) to recover the tokens by forwarding them to the corresponding sending chain (Osmosis or Cosmos Hub)
2. Osmosis/Cosmos Hub account with claims record sent IBC transfer to canto `secp256k1` receiver address

    **Consequences**

    - IBC vouchers  from IBC transfer got stuck in the receiver’s balance
    - IBC Transfer Action was claimed and the canto rewards were transferred to the receiver’s canto `secp256k1` account, resulting in stuck canto tokens.
    - Claims record of the sender was migrated to the receiver’s canto `secp256k1` account

    **Recovery procedure**

    - The receiver can send an IBC transfer from their Osmosis / Cosmos Hub  account (i.e `osmo1...` or `cosmos1...`) to its same canto account (`canto1...`)  to recover the tokens by forwarding them to the corresponding sending chain (Osmosis or Cosmos Hub)
    - Migrate once again the claims record to a valid account so that the remaining 3 actions can be claimed
    - Chain is restarted with restored Claims records

## IBC Middleware Stack

### Middleware ordering

The IBC middleware adds custom logic between the core IBC and the underlying application. Middlewares are implemented as stacks so that applications can define multiple layers of custom behavior.

The order of middleware matters. Function calls from IBC core to the application travel from top-level middleware to the bottom middleware and then to the application, whereas function calls from the application to IBC core go through the bottom middleware first and then in order to the top middleware and then to core IBC handlers. Thus, the same set of middleware put in different orders may produce different effects.

During packet execution each middleware in the stack will be executed in the order defined on creation (from top to bottom).

For canto the middleware stack ordering is defined as follows (from top to bottom):

1. IBC Transfer
2. Claims Middleware
3. Recovery Middleware

This means that the IBC transfer will be executed first, then the claim will be attempted and lastly the recovery will be executed. By performing the actions in this order we allow the users to receive back the coins used to trigger the recover.

**Example execution order**

1. User attempts to recover `1000acanto` that are stuck on the canto chain.
2. User sends `100uosmo` from Osmosis to canto through an IBC transaction.
3. canto receives the transaction, and goes through the IBC stack:
    1. **IBC transfer**: the `100uosmo` IBC vouchers are added to the user balance on canto.
    2. **Claims Middleware**: since `sender=receiver` -> perform no-op
    3. **Recovery Middleware**: since `sender=receiver` -> recover user balance (`1000acanto` and `100uosmo`) by sending an IBC transfer from `receiver` to the `sender` on the Osmosis chain.
4. User receives `100uosmo` and `1000acanto` (IBC voucher) on Osmosis.

### Execution errors

It is possible that the IBC transaction fails in any point of the stack execution and in that case the recovery will not be triggered by the transaction, as it will rollback to the previous state.

So if at any point either the IBC transfer or the claims middleware return an error, then the recovery middleware will not be executed.