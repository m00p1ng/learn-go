package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/m00p1ng/learn-go/nomadcoin/utils"
	"github.com/m00p1ng/learn-go/nomadcoin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

var m *mempool = &mempool{}
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address)

		if !valid {
			break
		}
	}

	return valid
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}

	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()

	return &tx
}

var ErrorNoMoney = errors.New("not enough money")
var ErrorNotValid = errors.New("Tx Invalid")

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())

	for _, uTxOut := range uTxOuts {
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		if total >= amount {
			break
		}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}

	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)

	if err != nil {
		return nil, err
	}

	m.Txs[tx.ID] = tx
	return tx, nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = nil

	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()

	m.Txs[tx.ID] = tx
}
