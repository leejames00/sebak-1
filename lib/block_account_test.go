package sebak

import (
	"testing"

	"github.com/spikeekips/sebak/lib/storage"
)

func TestSaveNewBlockAccount(t *testing.T) {
	st, _ := storage.NewTestMemoryLevelDBBackend()

	b := makeBlockAccount()
	err := b.Save(st)
	if err != nil {
		t.Errorf("failed to save BlockAccount, %v", err)
		return
	}

	exists, err := ExistBlockAccount(st, b.Address)
	if err != nil {
		t.Errorf("failed to get BlockAccount, %v", err)
		return
	}

	if !exists {
		t.Errorf("failed to get BlockAccount, does not exists")
		return
	}
}

func TestSaveExistingBlockAccount(t *testing.T) {
	st, _ := storage.NewTestMemoryLevelDBBackend()

	b := makeBlockAccount()
	b.Save(st)

	b.EnsureUpdate(100, "fake-checkpoint", b.GetBalance()+100)
	b.Save(st)

	fetched, _ := GetBlockAccount(st, b.Address)
	if b.Balance != fetched.Balance {
		t.Error("failed to update `BlockAccount.Balance`")
		return
	}
}

func TestSortMultipleBlockAccount(t *testing.T) {
	st, _ := storage.NewTestMemoryLevelDBBackend()

	var createdOrder []string
	for i := 0; i < 50; i++ {
		b := makeBlockAccount()
		b.Save(st)

		createdOrder = append(createdOrder, b.Address)
	}

	var saved []string
	iterFunc, closeFunc := GetBlockAccountAddressesByCreated(st, false)
	for {
		address, hasNext := iterFunc()
		if !hasNext {
			break
		}

		saved = append(saved, address)
	}
	closeFunc()

	for i, a := range createdOrder {
		if a != saved[i] {
			t.Error("failed to save `BlockAccount` by creation order")
			break
		}
	}
}

func TestGetSortedBlockAccounts(t *testing.T) {
	st, _ := storage.NewTestMemoryLevelDBBackend()

	var createdOrder []string
	for i := 0; i < 50; i++ {
		b := makeBlockAccount()
		b.Save(st)

		createdOrder = append(createdOrder, b.Address)
	}

	var saved []string
	iterFunc, closeFunc := GetBlockAccountsByCreated(st, false)
	for {
		ba, hasNext := iterFunc()
		if !hasNext {
			break
		}

		saved = append(saved, ba.Address)
	}
	closeFunc()

	for i, a := range createdOrder {
		if a != saved[i] {
			t.Error("failed to save `BlockAccount` by creation order")
			break
		}
	}
}