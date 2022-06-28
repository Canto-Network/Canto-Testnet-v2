package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/Canto-Network/Canto-Testnet-v2/v0/x/incentives/client/cli"
	"github.com/Canto-Network/Canto-Testnet-v2/v0/x/incentives/client/rest"
)

var (
	RegisterIncentiveProposalHandler = govclient.NewProposalHandler(cli.NewRegisterIncentiveProposalCmd, rest.RegisterIncentiveProposalRESTHandler)
	CancelIncentiveProposalHandler   = govclient.NewProposalHandler(cli.NewCancelIncentiveProposalCmd, rest.CancelIncentiveProposalRequestRESTHandler)
)
