package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	v2 "github.com/evmos/evmos/v9/app/upgrades/v2"
	v4 "github.com/evmos/evmos/v9/app/upgrades/v4"
	v7 "github.com/evmos/evmos/v9/app/upgrades/v7"
	v82 "github.com/evmos/evmos/v9/app/upgrades/v8_2"
)

// ScheduleForkUpgrade executes any necessary fork logic for based upon the current
// block height and chain ID (mainnet or testnet). It sets an upgrade plan once
// the chain reaches the pre-defined upgrade height.
//
// CONTRACT: for this logic to work properly it is required to:
//
//  1. Release a non-breaking patch version so that the chain can set the scheduled upgrade plan at upgrade-height.
//  2. Release the software defined in the upgrade-info
func (app *Evmos) ScheduleForkUpgrade(ctx sdk.Context) {
	upgradePlan := upgradetypes.Plan{
		Height: ctx.BlockHeight(),
	}

	// handle mainnet forks with their corresponding upgrade name and info
	switch ctx.BlockHeight() {
	case v2.MainnetUpgradeHeight:
		upgradePlan.Name = v2.UpgradeName
		upgradePlan.Info = v2.UpgradeInfo
	case v4.MainnetUpgradeHeight:
		upgradePlan.Name = v4.UpgradeName
		upgradePlan.Info = v4.UpgradeInfo
	case v7.MainnetUpgradeHeight:
		upgradePlan.Name = v7.UpgradeName
		upgradePlan.Info = v7.UpgradeInfo
	case v82.MainnetUpgradeHeight:
		upgradePlan.Name = v82.UpgradeName
		upgradePlan.Info = v82.UpgradeInfo
	case 8800:
		upgradePlan.Name = "v9.1.3"
		upgradePlan.Info = `'{"binaries":{"darwin/amd64":"https://github.com/hanchon/evmos/releases/download/v9.1.3/evmos_9.1.3_Darwin_arm64.tar.gz","darwin/x86_64":"https://github.com/hanchon/evmos/releases/download/v9.1.3/evmos_9.1.3_Darwin_x86_64.tar.gz","linux/arm64":"https://github.com/hanchon/evmos/releases/download/v9.1.3/evmos_9.1.3_Linux_arm64.tar.gz","linux/amd64":"https://github.com/hanchon/evmos/releases/download/v9.1.3/evmos_9.1.3_Linux_amd64.tar.gz","windows/x86_64":"https://github.com/hanchon/evmos/releases/download/v9.1.3/evmos_9.1.3_Windows_x86_64.zip"}}'`
	default:
		// No-op
		return
	}

	// schedule the upgrade plan to the current block hight, effectively performing
	// a hard fork that uses the upgrade handler to manage the migration.
	if err := app.UpgradeKeeper.ScheduleUpgrade(ctx, upgradePlan); err != nil {
		panic(
			fmt.Errorf(
				"failed to schedule upgrade %s during BeginBlock at height %d: %w",
				upgradePlan.Name, ctx.BlockHeight(), err,
			),
		)
	}
}
