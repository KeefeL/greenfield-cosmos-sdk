package group_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/group"
)

var (
	admin   = sdk.AccAddress("admin")
	member1 = sdk.AccAddress("member1")
	member2 = sdk.AccAddress("member2")
	member3 = sdk.AccAddress("member3")
	member4 = sdk.AccAddress("member4")
	member5 = sdk.AccAddress("member5")
)

func TestMsgCreateGroup(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgCreateGroup
		expErr bool
		errMsg string
	}{
		{
			"invalid admin address",
			&group.MsgCreateGroup{
				Admin: "admin",
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"invalid member address",
			&group.MsgCreateGroup{
				Admin: admin.String(),
				Members: []group.MemberRequest{
					{
						Address: "invalid address",
					},
				},
			},
			true,
			"address: invalid address hex length",
		},
		{
			"negitive member's weight not allowed",
			&group.MsgCreateGroup{
				Admin: admin.String(),
				Members: []group.MemberRequest{
					{
						Address: member1.String(),
						Weight:  "-1",
					},
				},
			},
			true,
			"expected a non-negative decimal",
		},
		{
			"zero member's weight not allowed",
			&group.MsgCreateGroup{
				Admin: admin.String(),
				Members: []group.MemberRequest{
					{
						Address: member1.String(),
						Weight:  "0",
					},
				},
			},
			true,
			"expected a positive decimal",
		},
		{
			"duplicate member not allowed",
			&group.MsgCreateGroup{
				Admin: admin.String(),
				Members: []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				},
			},
			true,
			"duplicate value",
		},
		{
			"valid test case with single member",
			&group.MsgCreateGroup{
				Admin: admin.String(),
				Members: []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				},
			},
			false,
			"",
		},
		{
			"minimum fields",
			&group.MsgCreateGroup{
				Admin:   admin.String(),
				Members: []group.MemberRequest{},
			},
			false,
			"",
		},
		{
			"valid test case with multiple members",
			&group.MsgCreateGroup{
				Admin: admin.String(),
				Members: []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
					{
						Address:  member2.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				},
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.Type(), sdk.MsgTypeURL(&group.MsgCreateGroup{}))
			}
		})
	}
}

func TestMsgUpdateGroupAdmin(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgUpdateGroupAdmin
		expErr bool
		errMsg string
	}{
		{
			"empty group id",
			&group.MsgUpdateGroupAdmin{
				Admin:    admin.String(),
				NewAdmin: member1.String(),
			},
			true,
			"group id: value is empty",
		},
		{
			"admin: invalid address",
			&group.MsgUpdateGroupAdmin{
				GroupId: 1,
				Admin:   "admin",
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"new admin: invalid address",
			&group.MsgUpdateGroupAdmin{
				GroupId:  1,
				Admin:    admin.String(),
				NewAdmin: "new-admin",
			},
			true,
			"new admin: invalid address hex length",
		},
		{
			"admin & new admin is same",
			&group.MsgUpdateGroupAdmin{
				GroupId:  1,
				Admin:    admin.String(),
				NewAdmin: admin.String(),
			},
			true,
			"new and old admin are the same",
		},
		{
			"valid case",
			&group.MsgUpdateGroupAdmin{
				GroupId:  1,
				Admin:    admin.String(),
				NewAdmin: member1.String(),
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.Type(), sdk.MsgTypeURL(&group.MsgUpdateGroupAdmin{}))
			}
		})
	}
}

func TestMsgUpdateGroupMetadata(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgUpdateGroupMetadata
		expErr bool
		errMsg string
	}{
		{
			"empty group id",
			&group.MsgUpdateGroupMetadata{
				Admin: admin.String(),
			},
			true,
			"group id: value is empty",
		},
		{
			"admin: invalid address",
			&group.MsgUpdateGroupMetadata{
				GroupId: 1,
				Admin:   "admin",
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"valid test",
			&group.MsgUpdateGroupMetadata{
				GroupId:  1,
				Admin:    admin.String(),
				Metadata: "metadata",
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.Type(), sdk.MsgTypeURL(&group.MsgUpdateGroupMetadata{}))
			}
		})
	}
}

func TestMsgUpdateGroupMembers(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgUpdateGroupMembers
		expErr bool
		errMsg string
	}{
		{
			"empty group id",
			&group.MsgUpdateGroupMembers{},
			true,
			"group id: value is empty",
		},
		{
			"admin: invalid address",
			&group.MsgUpdateGroupMembers{
				GroupId: 1,
				Admin:   "admin",
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"empty member list",
			&group.MsgUpdateGroupMembers{
				GroupId:       1,
				Admin:         admin.String(),
				MemberUpdates: []group.MemberRequest{},
			},
			true,
			"member updates: value is empty",
		},
		{
			"valid test",
			&group.MsgUpdateGroupMembers{
				GroupId: 1,
				Admin:   admin.String(),
				MemberUpdates: []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				},
			},
			false,
			"",
		},
		{
			"valid test with zero weight",
			&group.MsgUpdateGroupMembers{
				GroupId: 1,
				Admin:   admin.String(),
				MemberUpdates: []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "0",
						Metadata: "metadata",
					},
				},
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.Type(), sdk.MsgTypeURL(&group.MsgUpdateGroupMembers{}))
			}
		})
	}
}

func TestMsgCreateGroupWithPolicy(t *testing.T) {
	testCases := []struct {
		name   string
		msg    func() *group.MsgCreateGroupWithPolicy
		expErr bool
		errMsg string
	}{
		{
			"invalid admin address",
			func() *group.MsgCreateGroupWithPolicy {
				admin := "admin"
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				members := []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				}
				req, err := group.NewMsgCreateGroupWithPolicy(admin, members, "group_metadata", "group_policy_metadata", false, policy)
				require.NoError(t, err)
				return req
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"invalid member address",
			func() *group.MsgCreateGroupWithPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				members := []group.MemberRequest{
					{
						Address:  "invalid_address",
						Weight:   "1",
						Metadata: "metadata",
					},
				}
				req, err := group.NewMsgCreateGroupWithPolicy(admin.String(), members, "group_metadata", "group_policy_metadata", false, policy)
				require.NoError(t, err)
				return req
			},
			true,
			"address: invalid address hex length",
		},
		{
			"negative member's weight not allowed",
			func() *group.MsgCreateGroupWithPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				members := []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "-1",
						Metadata: "metadata",
					},
				}
				req, err := group.NewMsgCreateGroupWithPolicy(admin.String(), members, "group_metadata", "group_policy_metadata", false, policy)
				require.NoError(t, err)
				return req
			},
			true,
			"expected a non-negative decimal",
		},
		{
			"zero member's weight not allowed",
			func() *group.MsgCreateGroupWithPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				members := []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "0",
						Metadata: "metadata",
					},
				}
				req, err := group.NewMsgCreateGroupWithPolicy(admin.String(), members, "group_metadata", "group_policy_metadata", false, policy)
				require.NoError(t, err)
				return req
			},
			true,
			"expected a positive decimal",
		},
		{
			"duplicate member not allowed",
			func() *group.MsgCreateGroupWithPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				members := []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				}
				req, err := group.NewMsgCreateGroupWithPolicy(admin.String(), members, "group_metadata", "group_policy_metadata", false, policy)
				require.NoError(t, err)
				return req
			},
			true,
			"duplicate value",
		},
		{
			"invalid threshold policy",
			func() *group.MsgCreateGroupWithPolicy {
				policy := group.NewThresholdDecisionPolicy("-1", time.Second, 0)
				members := []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				}
				req, err := group.NewMsgCreateGroupWithPolicy(admin.String(), members, "group_metadata", "group_policy_metadata", false, policy)
				require.NoError(t, err)
				return req
			},
			true,
			"expected a positive decimal",
		},
		{
			"valid test case with single member",
			func() *group.MsgCreateGroupWithPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				members := []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				}
				req, err := group.NewMsgCreateGroupWithPolicy(admin.String(), members, "group_metadata", "group_policy_metadata", false, policy)
				require.NoError(t, err)
				return req
			},
			false,
			"",
		},
		{
			"valid test case with multiple members",
			func() *group.MsgCreateGroupWithPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				members := []group.MemberRequest{
					{
						Address:  member1.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
					{
						Address:  member2.String(),
						Weight:   "1",
						Metadata: "metadata",
					},
				}
				req, err := group.NewMsgCreateGroupWithPolicy(admin.String(), members, "group_metadata", "group_policy_metadata", false, policy)
				require.NoError(t, err)
				return req
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg()
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgCreateGroupWithPolicy{}))
			}
		})
	}
}

func TestMsgCreateGroupPolicy(t *testing.T) {
	testCases := []struct {
		name   string
		msg    func() *group.MsgCreateGroupPolicy
		expErr bool
		errMsg string
	}{
		{
			"empty group id",
			func() *group.MsgCreateGroupPolicy {
				return &group.MsgCreateGroupPolicy{
					Admin: admin.String(),
				}
			},
			true,
			"group id: value is empty",
		},
		{
			"admin: invalid address",
			func() *group.MsgCreateGroupPolicy {
				return &group.MsgCreateGroupPolicy{
					Admin:   "admin",
					GroupId: 1,
				}
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"invalid threshold policy",
			func() *group.MsgCreateGroupPolicy {
				policy := group.NewThresholdDecisionPolicy("-1", time.Second, 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", policy)
				require.NoError(t, err)
				return req
			},
			true,
			"expected a positive decimal",
		},
		{
			"invalid voting period",
			func() *group.MsgCreateGroupPolicy {
				policy := group.NewThresholdDecisionPolicy("-1", time.Duration(0), 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", policy)
				require.NoError(t, err)
				return req
			},
			true,
			"expected a positive decimal",
		},
		{
			"invalid execution period",
			func() *group.MsgCreateGroupPolicy {
				policy := group.NewThresholdDecisionPolicy("-1", time.Minute, 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", policy)
				require.NoError(t, err)
				return req
			},
			true,
			"expected a positive decimal",
		},
		{
			"valid test case, only voting period",
			func() *group.MsgCreateGroupPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", policy)
				require.NoError(t, err)
				return req
			},
			false,
			"",
		},
		{
			"valid test case, voting and execution, empty min exec period",
			func() *group.MsgCreateGroupPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", policy)
				require.NoError(t, err)
				return req
			},
			false,
			"",
		},
		{
			"valid test case, voting and execution, non-empty min exec period",
			func() *group.MsgCreateGroupPolicy {
				policy := group.NewThresholdDecisionPolicy("1", time.Second, time.Minute)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", policy)
				require.NoError(t, err)
				return req
			},
			false,
			"",
		},
		{
			"invalid percentage decision policy with zero value",
			func() *group.MsgCreateGroupPolicy {
				percentagePolicy := group.NewPercentageDecisionPolicy("0", time.Second, 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", percentagePolicy)
				require.NoError(t, err)
				return req
			},
			true,
			"expected a positive decimal",
		},
		{
			"invalid percentage decision policy with negative value",
			func() *group.MsgCreateGroupPolicy {
				percentagePolicy := group.NewPercentageDecisionPolicy("-0.2", time.Second, 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", percentagePolicy)
				require.NoError(t, err)
				return req
			},
			true,
			"expected a positive decimal",
		},
		{
			"invalid percentage decision policy with value greater than 1",
			func() *group.MsgCreateGroupPolicy {
				percentagePolicy := group.NewPercentageDecisionPolicy("2", time.Second, 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", percentagePolicy)
				require.NoError(t, err)
				return req
			},
			true,
			"percentage must be > 0 and <= 1",
		},
		{
			"valid test case with percentage decision policy",
			func() *group.MsgCreateGroupPolicy {
				percentagePolicy := group.NewPercentageDecisionPolicy("0.5", time.Second, 0)
				req, err := group.NewMsgCreateGroupPolicy(admin, 1, "metadata", percentagePolicy)
				require.NoError(t, err)
				return req
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg()
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgCreateGroupPolicy{}))
			}
		})
	}
}

func TestMsgUpdateGroupPolicyDecisionPolicy(t *testing.T) {
	validPolicy := group.NewThresholdDecisionPolicy("1", time.Second, 0)
	msg1, err := group.NewMsgUpdateGroupPolicyDecisionPolicy(admin, member1, validPolicy)
	require.NoError(t, err)

	invalidPolicy := group.NewThresholdDecisionPolicy("-1", time.Second, 0)
	msg2, err := group.NewMsgUpdateGroupPolicyDecisionPolicy(admin, member2, invalidPolicy)
	require.NoError(t, err)

	validPercentagePolicy := group.NewPercentageDecisionPolicy("0.7", time.Second, 0)
	msg3, err := group.NewMsgUpdateGroupPolicyDecisionPolicy(admin, member3, validPercentagePolicy)
	require.NoError(t, err)

	invalidPercentagePolicy := group.NewPercentageDecisionPolicy("-0.1", time.Second, 0)
	msg4, err := group.NewMsgUpdateGroupPolicyDecisionPolicy(admin, member4, invalidPercentagePolicy)
	require.NoError(t, err)

	invalidPercentagePolicy2 := group.NewPercentageDecisionPolicy("2", time.Second, 0)
	msg5, err := group.NewMsgUpdateGroupPolicyDecisionPolicy(admin, member5, invalidPercentagePolicy2)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		msg    *group.MsgUpdateGroupPolicyDecisionPolicy
		expErr bool
		errMsg string
	}{
		{
			"admin: invalid address",
			&group.MsgUpdateGroupPolicyDecisionPolicy{
				Admin: "admin",
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"group policy: invalid address",
			&group.MsgUpdateGroupPolicyDecisionPolicy{
				Admin:              admin.String(),
				GroupPolicyAddress: "address",
			},
			true,
			"group policy: invalid address hex length",
		},
		{
			"group policy: invalid address",
			&group.MsgUpdateGroupPolicyDecisionPolicy{
				Admin:              admin.String(),
				GroupPolicyAddress: "address",
			},
			true,
			"group policy: invalid address hex length",
		},
		{
			"invalid decision policy",
			msg2,
			true,
			"decision policy: threshold: expected a positive decimal",
		},
		{
			"valid decision policy",
			msg1,
			false,
			"",
		},
		{
			"valid percentage decision policy",
			msg3,
			false,
			"",
		},
		{
			"invalid percentage decision policy with negative value",
			msg4,
			true,
			"decision policy: percentage threshold: expected a positive decimal",
		},
		{
			"invalid percentage decision policy with value greater than 1",
			msg5,
			true,
			"decision policy: percentage must be > 0 and <= 1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgUpdateGroupPolicyDecisionPolicy{}))
			}
		})
	}
}

func TestMsgUpdateGroupPolicyAdmin(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgUpdateGroupPolicyAdmin
		expErr bool
		errMsg string
	}{
		{
			"admin: invalid address",
			&group.MsgUpdateGroupPolicyAdmin{
				Admin: "admin",
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"policy address: invalid address",
			&group.MsgUpdateGroupPolicyAdmin{
				Admin:              admin.String(),
				NewAdmin:           member1.String(),
				GroupPolicyAddress: "address",
			},
			true,
			"group policy: invalid address hex length",
		},
		{
			"new admin: invalid address",
			&group.MsgUpdateGroupPolicyAdmin{
				Admin:              admin.String(),
				GroupPolicyAddress: admin.String(),
				NewAdmin:           "new-admin",
			},
			true,
			"new admin: invalid address hex length",
		},
		{
			"same old and new admin",
			&group.MsgUpdateGroupPolicyAdmin{
				Admin:              admin.String(),
				GroupPolicyAddress: admin.String(),
				NewAdmin:           admin.String(),
			},
			true,
			"new and old admin are same",
		},
		{
			"valid test",
			&group.MsgUpdateGroupPolicyAdmin{
				Admin:              admin.String(),
				GroupPolicyAddress: admin.String(),
				NewAdmin:           member1.String(),
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgUpdateGroupPolicyAdmin{}))
			}
		})
	}
}

func TestMsgUpdateGroupPolicyMetadata(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgUpdateGroupPolicyMetadata
		expErr bool
		errMsg string
	}{
		{
			"admin: invalid address",
			&group.MsgUpdateGroupPolicyMetadata{
				Admin: "admin",
			},
			true,
			"admin: invalid address hex length",
		},
		{
			"group policy address: invalid address",
			&group.MsgUpdateGroupPolicyMetadata{
				Admin:              admin.String(),
				GroupPolicyAddress: "address",
			},
			true,
			"group policy: invalid address hex length",
		},
		{
			"valid testcase",
			&group.MsgUpdateGroupPolicyMetadata{
				Admin:              admin.String(),
				GroupPolicyAddress: member1.String(),
				Metadata:           "metadata",
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgUpdateGroupPolicyMetadata{}))
			}
		})
	}
}

func TestMsgSubmitProposal(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgSubmitProposal
		expErr bool
		errMsg string
	}{
		{
			"invalid group policy address",
			&group.MsgSubmitProposal{
				GroupPolicyAddress: "address",
			},
			true,
			"group policy: invalid address hex length",
		},
		{
			"proposers required",
			&group.MsgSubmitProposal{
				GroupPolicyAddress: admin.String(),
			},
			true,
			"proposers: value is empty",
		},
		{
			"valid testcase",
			&group.MsgSubmitProposal{
				GroupPolicyAddress: admin.String(),
				Proposers:          []string{member1.String(), member2.String()},
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgSubmitProposal{}))
			}
		})
	}
}

func TestMsgVote(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgVote
		expErr bool
		errMsg string
	}{
		{
			"invalid voter address",
			&group.MsgVote{
				Voter: "voter",
			},
			true,
			"voter: invalid address hex length",
		},
		{
			"proposal id is required",
			&group.MsgVote{
				Voter: member1.String(),
			},
			true,
			"proposal id: value is empty",
		},
		{
			"unspecified vote option",
			&group.MsgVote{
				Voter:      member1.String(),
				ProposalId: 1,
			},
			true,
			"vote option: value is empty",
		},
		{
			"valid test case",
			&group.MsgVote{
				Voter:      member1.String(),
				ProposalId: 1,
				Option:     group.VOTE_OPTION_YES,
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgVote{}))
			}
		})
	}
}

func TestMsgWithdrawProposal(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgWithdrawProposal
		expErr bool
		errMsg string
	}{
		{
			"invalid address",
			&group.MsgWithdrawProposal{
				Address: "address",
			},
			true,
			"invalid address hex length",
		},
		{
			"proposal id is required",
			&group.MsgWithdrawProposal{
				Address: member1.String(),
			},
			true,
			"proposal id: value is empty",
		},
		{
			"valid msg",
			&group.MsgWithdrawProposal{
				Address:    member1.String(),
				ProposalId: 1,
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgWithdrawProposal{}))
			}
		})
	}
}

func TestMsgExec(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgExec
		expErr bool
		errMsg string
	}{
		{
			"invalid signer address",
			&group.MsgExec{
				Executor: "signer",
			},
			true,
			"signer: invalid address hex length",
		},
		{
			"proposal is required",
			&group.MsgExec{
				Executor: admin.String(),
			},
			true,
			"proposal id: value is empty",
		},
		{
			"valid testcase",
			&group.MsgExec{
				Executor:   admin.String(),
				ProposalId: 1,
			},
			false,
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgExec{}))
			}
		})
	}
}

func TestMsgLeaveGroup(t *testing.T) {
	testCases := []struct {
		name   string
		msg    *group.MsgLeaveGroup
		expErr bool
		errMsg string
	}{
		{
			"invalid group member address",
			&group.MsgLeaveGroup{
				Address: "member",
			},
			true,
			"group member: invalid address hex length",
		},
		{
			"group id is required",
			&group.MsgLeaveGroup{
				Address: admin.String(),
			},
			true,
			"group-id: value is empty",
		},
		{
			"valid testcase",
			&group.MsgLeaveGroup{
				Address: admin.String(),
				GroupId: 1,
			},
			false,
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			err := msg.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, msg.Type(), sdk.MsgTypeURL(&group.MsgLeaveGroup{}))
			}
		})
	}
}
