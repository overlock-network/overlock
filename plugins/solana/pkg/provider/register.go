package provider

import (
	"context"
	"crypto/sha256"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	crossplanev1beta1 "github.com/overlock-network/api/go/node/overlock/crossplane/v1beta1"
	"go.uber.org/zap"
)

type RegisterProviderArgs struct {
	Name            string
	Ip              string
	Port            uint16
	Country         string
	EnvironmentType string
	Availability    bool
}

const maxInstructionSize = 10240

func Register(logger zap.SugaredLogger, grpcAddress, programId, keyPath string, provider crossplanev1beta1.MsgCreateProvider) {
	client := rpc.New(grpcAddress)
	var availability bool
	if provider.Availability == "available" {
		availability = true
	} else {
		availability = false
	}

	args := RegisterProviderArgs{
		Name:            provider.Metadata.Name,
		Ip:              provider.Ip,
		Port:            uint16(provider.Port),
		Country:         provider.CountryCode,
		EnvironmentType: provider.EnvironmentType,
		Availability:    availability,
	}

	hash := sha256.Sum256([]byte("global:register_provider"))
	discriminator := hash[:8]

	argsEncoded, err := borsh.Serialize(args)
	if err != nil {
		panic(err)
	}

	fullData := append(discriminator, argsEncoded...)

	if len(fullData) > maxInstructionSize {
		logger.Warnf("Instruction too large (%d bytes). Trimming", len(fullData))
		fullData = fullData[:maxInstructionSize]
	}

	prikey, err := solana.PrivateKeyFromSolanaKeygenFile(keyPath)
	if err != nil {
		panic(err)
	}
	payer := prikey.PublicKey()

	providerAccount := solana.NewWallet()
	providerPubkey := providerAccount.PublicKey()

	ID := solana.MustPublicKeyFromBase58(programId)

	instruction := solana.NewInstruction(
		ID,
		solana.AccountMetaSlice{
			{PublicKey: payer, IsSigner: true, IsWritable: true},
			{PublicKey: providerPubkey, IsSigner: true, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		fullData,
	)

	out, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	txBuilder := solana.NewTransactionBuilder().
		AddInstruction(instruction).
		SetFeePayer(payer).
		SetRecentBlockHash(out.Value.Blockhash)

	transaction, err := txBuilder.Build()
	if err != nil {
		panic(err)
	}

	transaction.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		switch {
		case key.Equals(payer):
			return &prikey
		case key.Equals(providerPubkey):
			return &providerAccount.PrivateKey
		default:
			return nil
		}
	})

	sig, err := client.SendTransaction(context.Background(), transaction)
	if err != nil {
		panic(err)
	}
	logger.Infof("Transaction Signature: %s", sig.String())
}
