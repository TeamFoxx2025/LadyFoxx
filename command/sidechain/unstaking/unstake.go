package unstaking

import (
	"fmt"
	"time"

	"github.com/TeamFoxx2025/LadyFoxx/command"
	"github.com/TeamFoxx2025/LadyFoxx/command/helper"
	"github.com/TeamFoxx2025/LadyFoxx/command/polybftsecrets"
	sidechainHelper "github.com/TeamFoxx2025/LadyFoxx/command/sidechain"
	"github.com/TeamFoxx2025/LadyFoxx/consensus/polybft/contractsapi"
	"github.com/TeamFoxx2025/LadyFoxx/contracts"
	"github.com/TeamFoxx2025/LadyFoxx/txrelayer"
	"github.com/TeamFoxx2025/LadyFoxx/types"
	"github.com/spf13/cobra"
	"github.com/umbracle/ethgo"
)

var params unstakeParams

func GetCommand() *cobra.Command {
	unstakeCmd := &cobra.Command{
		Use:     "unstake",
		Short:   "Unstakes the amount sent for validator or undelegates amount from validator",
		PreRunE: runPreRun,
		RunE:    runCommand,
	}

	helper.RegisterJSONRPCFlag(unstakeCmd)
	setFlags(unstakeCmd)

	return unstakeCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.accountDir,
		polybftsecrets.AccountDirFlag,
		"",
		polybftsecrets.AccountDirFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.accountConfig,
		polybftsecrets.AccountConfigFlag,
		"",
		polybftsecrets.AccountConfigFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.amount,
		sidechainHelper.AmountFlag,
		"",
		"amount to unstake from validator",
	)

	cmd.MarkFlagsMutuallyExclusive(polybftsecrets.AccountDirFlag, polybftsecrets.AccountConfigFlag)
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	params.jsonRPC = helper.GetJSONRPCAddress(cmd)

	return params.validateFlags()
}

func runCommand(cmd *cobra.Command, _ []string) error {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	validatorAccount, err := sidechainHelper.GetAccount(params.accountDir, params.accountConfig)
	if err != nil {
		return err
	}

	txRelayer, err := txrelayer.NewTxRelayer(txrelayer.WithIPAddress(params.jsonRPC),
		txrelayer.WithReceiptTimeout(150*time.Millisecond))
	if err != nil {
		return err
	}

	unstakeFn := &contractsapi.UnstakeValidatorSetFn{
		Amount: params.amountValue,
	}

	encoded, err := unstakeFn.EncodeAbi()
	if err != nil {
		return err
	}

	txn := &ethgo.Transaction{
		From:  validatorAccount.Ecdsa.Address(),
		Input: encoded,
		To:    (*ethgo.Address)(&contracts.ValidatorSetContract),
	}

	receipt, err := txRelayer.SendTransaction(txn, validatorAccount.Ecdsa)
	if err != nil {
		return err
	}

	if receipt.Status != uint64(types.ReceiptSuccess) {
		return fmt.Errorf("unstake transaction failed on block: %d", receipt.BlockNumber)
	}

	var (
		withdrawalRegisteredEvent contractsapi.WithdrawalRegisteredEvent
		foundLog                  bool
	)

	result := &unstakeResult{
		ValidatorAddress: validatorAccount.Ecdsa.Address().String(),
	}

	// check the logs to check for the result
	for _, log := range receipt.Logs {
		doesMatch, err := withdrawalRegisteredEvent.ParseLog(log)
		if err != nil {
			return err
		}

		if doesMatch {
			foundLog = true
			result.Amount = withdrawalRegisteredEvent.Amount.Uint64()

			break
		}
	}

	if !foundLog {
		return fmt.Errorf("could not find an appropriate log in receipt that unstake happened (withdrawal registered)")
	}

	outputter.WriteCommandResult(result)

	return nil
}
