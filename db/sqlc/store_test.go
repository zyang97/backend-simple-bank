package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run a concurrent transfer transactions.
	n := 5
	amount := int64(10)

	results := make(chan TransferTxResult)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// Check results.
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer.
		transfer := result.Transfer
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreateAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check entries.
		from_entry := result.FromEntry
		require.Equal(t, from_entry.AccountID, account1.ID)
		require.Equal(t, from_entry.Amount, -amount)
		require.NotZero(t, from_entry.ID)
		require.NotZero(t, from_entry.CreateAt)

		_, err = store.GetEntry(context.Background(), from_entry.ID)
		require.NoError(t, err)

		to_entry := result.ToEntry
		require.Equal(t, to_entry.AccountID, account2.ID)
		require.Equal(t, to_entry.Amount, -amount)
		require.NotZero(t, to_entry.ID)
		require.NotZero(t, to_entry.CreateAt)

		_, err = store.GetEntry(context.Background(), to_entry.ID)
		require.NoError(t, err)

		// TODO: check accounts' balance.

	}
}
