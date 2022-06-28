<!--
order: 8
-->

# Clients

A user can query the `x/incentives` module using the CLI, JSON-RPC, gRPC or
REST.

## CLI

Find below a list of `cantod` commands added with the `x/inflation` module. You
can obtain the full list by using the `cantod -h` command.

### Queries

The `query` commands allow users to query `inflation` state.

**`period`**

Allows users to query the current inflation period.

```go
cantod query inflation period [flags]
```

**`epoch-mint-provision`**

Allows users to query the current inflation epoch provisions value.

```go
cantod query inflation epoch-mint-provision [flags]
```

**`skipped-epochs`**

Allows users to query the current number of skipped epochs.

```go
cantod query inflation skipped-epochs [flags]
```

**`total-supply`**

Allows users to query the total supply of tokens in circulation.

```go
cantod query inflation total-supply [flags]
```

**`inflation-rate`**

Allows users to query the inflation rate of the current period.

```go
cantod query inflation inflation-rate [flags]
```

**`params`**

Allows users to query the current inflation parameters.

```go
cantod query inflation params [flags]
```

### Proposals

The `tx gov submit-proposal` commands allow users to query create a proposal
using the governance module CLI:

**`param-change`**

Allows users to submit a `ParameterChangeProposal`.

```bash
cantod tx gov submit-proposal param-change [proposal-file] [flags]
```

## gRPC

### Queries

| Verb   | Method                                        | Description                                   |
| ------ | --------------------------------------------- | --------------------------------------------- |
| `gRPC` | `canto.inflation.v1.Query/Period`             | Gets current inflation period                 |
| `gRPC` | `canto.inflation.v1.Query/EpochMintProvision` | Gets current inflation epoch provisions value |
| `gRPC` | `canto.inflation.v1.Query/Params`             | Gets current inflation parameters             |
| `gRPC` | `canto.inflation.v1.Query/SkippedEpochs`      | Gets current number of skipped epochs         |
| `gRPC` | `canto.inflation.v1.Query/TotalSupply`        | Gets current total supply                     |
| `gRPC` | `canto.inflation.v1.Query/InflationRate`      | Gets current inflation rate                   |
| `GET`  | `/canto/inflation/v1/period`                  | Gets current inflation period                 |
| `GET`  | `/canto/inflation/v1/epoch_mint_provision`    | Gets current inflation epoch provisions value |
| `GET`  | `/canto/inflation/v1/skipped_epochs`          | Gets current number of skipped epochs         |
| `GET`  | `/canto/inflation/v1/total_supply`          | Gets current total supply                     |
| `GET`  | `/canto/inflation/v1/inflation_rate`          | Gets current inflation rate                   |
| `GET`  | `/canto/inflation/v1/params`                  | Gets current inflation parameters             |
