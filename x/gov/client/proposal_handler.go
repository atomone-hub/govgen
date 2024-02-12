package client

import (
	"github.com/spf13/cobra"

	"github.com/atomone-hub/govgen/v1/x/gov/client/rest"

	"github.com/cosmos/cosmos-sdk/client"
	legacyclient "github.com/cosmos/cosmos-sdk/x/gov/client"
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

func WrapPropposalHandler(h legacyclient.ProposalHandler) ProposalHandler {
	return ProposalHandler{
		CLIHandler: func() *cobra.Command {
			return h.CLIHandler()
		},
		RESTHandler: func(ctx client.Context) rest.ProposalRESTHandler {
			handler := h.RESTHandler(ctx)
			return rest.ProposalRESTHandler{
				SubRoute: handler.SubRoute,
				Handler:  handler.Handler,
			}
		},
	}
}
