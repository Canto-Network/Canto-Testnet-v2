package keeper_test

import "github.com/Canto-Network/Canto-Testnet-v2/v0/x/incentives/types"

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.IncentivesKeeper.GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)
	params.EnableIncentives = false
	suite.app.IncentivesKeeper.SetParams(suite.ctx, params)
	newParams := suite.app.IncentivesKeeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}
