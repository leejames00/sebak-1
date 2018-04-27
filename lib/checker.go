package sebak

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/spikeekips/sebak/lib/error"
	"github.com/stellar/go/keypair"
)

func checkTransactionSource(target interface{}, args ...interface{}) error {
	if _, err := keypair.Parse(target.(Transaction).B.Source); err != nil {
		return sebak_error.ErrorBadPublicAddress
	}

	return nil
}

func checkTransactionBaseFee(target interface{}, args ...interface{}) error {
	if int64(target.(Transaction).B.Fee) < BaseFee {
		return sebak_error.ErrorInvalidFee
	}

	return nil
}

func checkTransactionOperationIsWellFormed(target interface{}, args ...interface{}) error {
	tx := target.(Transaction)
	for _, op := range tx.B.Operations {
		if ta := op.B.GetTargetAddress(); tx.B.Source == ta {
			return sebak_error.ErrorInvalidOperation
		}
		if err := op.IsWellFormed(); err != nil {
			return err
		}
	}

	return nil
}

func checkTransactionVerifySignature(target interface{}, args ...interface{}) error {
	tx := target.(Transaction)
	err := keypair.MustParse(tx.B.Source).Verify([]byte(tx.H.Hash), base58.Decode(tx.H.Signature))
	if err != nil {
		return err
	}
	return nil
}

func checkTransactionHashMatch(target interface{}, args ...interface{}) error {
	tx := target.(Transaction)
	if tx.H.Hash != tx.B.MakeHashString() {
		return sebak_error.ErrorHashDoesNotMatch
	}

	return nil
}