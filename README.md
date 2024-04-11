# GovGen

> [!WARNING]
> THE CODE IN THIS REPOSITORY **HAS NOT BEEN AUDITED YET**.
>
> PLEASE USE **EXTREME CAUTION** WHEN USING THIS SOFTWARE, AND USE IT AT YOUR OWN RISK.
> FOR THE TIME BEING, WE ADVISE THAT YOU NOT USE IT WITH YOUR PERSONAL PRIVATE KEY(S).
>
> THIS IS **ESPECIALLY IMPORTANT** AS GOVGEN RELIES ON AND USES ACCOUNTS DERIVED FROM
> THE COSMOS HUB, AND THEREFORE THERE IS RISK OF COMPROMISING YOUR COSMOS HUB
> ACCOUNT AS WELL.

GovGen is built using the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) as a fork of the
[Cosmos Hub](https://github.com/cosmos/gaia) at version [v14.1.0](https://github.com/cosmos/gaia/releases/tag/v14.1.0).

The following modifications have been made to the Cosmos Hub software to create GovGen:

1. Removed x/globalfee module and revert to older and simpler fee decorator
2. Removed IBC and related modules (e.g. ICA, Packet Forwarding Middleware, etc.)
3. Removed Interchain Security module
4. Reverted to standard Cosmos SDK v0.46.16 without the Liquid Staking Module (LSM)
5. Changed Bech32 prefixes to `govgen` (see `cmd/govgend/cmd/config.go`)
6. Reduced hard-coded ante min-deposit percentage to 1% (see `ante/gov_ante.go:minInitialDepositFraction`)
7. Removed ability for validators to vote on proposals with delegations, they can only use their own stake
8. Removed community spend proposal
9. Allowed setting different voting periods for different proposal types
10. Stake automatically 50% of balance for accounts that have more than 25 $GOVGEN at genesis initialization. The resulting stake distribution will provide approximately the same voting power to all genesis validators. Accounts will automatically stake to a maximum of 5 validators if 50% of the balance is less than 500 $GOVGEN, a maximum of 10 validators if less than 10,000 $GOVGEN and a maximum of 20 validators if more, uniformly. The number of validators elected for the delegations is not a constant because it depends on the state of the distribution.
