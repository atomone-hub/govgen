package client

import (
	"github.com/spf13/cobra"

	"github.com/atomone-hub/govgen/v1/x/gov/client/rest"

	"github.com/cosmos/cosmos-sdk/client"
)

// function to create the rest handler
type RESTHandlerFn func(client.Context) rest.ProposalRESTHandler

// function to create the cli handler
type CLIHandlerFn func() *cobra.Command

// The combined type for a proposal handler for both cli and rest
type ProposalHandler struct {
	CLIHandler  CLIHandlerFn
	RESTHandler RESTHandlerFn
}

// NewProposalHandler creates a new ProposalHandler object
func NewProposalHandler(cliHandler CLIHandlerFn, restHandler RESTHandlerFn) ProposalHandler {
	return ProposalHandler{
		CLIHandler:  cliHandler,
		RESTHandler: restHandler,
	}
}
