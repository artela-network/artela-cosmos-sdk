package gov

import (
	"fmt"
	"strings"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	govv1 "cosmossdk.io/api/cosmos/gov/v1"
	govv1beta1 "cosmossdk.io/api/cosmos/gov/v1beta1"

	"github.com/cosmos/cosmos-sdk/version"
)

const (
	FlagTitle     = "title"
	FlagDeposit   = "deposit"
	flagVoter     = "voter"
	flagDepositor = "depositor"
	flagStatus    = "status"
	FlagMetadata  = "metadata"
	FlagSummary   = "summary"
	// Deprecated: only used for v1beta1 legacy proposals.
	FlagProposal = "proposal"
	// Deprecated: only used for v1beta1 legacy proposals.
	FlagDescription = "description"
	// Deprecated: only used for v1beta1 legacy proposals.
	FlagProposalType = "type"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: govv1.Msg_ServiceDesc.ServiceName,
			// map v1beta1 as a sub-command
			SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
				"v1beta1": {
					Service: govv1beta1.Msg_ServiceDesc.ServiceName,
				},
			},
		},
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: govv1.Query_ServiceDesc.ServiceName,
			// map v1beta1 as a sub-command
			SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
				"v1beta1": {
					Service: govv1beta1.Query_ServiceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod: "Proposal",
							Use:       "proposal [proposal-id]",
							Short:     "Query details of a single proposal",
							Long: strings.TrimSpace(
								fmt.Sprintf(`Query details for a proposal. You can find the
proposal-id by running "%s query gov proposals".

Example:
$ %s query gov proposal 1
`,
									version.AppName, version.AppName,
								),
							),
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "proposal_id"},
							},
						},
						{
							RpcMethod: "Proposals",
							Use:       "proposals",
							Short:     "Query proposals with optional filters",
							Long: strings.TrimSpace(
								fmt.Sprintf(`Query for a all paginated proposals that match optional filters:

Example:
$ %s query gov proposals --depositor cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query gov proposals --voter cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query gov proposals --status (DepositPeriod|VotingPeriod|Passed|Rejected)
$ %s query gov proposals --page=2 --limit=100
`,
									version.AppName, version.AppName, version.AppName, version.AppName,
								),
							),
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
							FlagOptions: map[string]*autocliv1.FlagOptions{
								"proposal_status": {
									Name:  flagStatus,
									Usage: "(optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected",
								},
								flagVoter: {
									Name:  flagVoter,
									Usage: "(optional) filter by proposals voted on by voted",
								},
								flagDepositor: {
									Name:  flagDepositor,
									Usage: "(optional) filter by proposals deposited on by depositor",
								},
							},
						},
						{
							RpcMethod: "Vote",
							Use:       "vote [proposal-id] [voter-addr]",
							Short:     "Query details of a single vote",
							Long: strings.TrimSpace(
								fmt.Sprintf(`Query details for a single vote on a proposal given its identifier.

Example:
$ %s query gov vote 1 cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`,
									version.AppName,
								),
							),
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "proposal_id"},
								{ProtoField: "voter"},
							},
						},
						{
							RpcMethod: "Votes",
							Use:       "votes [proposal-id]",
							Short:     "Query votes on a proposal",
							Long: strings.TrimSpace(
								fmt.Sprintf(`Query vote details for a single proposal by its identifier.

Example:
$ %[1]s query gov votes 1
$ %[1]s query gov votes 1 --page=2 --limit=100
`,
									version.AppName,
								),
							),
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "proposal_id"},
							},
						},
						{
							RpcMethod: "Params",
							Use:       "param [param-type]",
							Short:     "Query the parameters (voting|tallying|deposit) of the governance process",
							Long: strings.TrimSpace(
								fmt.Sprintf(`Query the all the parameters for the governance process.
Example:
$ %s query gov param voting
$ %s query gov param tallying
$ %s query gov param deposit
`,
									version.AppName, version.AppName, version.AppName,
								),
							),
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "params_type"},
							},
						},
						//	TODO: Params uses custom logic to set a specific type of params \"deposits\""
						{
							RpcMethod: "Deposit",
							Use:       "deposit [proposal-id] [depositer-addr]",
							Short:     "Query details of a deposit",
							Long: strings.TrimSpace(
								fmt.Sprintf(`Query details for a single proposal deposit on a proposal by its identifier.

Example:
$ %s query gov deposit 1 cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`,
									version.AppName,
								),
							),
						},
					},
				},
			},
		},
	}
}
