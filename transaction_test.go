package flow_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/crypto"
)

type MockSigner flow.AccountKey

func (s MockSigner) Sign(crypto.Signable) ([]byte, error) {
	return []byte{uint8(s.Index)}, nil
}

func ExampleTransaction() {
	// Mock user accounts

	adrianLaptopKey := flow.AccountKey{
		Index:          3,
		SequenceNumber: 42,
	}

	adrianPhoneKey := flow.AccountKey{Index: 2}

	adrian := flow.Account{
		Address: flow.HexToAddress("01"),
		Keys: []flow.AccountKey{
			adrianLaptopKey,
			adrianPhoneKey,
		},
	}

	blaineHardwareKey := flow.AccountKey{Index: 7}

	blaine := flow.Account{
		Address: flow.HexToAddress("02"),
		Keys: []flow.AccountKey{
			blaineHardwareKey,
		},
	}

	// Transaction preparation

	tx := flow.NewTransaction().
		SetScript([]byte(`transaction { execute { log("Hello, World!") } }`)).
		SetReferenceBlockID(flow.Identifier{0x01, 0x02}).
		SetGasLimit(42).
		SetProposalKey(adrian.Address, adrianLaptopKey.Index, adrianLaptopKey.SequenceNumber).
		SetPayer(blaine.Address, blaineHardwareKey.Index).
		AddAuthorizer(adrian.Address, adrianLaptopKey.Index, adrianPhoneKey.Index)

	fmt.Println("Signers:")
	for _, signer := range tx.Signers() {
		fmt.Printf(
			"Address: %s, Roles: %s, Key Indices: %d\n",
			signer.Address,
			signer.Roles,
			signer.KeyIndices,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (before signing): %x\n\n", tx.ID())

	// Signing

	err := tx.SignPayload(adrian.Address, adrianLaptopKey.Index, MockSigner(adrianLaptopKey))
	if err != nil {
		panic(err)
	}

	err = tx.SignPayload(adrian.Address, adrianPhoneKey.Index, MockSigner(adrianPhoneKey))
	if err != nil {
		panic(err)
	}

	err = tx.SignContainer(blaine.Address, blaineHardwareKey.Index, MockSigner(blaineHardwareKey))
	if err != nil {
		panic(err)
	}

	fmt.Println("Signatures:")
	for _, sig := range tx.Signatures {
		fmt.Printf(
			"Kind: %s, Address: %s, Key Index: %d, Signature: %x\n",
			sig.Kind,
			sig.Address,
			sig.KeyIndex,
			sig.Signature,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (after signing): %x\n", tx.ID())

	// Output:
	// Signers:
	// Address: 0000000000000000000000000000000000000001, Roles: [PROPOSER AUTHORIZER], Key Indices: [2 3]
	// Address: 0000000000000000000000000000000000000002, Roles: [PAYER], Key Indices: [7]
	//
	// Transaction ID (before signing): 4cd86595c7dc854b371644060c1b4cbc478726b7e3c8be2176353c169e1a76d3
	//
	// Signatures:
	// Kind: PAYLOAD, Address: 0000000000000000000000000000000000000001, Key Index: 3, Signature: 03
	// Kind: PAYLOAD, Address: 0000000000000000000000000000000000000001, Key Index: 2, Signature: 02
	// Kind: CONTAINER, Address: 0000000000000000000000000000000000000002, Key Index: 7, Signature: 07
	//
	// Transaction ID (after signing): 63271c5cb5429bcabbb3fd0f174afd1d22ca4c2e5fb237cf940ce1c61e2176f3
}

var (
	AddressA flow.Address
	AddressB flow.Address
	AddressC flow.Address
	AddressD flow.Address
	AddressE flow.Address

	RolesProposerPayerAuthorizer []flow.SignerRole
	RolesProposerPayer           []flow.SignerRole
	RolesProposerAuthorizer      []flow.SignerRole
	RolesPayerAuthorizer         []flow.SignerRole
	RolesProposer                []flow.SignerRole
	RolesPayer                   []flow.SignerRole
	RolesAuthorizer              []flow.SignerRole
)

func init() {
	AddressA = flow.HexToAddress("01")
	AddressB = flow.HexToAddress("02")
	AddressC = flow.HexToAddress("03")
	AddressD = flow.HexToAddress("04")
	AddressE = flow.HexToAddress("05")

	RolesProposerPayerAuthorizer = []flow.SignerRole{flow.SignerRoleProposer, flow.SignerRolePayer, flow.SignerRoleAuthorizer}
	RolesProposerPayer = []flow.SignerRole{flow.SignerRoleProposer, flow.SignerRolePayer}
	RolesProposerAuthorizer = []flow.SignerRole{flow.SignerRoleProposer, flow.SignerRoleAuthorizer}
	RolesPayerAuthorizer = []flow.SignerRole{flow.SignerRolePayer, flow.SignerRoleAuthorizer}
	RolesProposer = []flow.SignerRole{flow.SignerRoleProposer}
	RolesPayer = []flow.SignerRole{flow.SignerRolePayer}
	RolesAuthorizer = []flow.SignerRole{flow.SignerRoleAuthorizer}
}

func TestTransaction_Signers_SeparateSigners(t *testing.T) {
	t.Run("No authorizers", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressB, 1).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIndices)
	})

	t.Run("With authorizer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1).
			SetPayer(AddressC, 1).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIndices)

		assert.Equal(t, AddressC, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIndices)
	})
}

func TestTransaction_Signers_DeclarationOrder(t *testing.T) {
	t.Run("Payer before proposer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetPayer(AddressB, 1).
			SetProposalKey(AddressA, 1, 42).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIndices)
	})

	t.Run("Authorizer before proposer", func(t *testing.T) {
		signers := flow.NewTransaction().
			AddAuthorizer(AddressB, 1).
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressC, 1).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIndices)

		assert.Equal(t, AddressC, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIndices)
	})

	t.Run("Authorizer after payer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressC, 1).
			AddAuthorizer(AddressB, 1).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIndices)

		assert.Equal(t, AddressC, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIndices)
	})

	t.Run("Authorizer before and after payer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1).
			SetPayer(AddressD, 1).
			AddAuthorizer(AddressC, 1).
			Signers()

		require.Len(t, signers, 4)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIndices)

		assert.Equal(t, AddressC, signers[2].Address)
		assert.Equal(t, RolesAuthorizer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIndices)

		assert.Equal(t, AddressD, signers[3].Address)
		assert.Equal(t, RolesPayer, signers[3].Roles)
		assert.Equal(t, []int{1}, signers[3].KeyIndices)
	})
}

func TestTransaction_Signers_KeysOutOfOrder(t *testing.T) {
	signers := flow.NewTransaction().
		SetProposalKey(AddressA, 1, 42).
		SetPayer(AddressA, 4, 2, 1, 3).
		Signers()

	require.Len(t, signers, 1)

	assert.Equal(t, AddressA, signers[0].Address)
	assert.Equal(t, RolesProposerPayer, signers[0].Roles)
	assert.Equal(t, []int{1, 2, 3, 4}, signers[0].KeyIndices)
}

func TestTransaction_Signers_MultipleAuthorizers(t *testing.T) {
	signers := flow.NewTransaction().
		SetProposalKey(AddressA, 1, 42).
		AddAuthorizer(AddressB, 1).
		AddAuthorizer(AddressC, 2).
		AddAuthorizer(AddressD, 3).
		SetPayer(AddressE, 1).
		Signers()

	require.Len(t, signers, 5)

	assert.Equal(t, AddressB, signers[1].Address)
	assert.Equal(t, RolesAuthorizer, signers[1].Roles)
	assert.Equal(t, []int{1}, signers[1].KeyIndices)

	assert.Equal(t, AddressC, signers[2].Address)
	assert.Equal(t, RolesAuthorizer, signers[2].Roles)
	assert.Equal(t, []int{2}, signers[2].KeyIndices)

	assert.Equal(t, AddressD, signers[3].Address)
	assert.Equal(t, RolesAuthorizer, signers[3].Roles)
	assert.Equal(t, []int{3}, signers[3].KeyIndices)
}

func TestTransaction_Signers_ProposerPayerAuthorizerSameAddress(t *testing.T) {
	t.Run("Single key", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 1).
			AddAuthorizer(AddressA, 1).
			Signers()

		require.Len(t, signers, 1)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerPayerAuthorizer, signers[0].Roles)
		assert.Equal(
			t,
			&flow.ProposalKey{
				Address:        AddressA,
				KeyIndex:       1,
				SequenceNumber: 42,
			},
			signers[0].ProposalKey,
		)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)
	})

	t.Run("Identical key-sets", func(t *testing.T) {
		// All key-sets contain the elements [1, 2]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 1, 2).
			AddAuthorizer(AddressA, 1, 2).
			Signers()

		require.Len(t, signers, 1)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerPayerAuthorizer, signers[0].Roles)
		assert.Equal(t, []int{1, 2}, signers[0].KeyIndices)
	})

	t.Run("Subset of payer key-set", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [1]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 1, 2).
			AddAuthorizer(AddressA, 1).
			Signers()

		require.Len(t, signers, 1)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerPayerAuthorizer, signers[0].Roles)
		assert.Equal(t, []int{1, 2}, signers[0].KeyIndices)
	})
}

func TestTransaction_Signers_ProposerPayerSameAddress(t *testing.T) {
	t.Run("No authorizers", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 1, 2).
			Signers()

		require.Len(t, signers, 1)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerPayer, signers[0].Roles)
		assert.Equal(t, []int{1, 2}, signers[0].KeyIndices)
	})

	t.Run("With authorizer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1).
			SetPayer(AddressA, 1, 2).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressB, signers[0].Address)
		assert.Equal(t, RolesAuthorizer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressA, signers[1].Address)
		assert.Equal(t, RolesProposerPayer, signers[1].Roles)
		assert.Equal(t, []int{1, 2}, signers[1].KeyIndices)
	})

	t.Run("Disjoint key-sets", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 2, 3).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressA, signers[1].Address)
		assert.Equal(t, RolesPayer, signers[1].Roles)
		assert.Equal(t, []int{2, 3}, signers[1].KeyIndices)
	})
}

func TestTransaction_Signers_PayerAuthorizerSameAddress(t *testing.T) {
	t.Run("Identical key-sets", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [1, 2]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1, 2).
			SetPayer(AddressB, 1, 2).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayerAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1, 2}, signers[1].KeyIndices)
	})

	t.Run("Subset of payer key-set", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [1]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1).
			SetPayer(AddressB, 1, 2).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayerAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1, 2}, signers[1].KeyIndices)
	})

	t.Run("Disjoint key-sets", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [3, 4]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 3, 4).
			SetPayer(AddressB, 1, 2).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{3, 4}, signers[1].KeyIndices)

		assert.Equal(t, AddressB, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1, 2}, signers[2].KeyIndices)
	})

	t.Run("Overlapping key-sets", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [2, 3]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 2, 3).
			SetPayer(AddressB, 1, 2).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{2, 3}, signers[1].KeyIndices)

		assert.Equal(t, AddressB, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1, 2}, signers[2].KeyIndices)
	})
}

func TestTransaction_Signers_ProposerAuthorizerSameAddress(t *testing.T) {
	t.Run("Overlapping key-sets", func(t *testing.T) {
		// Proposal key: 1
		// Authorizer key-set: [1, 2]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressA, 1, 2).
			SetPayer(AddressB, 1).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerAuthorizer, signers[0].Roles)
		assert.Equal(t, []int{1, 2}, signers[0].KeyIndices)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIndices)
	})

	t.Run("Disjoint key-sets", func(t *testing.T) {
		// Proposal key: 1
		// Authorizer key-set: [2]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressA, 2).
			SetPayer(AddressB, 1).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIndices)

		assert.Equal(t, AddressA, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{2}, signers[1].KeyIndices)

		assert.Equal(t, AddressB, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIndices)
	})
}
