package jsonrpc

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/TeamFoxx2025/LadyFoxx/types"

	"github.com/stretchr/testify/assert"
)

func TestContentEndpoint(t *testing.T) {
	t.Parallel()

	t.Run("returns empty ContentResponse if tx pool has no transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Content()
		//nolint:forcetypeassert
		response := result.(ContentResponse)

		assert.True(t, mockStore.includeQueued)
		assert.Equal(t, 0, len(response.Pending))
		assert.Equal(t, 0, len(response.Queued))
	})

	t.Run("returns correct data for pending transaction", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		address1 := types.Address{0x1}
		testTx1 := newTestTransaction(2, address1)
		testTx2 := newTestDynamicFeeTransaction(3, address1)
		mockStore.pending[address1] = []*types.Transaction{testTx1, testTx2}
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Content()
		//nolint:forcetypeassert
		response := result.(ContentResponse)

		assert.Equal(t, 1, len(response.Pending))
		assert.Equal(t, 0, len(response.Queued))
		assert.Equal(t, 2, len(response.Pending[address1]))

		txData := response.Pending[address1][testTx1.Nonce]
		assert.NotNil(t, txData)
		assert.Equal(t, testTx1.Gas, uint64(txData.Gas))
		assert.Equal(t, *testTx1.GasPrice, big.Int(*txData.GasPrice))
		assert.Equal(t, (*argBig)(nil), txData.GasFeeCap)
		assert.Equal(t, (*argBig)(nil), txData.GasTipCap)
		assert.Equal(t, testTx1.To, txData.To)
		assert.Equal(t, testTx1.From, txData.From)
		assert.Equal(t, *testTx1.Value, big.Int(txData.Value))
		assert.Equal(t, testTx1.Input, []byte(txData.Input))
		assert.Equal(t, (*argUint64)(nil), txData.BlockNumber)
		assert.Equal(t, (*argUint64)(nil), txData.TxIndex)

		txData = response.Pending[address1][testTx2.Nonce]
		assert.NotNil(t, txData)
		assert.Equal(t, (argUint64)(types.DynamicFeeTx), txData.Type)
		assert.Equal(t, testTx2.Gas, uint64(txData.Gas))
		assert.Equal(t, (*argBig)(nil), txData.GasPrice)
		assert.Equal(t, *testTx2.GasFeeCap, big.Int(*txData.GasFeeCap))
		assert.Equal(t, *testTx2.GasTipCap, big.Int(*txData.GasTipCap))
		assert.Equal(t, testTx2.To, txData.To)
		assert.Equal(t, testTx2.From, txData.From)
		assert.Equal(t, *testTx2.ChainID, big.Int(*txData.ChainID))
		assert.Equal(t, *testTx2.Value, big.Int(txData.Value))
		assert.Equal(t, testTx2.Input, []byte(txData.Input))
		assert.Equal(t, (*argUint64)(nil), txData.BlockNumber)
		assert.Equal(t, (*argUint64)(nil), txData.TxIndex)
	})

	t.Run("returns correct data for queued transaction", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		address1, address2 := types.Address{0x1}, types.Address{0x2}
		testTx1 := newTestTransaction(2, address1)
		testTx2 := newTestDynamicFeeTransaction(1, address2)
		mockStore.queued[address1] = []*types.Transaction{testTx1}
		mockStore.queued[address2] = []*types.Transaction{testTx2}
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Content()
		//nolint:forcetypeassert
		response := result.(ContentResponse)

		assert.Equal(t, 0, len(response.Pending))
		assert.Equal(t, 2, len(response.Queued))
		assert.Equal(t, 1, len(response.Queued[address1]))
		assert.Equal(t, 1, len(response.Queued[address2]))

		txData := response.Queued[address1][testTx1.Nonce]
		assert.NotNil(t, txData)
		assert.Equal(t, testTx1.Gas, uint64(txData.Gas))
		assert.Equal(t, *testTx1.GasPrice, big.Int(*txData.GasPrice))
		assert.Equal(t, (*argBig)(nil), txData.GasFeeCap)
		assert.Equal(t, (*argBig)(nil), txData.GasTipCap)
		assert.Equal(t, testTx1.To, txData.To)
		assert.Equal(t, testTx1.From, txData.From)
		assert.Equal(t, *testTx1.Value, big.Int(txData.Value))
		assert.Equal(t, testTx1.Input, []byte(txData.Input))
		assert.Equal(t, (*argUint64)(nil), txData.BlockNumber)
		assert.Equal(t, (*argUint64)(nil), txData.TxIndex)

		txData = response.Queued[address2][testTx2.Nonce]
		assert.NotNil(t, txData)
		assert.Equal(t, (argUint64)(types.DynamicFeeTx), txData.Type)
		assert.Equal(t, testTx2.Gas, uint64(txData.Gas))
		assert.Equal(t, (*argBig)(nil), txData.GasPrice)
		assert.Equal(t, *testTx2.GasFeeCap, big.Int(*txData.GasFeeCap))
		assert.Equal(t, *testTx2.GasTipCap, big.Int(*txData.GasTipCap))
		assert.Equal(t, testTx2.To, txData.To)
		assert.Equal(t, testTx2.From, txData.From)
		assert.Equal(t, *testTx2.ChainID, big.Int(*txData.ChainID))
		assert.Equal(t, *testTx2.Value, big.Int(txData.Value))
		assert.Equal(t, testTx2.Input, []byte(txData.Input))
		assert.Equal(t, (*argUint64)(nil), txData.BlockNumber)
		assert.Equal(t, (*argUint64)(nil), txData.TxIndex)
	})

	t.Run("returns correct ContentResponse data for multiple transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		address1 := types.Address{0x1}
		testTx1 := newTestTransaction(2, address1)
		testTx2 := newTestTransaction(4, address1)
		testTx3 := newTestTransaction(11, address1)
		address2 := types.Address{0x2}
		testTx4 := newTestTransaction(7, address2)
		testTx5 := newTestTransaction(8, address2)
		mockStore.pending[address1] = []*types.Transaction{testTx1, testTx2}
		mockStore.pending[address2] = []*types.Transaction{testTx4}
		mockStore.queued[address1] = []*types.Transaction{testTx3}
		mockStore.queued[address2] = []*types.Transaction{testTx5}
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Content()
		//nolint:forcetypeassert
		response := result.(ContentResponse)

		assert.True(t, mockStore.includeQueued)
		assert.Equal(t, 2, len(response.Pending))
		assert.Equal(t, 2, len(response.Pending[address1]))
		assert.Equal(t, 1, len(response.Pending[address2]))
		assert.Equal(t, 2, len(response.Queued))
	})
}

func TestInspectEndpoint(t *testing.T) {
	t.Parallel()

	t.Run("returns empty InspectResponse if tx pool has no transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		mockStore.maxSlots = 1024
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Inspect()
		//nolint:forcetypeassert
		response := result.(InspectResponse)

		assert.True(t, mockStore.includeQueued)
		assert.Equal(t, 0, len(response.Pending))
		assert.Equal(t, 0, len(response.Queued))
		assert.Equal(t, uint64(0), response.CurrentCapacity)
		assert.Equal(t, mockStore.maxSlots, response.MaxCapacity)
	})

	t.Run("returns correct data for queued transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		mockStore.capacity = 1
		address1 := types.Address{0x1}
		testTx := newTestTransaction(2, address1)
		mockStore.queued[address1] = []*types.Transaction{testTx}
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Inspect()
		//nolint:forcetypeassert
		response := result.(InspectResponse)

		assert.Equal(t, 0, len(response.Pending))
		assert.Equal(t, 1, len(response.Queued))
		assert.Equal(t, uint64(1), response.CurrentCapacity)
		transactionInfo := response.Queued[testTx.From.String()]
		assert.NotNil(t, transactionInfo)
		assert.NotNil(t, transactionInfo[strconv.FormatUint(testTx.Nonce, 10)])
	})

	t.Run("returns correct data for pending transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		mockStore.capacity = 2
		address1 := types.Address{0x1}
		testTx := newTestTransaction(2, address1)
		testTx2 := newTestTransaction(3, address1)
		mockStore.pending[address1] = []*types.Transaction{testTx, testTx2}
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Inspect()
		//nolint:forcetypeassert
		response := result.(InspectResponse)

		assert.Equal(t, 1, len(response.Pending))
		assert.Equal(t, 0, len(response.Queued))
		assert.Equal(t, uint64(2), response.CurrentCapacity)
		transactionInfo := response.Pending[testTx.From.String()]
		assert.NotNil(t, transactionInfo)
		assert.NotNil(t, transactionInfo[strconv.FormatUint(testTx.Nonce, 10)])
		assert.NotNil(t, transactionInfo[strconv.FormatUint(testTx2.Nonce, 10)])
	})
}

func TestStatusEndpoint(t *testing.T) {
	t.Parallel()

	t.Run("returns empty StatusResponse if tx pool has no transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Status()
		//nolint:forcetypeassert
		response := result.(StatusResponse)

		assert.Equal(t, uint64(0), response.Pending)
		assert.Equal(t, uint64(0), response.Queued)
	})

	t.Run("returns correct count of pending/queued transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		address1 := types.Address{0x1}
		testTx1 := newTestTransaction(2, address1)
		testTx2 := newTestTransaction(4, address1)
		testTx3 := newTestTransaction(11, address1)
		address2 := types.Address{0x2}
		testTx4 := newTestTransaction(7, address2)
		testTx5 := newTestTransaction(8, address2)
		mockStore.pending[address1] = []*types.Transaction{testTx1, testTx2}
		mockStore.pending[address2] = []*types.Transaction{testTx4}
		mockStore.queued[address1] = []*types.Transaction{testTx3}
		mockStore.queued[address2] = []*types.Transaction{testTx5}
		txPoolEndpoint := &TxPool{mockStore}

		result, _ := txPoolEndpoint.Status()
		//nolint:forcetypeassert
		response := result.(StatusResponse)

		assert.Equal(t, uint64(3), response.Pending)
		assert.Equal(t, uint64(2), response.Queued)
	})
}

type mockTxPoolStore struct {
	pending       map[types.Address][]*types.Transaction
	queued        map[types.Address][]*types.Transaction
	capacity      uint64
	maxSlots      uint64
	baseFee       uint64
	includeQueued bool
}

func newMockTxPoolStore() *mockTxPoolStore {
	return &mockTxPoolStore{
		pending: make(map[types.Address][]*types.Transaction),
		queued:  make(map[types.Address][]*types.Transaction),
	}
}

func (s *mockTxPoolStore) GetTxs(inclQueued bool) (map[types.Address][]*types.Transaction, map[types.Address][]*types.Transaction) {
	s.includeQueued = inclQueued

	return s.pending, s.queued
}

func (s *mockTxPoolStore) GetCapacity() (uint64, uint64) {
	return s.capacity, s.maxSlots
}

func (s *mockTxPoolStore) GetBaseFee() uint64 {
	return s.baseFee
}

func newTestTransaction(nonce uint64, from types.Address) *types.Transaction {
	txn := &types.Transaction{
		Nonce:    nonce,
		GasPrice: big.NewInt(1),
		Gas:      nonce * 100,
		Value:    big.NewInt(200),
		Input:    []byte{0xff},
		From:     from,
		To:       &addr1,
		V:        big.NewInt(1),
		R:        big.NewInt(1),
		S:        big.NewInt(1),
	}

	txn.ComputeHash(1)

	return txn
}

func newTestDynamicFeeTransaction(nonce uint64, from types.Address) *types.Transaction {
	txn := &types.Transaction{
		Type:      types.DynamicFeeTx,
		Nonce:     nonce,
		GasTipCap: big.NewInt(2),
		GasFeeCap: big.NewInt(4),
		ChainID:   big.NewInt(100),
		Gas:       nonce * 100,
		Value:     big.NewInt(200),
		Input:     []byte{0xff},
		From:      from,
		To:        &addr1,
		V:         big.NewInt(1),
		R:         big.NewInt(1),
		S:         big.NewInt(1),
	}

	txn.ComputeHash(1)

	return txn
}
