package wallet

import (
	"testing"

	"github.com/tarality/0xTaral/messages/proto"
	"github.com/stretchr/testify/require"

	"github.com/TeamFoxx2025/LadyFoxx/bls"
	"github.com/TeamFoxx2025/LadyFoxx/consensus/polybft/signer"
)

func Test_RecoverAddressFromSignature(t *testing.T) {
	t.Parallel()

	for _, account := range []*Account{generateTestAccount(t), generateTestAccount(t), generateTestAccount(t)} {
		key := NewKey(account)
		msgNoSig := &proto.Message{
			From:    key.Address().Bytes(),
			Type:    proto.MessageType_COMMIT,
			Payload: &proto.Message_CommitData{},
		}

		msg, err := key.SignIBFTMessage(msgNoSig)
		require.NoError(t, err)

		payload, err := msgNoSig.PayloadNoSig()
		require.NoError(t, err)

		address, err := RecoverAddressFromSignature(msg.Signature, payload)
		require.NoError(t, err)
		require.Equal(t, key.Address().Bytes(), address.Bytes())
	}
}

func Test_Sign(t *testing.T) {
	t.Parallel()

	msg := []byte("some message")

	for _, account := range []*Account{generateTestAccount(t), generateTestAccount(t)} {
		key := NewKey(account)
		ser, err := key.SignWithDomain(msg, signer.DomainCheckpointManager)

		require.NoError(t, err)

		sig, err := bls.UnmarshalSignature(ser)
		require.NoError(t, err)

		require.True(t, sig.Verify(key.raw.Bls.PublicKey(), msg, signer.DomainCheckpointManager))
	}
}

func Test_String(t *testing.T) {
	t.Parallel()

	for _, account := range []*Account{generateTestAccount(t), generateTestAccount(t), generateTestAccount(t)} {
		key := NewKey(account)
		require.Equal(t, key.Address().String(), key.String())
	}
}
