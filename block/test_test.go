package block

import (
	"bytes"
	"math/big"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/db"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/goloop/service/txresult"
	"github.com/pkg/errors"
)

const (
	gheight           int64 = 0
	defaultValidators       = 1
)

type testChain struct {
	wallet   module.Wallet
	database db.Database
	gtx      *testTransaction
	vld      module.CommitVoteSetDecoder
}

func (c *testChain) Database() db.Database {
	return c.database
}

func (c *testChain) Wallet() module.Wallet {
	return c.wallet
}

func (c *testChain) NID() int {
	return 0
}

func (c *testChain) Genesis() []byte {
	return c.gtx.Bytes()
}

func (c *testChain) CommitVoteSetDecoder() module.CommitVoteSetDecoder {
	return c.vld
}

type testError struct {
}

func (e *testError) Error() string {
	return "testError"
}

type testTransactionEffect struct {
	WorldState     []byte
	NextValidators *testValidatorList
	LogBloom       txresult.LogBloom
}

type testReceiptData struct {
	To                 *common.Address
	CumulativeStepUsed *big.Int
	StepPrice          *big.Int
	StepUsed           *big.Int
	Status             module.Status
	SCOREAddress       *common.Address
}

type testReceipt struct {
	Data testReceiptData
}

func (r *testReceipt) Bytes() []byte {
	return codec.MustMarshalToBytes(r)
}

func (r *testReceipt) To() module.Address {
	return r.Data.To
}

func (r *testReceipt) CumulativeStepUsed() *big.Int {
	return r.Data.CumulativeStepUsed
}

func (r *testReceipt) StepPrice() *big.Int {
	return r.Data.StepPrice
}

func (r *testReceipt) StepUsed() *big.Int {
	return r.Data.StepUsed
}

func (r *testReceipt) Status() module.Status {
	return r.Data.Status
}

func (r *testReceipt) SCOREAddress() module.Address {
	return r.Data.SCOREAddress
}

func (r *testReceipt) Check(r2 module.Receipt) error {
	rct2, ok := r2.(*testReceipt)
	if !ok {
		return errors.Errorf("IncompatibleReceipt")
	}
	if !reflect.DeepEqual(&r.Data, &rct2.Data) {
		return errors.Errorf("DataIsn'tEqual")
	}
	return nil
}

func (r *testReceipt) ToJSON(int) (interface{}, error) {
	panic("not implemented")
}

type testTransactionData struct {
	Group                 module.TransactionGroup
	CreateError           *testError
	TransitionCreateError *testError
	ValidateError         *testError
	ExecuteError          *testError
	Effect                testTransactionEffect
	Receipt               testReceipt
	Nonce                 int64
}

type testTransaction struct {
	Data testTransactionData
}

var gnonce int64

func newTestTransaction() *testTransaction {
	tx := &testTransaction{}
	tx.Data.Group = module.TransactionGroupNormal
	tx.Data.Nonce = atomic.LoadInt64(&gnonce)
	atomic.AddInt64(&gnonce, 1)
	return tx
}

func (tx *testTransaction) Group() module.TransactionGroup {
	return tx.Data.Group
}

func (tx *testTransaction) ID() []byte {
	return crypto.SHA3Sum256(tx.Bytes())
}

func (tx *testTransaction) Bytes() []byte {
	return codec.MustMarshalToBytes(tx)
}

func (tx *testTransaction) Hash() []byte {
	panic("not implemented")
}

func (tx *testTransaction) Verify() error {
	return nil
}

func (tx *testTransaction) Version() int {
	return 0
}

func (tx *testTransaction) ToJSON(version int) (interface{}, error) {
	panic("not implemented")
}

type testTransactionIterator struct {
	*testTransactionList
	i int
}

func (it *testTransactionIterator) Has() bool {
	return it.i < len(it.Transactions)
}

func (it *testTransactionIterator) Next() error {
	if !it.Has() {
		return errors.Errorf("no more index=%v\n", it.i)
	}
	it.i++
	return nil
}

func (it *testTransactionIterator) Get() (module.Transaction, int, error) {
	if !it.Has() {
		return nil, -1, errors.Errorf("no more index=%v\n", it.i)
	}
	return it.Transactions[it.i], it.i, nil
}

type testTransactionList struct {
	Transactions []*testTransaction
	_effect      *testTransactionEffect
	_receipts    []*testReceipt
}

func newTestTransactionList(txs []*testTransaction) *testTransactionList {
	l := &testTransactionList{}
	l.Transactions = make([]*testTransaction, len(txs))
	copy(l.Transactions, txs)
	l.updateCache()
	return l
}

func (l *testTransactionList) effect() *testTransactionEffect {
	l.updateCache()
	return l._effect
}

func (l *testTransactionList) receipts() []*testReceipt {
	l.updateCache()
	return l._receipts
}

func (l *testTransactionList) updateCache() {
	if l._effect == nil {
		l._effect = &testTransactionEffect{}
		for _, tx := range l.Transactions {
			if tx.Data.Effect.WorldState != nil {
				l._effect.WorldState = tx.Data.Effect.WorldState
			}
			if tx.Data.Effect.NextValidators != nil {
				l._effect.NextValidators = tx.Data.Effect.NextValidators
			}
			l._effect.LogBloom.Merge(&tx.Data.Effect.LogBloom)
			l._receipts = append(l._receipts, &tx.Data.Receipt)
		}
	}
}

func (l *testTransactionList) Len() int {
	return len(l.Transactions)
}

func (l *testTransactionList) Get(i int) (module.Transaction, error) {
	if 0 <= i && i < len(l.Transactions) {
		return l.Transactions[i], nil
	}
	return nil, errors.Errorf("bad index=%v\n", i)
}

func (l *testTransactionList) Iterator() module.TransactionIterator {
	return &testTransactionIterator{l, 0}
}

func (l *testTransactionList) Hash() []byte {
	return crypto.SHA3Sum256(codec.MustMarshalToBytes(l))
}

func (l *testTransactionList) Equal(l2 module.TransactionList) bool {
	if tl, ok := l2.(*testTransactionList); ok {
		if len(l.Transactions) != len(tl.Transactions) {
			return false
		}
		for i := 0; i < len(l.Transactions); i++ {
			if !bytes.Equal(l.Transactions[i].Bytes(), tl.Transactions[i].Bytes()) {
				return false
			}
		}
		return true
	}
	it := l.Iterator()
	it2 := l2.Iterator()
	for {
		if it.Has() != it2.Has() {
			return false
		}
		if !it.Has() {
			return true
		}
		t, _, _ := it.Get()
		t2, _, _ := it2.Get()
		if !bytes.Equal(t.Bytes(), t2.Bytes()) {
			return false
		}
	}
}

func (l *testTransactionList) Flush() error {
	panic("not implemented")
}

type transitionStep int

const (
	transitionStepUnexecuted transitionStep = iota
	transitionStepExecuting
	transitionStepSucceed
	transitionStepFailed
)

type testTransitionResult struct {
	WorldState          []byte
	PatchTXReceiptHash  []byte
	NormalTXReceiptHash []byte
}

type testTransition struct {
	patchTransactions  *testTransactionList
	normalTransactions *testTransactionList
	baseValidators     *testValidatorList
	_result            []byte
	_logBloom          *txresult.LogBloom

	sync.Mutex
	step     transitionStep
	_exeChan chan struct{}
}

func (l *testTransition) setExeChan(ch chan struct{}) {
	l._exeChan = ch
}

func (l *testTransition) exeChan() chan<- struct{} {
	return l._exeChan
}

func (tr *testTransition) PatchTransactions() module.TransactionList {
	tr.Lock()
	defer tr.Unlock()

	return tr.patchTransactions
}

func (tr *testTransition) NormalTransactions() module.TransactionList {
	tr.Lock()
	defer tr.Unlock()

	return tr.normalTransactions
}

func (tr *testTransition) EffectiveTransactions() *testTransactionList {
	if tr.patchTransactions.Len() > 0 {
		return tr.patchTransactions
	}
	return tr.normalTransactions
}

func (tr *testTransition) getErrors() (error, error) {
	for _, tx := range tr.patchTransactions.Transactions {
		if tx.Data.ValidateError != nil {
			return tx.Data.ValidateError, nil
		}
	}
	for _, tx := range tr.normalTransactions.Transactions {
		if tx.Data.ValidateError != nil {
			return tx.Data.ValidateError, nil
		}
	}
	for _, tx := range tr.patchTransactions.Transactions {
		if tx.Data.ExecuteError != nil {
			return nil, tx.Data.ExecuteError
		}
	}
	for _, tx := range tr.normalTransactions.Transactions {
		if tx.Data.ExecuteError != nil {
			return nil, tx.Data.ExecuteError
		}
	}
	return nil, nil
}

func (tr *testTransition) Execute(cb module.TransitionCallback) (canceler func() bool, err error) {
	tr.Lock()
	defer tr.Unlock()

	if tr.step >= transitionStepExecuting {
		return nil, errors.Errorf("already executed")
	}
	verr, eerr := tr.getErrors()
	go func() {
		tr.Lock()
		defer tr.Unlock()

		if verr != nil {
			tr.step = transitionStepFailed
			return
		}

		if tr._exeChan != nil {
			tr.Unlock()
			<-tr._exeChan
			cb.OnValidate(tr, verr)
			tr.Lock()
		} else {
			tr.Unlock()
			cb.OnValidate(tr, verr)
			tr.Lock()
		}

		if tr.step == transitionStepFailed {
			// canceled
			return
		}
		if eerr != nil {
			tr.step = transitionStepFailed
			return
		}
		tr.step = transitionStepSucceed
		if tr._exeChan != nil {
			tr.Unlock()
			<-tr._exeChan
			cb.OnExecute(tr, eerr)
			tr.Lock()
		} else {
			tr.Unlock()
			cb.OnExecute(tr, eerr)
			tr.Lock()
		}
	}()
	return func() bool {
		tr.Lock()
		defer tr.Unlock()

		if tr.step == transitionStepExecuting {
			tr.step = transitionStepFailed
			return true
		}
		return false
	}, nil
}

func (tr *testTransition) Result() []byte {
	tr.Lock()
	defer tr.Unlock()

	if tr.step == transitionStepSucceed {
		if tr._result != nil {
			result := &testTransitionResult{}
			result.WorldState = tr.EffectiveTransactions().effect().WorldState
			tr._result = codec.MustMarshalToBytes(result)
		}
		return tr._result
	}
	return nil
}

func (tr *testTransition) NextValidators() module.ValidatorList {
	tr.Lock()
	defer tr.Unlock()

	if tr.step == transitionStepSucceed {
		nv := tr.EffectiveTransactions().effect().NextValidators
		if nv != nil {
			return nv
		}
		return tr.baseValidators
	}
	return nil
}

func (tr *testTransition) LogBloom() module.LogBloom {
	tr.Lock()
	defer tr.Unlock()

	if tr.step == transitionStepSucceed {
		if tr._logBloom == nil {
			tr._logBloom = txresult.NewLogBloom(nil)
			tr._logBloom.Merge(&tr.patchTransactions.effect().LogBloom)
			tr._logBloom.Merge(&tr.normalTransactions.effect().LogBloom)
		}
		return tr._logBloom
	}
	return nil
}

type testServiceManager struct {
	transactions [][]*testTransaction
	bucket       *bucket
	exeChan      chan struct{}
}

func newTestServiceManager(database db.Database) *testServiceManager {
	sm := &testServiceManager{}
	sm.transactions = make([][]*testTransaction, 2)
	sm.bucket = newBucket(database, db.BytesByHash, nil)
	return sm
}

func (sm *testServiceManager) setTransitionExeChan(ch chan struct{}) {
	sm.exeChan = ch
}

func (sm *testServiceManager) ProposeTransition(parent module.Transition, bi module.BlockInfo) (module.Transition, error) {
	tr := &testTransition{}
	tr.baseValidators = parent.NextValidators().(*testValidatorList)
	tr.patchTransactions = newTestTransactionList(sm.transactions[module.TransactionGroupPatch])
	tr.normalTransactions = newTestTransactionList(sm.transactions[module.TransactionGroupNormal])
	if sm.exeChan != nil {
		tr.setExeChan(sm.exeChan)
	}
	return tr, nil
}

func (sm *testServiceManager) CreateInitialTransition(result []byte, nextValidators module.ValidatorList) (module.Transition, error) {
	if nextValidators == nil {
		nextValidators = newTestValidatorList(nil)
	}
	nvl, ok := nextValidators.(*testValidatorList)
	if !ok {
		return nil, errors.Errorf("bad validator list type")
	}
	tr := &testTransition{}
	tr.patchTransactions = newTestTransactionList(nil)
	tr.normalTransactions = newTestTransactionList(nil)
	tr.normalTransactions._effect.WorldState = result
	tr.normalTransactions._effect.NextValidators = nvl
	tr.step = transitionStepSucceed
	return tr, nil
}

func (sm *testServiceManager) CreateTransition(parent module.Transition, txs module.TransactionList, bi module.BlockInfo) (module.Transition, error) {
	if ttxl, ok := txs.(*testTransactionList); ok {
		for _, ttx := range ttxl.Transactions {
			if ttx.Data.TransitionCreateError != nil {
				return nil, errors.Errorf("bad transaction list")
			}
		}
		tr := &testTransition{}
		tr.baseValidators = parent.NextValidators().(*testValidatorList)
		tr.patchTransactions = newTestTransactionList(nil)
		tr.normalTransactions = ttxl
		if sm.exeChan != nil {
			tr.setExeChan(sm.exeChan)
		}
		return tr, nil
	}
	return nil, errors.Errorf("bad type")
}

func (sm *testServiceManager) GetPatches(parent module.Transition) module.TransactionList {
	return newTestTransactionList(sm.transactions[module.TransactionGroupPatch])
}

func (sm *testServiceManager) PatchTransition(transition module.Transition, patches module.TransactionList) module.Transition {
	var ttxl *testTransactionList
	var ttr *testTransition
	var ok bool
	if ttxl, ok = patches.(*testTransactionList); !ok {
		return nil
	}
	if ttr, ok = transition.(*testTransition); !ok {
		return nil
	}
	tr := &testTransition{}
	tr.baseValidators = tr.baseValidators
	tr.patchTransactions = ttxl
	tr.normalTransactions = ttr.normalTransactions
	return tr
}

func (sm *testServiceManager) Finalize(transition module.Transition, opt int) {
	var tr *testTransition
	var ok bool
	if tr, ok = transition.(*testTransition); !ok {
		return
	}
	if opt&module.FinalizeNormalTransaction != 0 {
		sm.bucket.put(tr.normalTransactions)
	}
	if opt&module.FinalizePatchTransaction != 0 {
		sm.bucket.put(tr.patchTransactions)
	}
	if opt&module.FinalizeResult != 0 {
		sm.bucket.put(tr.NextValidators())
	}
}

func (sm *testServiceManager) TransactionFromBytes(b []byte, blockVersion int) (module.Transaction, error) {
	ttx := &testTransaction{}
	_, err := codec.UnmarshalFromBytes(b, ttx)
	if err != nil {
		return nil, err
	}
	if ttx.Data.CreateError != nil {
		return nil, ttx.Data.CreateError
	}
	return ttx, nil
}

func (sm *testServiceManager) TransactionListFromHash(hash []byte) module.TransactionList {
	ttxs := &testTransactionList{}
	err := sm.bucket.get(raw(hash), ttxs)
	if err != nil {
		return nil
	}
	return ttxs
}

func (sm *testServiceManager) TransactionListFromSlice(txs []module.Transaction, version int) module.TransactionList {
	ttxs := make([]*testTransaction, len(txs))
	for i, tx := range txs {
		var ok bool
		ttxs[i], ok = tx.(*testTransaction)
		if !ok {
			return nil
		}
	}
	return newTestTransactionList(ttxs)
}

func (sm *testServiceManager) ReceiptListFromResult(result []byte, g module.TransactionGroup) module.ReceiptList {
	panic("not implemented")
}

func (sm *testServiceManager) SendTransaction(tx interface{}) ([]byte, error) {
	if ttx, ok := tx.(*testTransaction); ok {
		if ttx.Data.CreateError != nil {
			return nil, ttx.Data.CreateError
		}
		sm.transactions[ttx.Group()] = append(sm.transactions[ttx.Group()], ttx)
		return ttx.ID(), nil
	}
	return nil, errors.Errorf("bad type")
}

func (sm *testServiceManager) Call(result []byte, js []byte, bi module.BlockInfo) (module.Status, interface{}, error) {
	panic("not implemented")
}

func (sm *testServiceManager) ValidatorListFromHash(hash []byte) module.ValidatorList {
	tvl := &testValidatorList{}
	err := sm.bucket.get(raw(hash), tvl)
	if err != nil {
		return nil
	}
	return tvl
}

func (sm *testServiceManager) GetBalance(result []byte, addr module.Address) *big.Int {
	panic("not implemented")
}

type testValidator struct {
	Address_ *common.Address
}

func newTestValidator(addr module.Address) *testValidator {
	v := &testValidator{}
	v.Address_ = common.NewAddress(addr.Bytes())
	return v
}

func (v *testValidator) Address() module.Address {
	return v.Address_
}

func (v *testValidator) PublicKey() []byte {
	return nil
}

func (v *testValidator) Bytes() []byte {
	return v.Address_.Bytes()
}

type testValidatorList struct {
	Validators []*testValidator
}

func newTestValidatorList(validators []*testValidator) *testValidatorList {
	vl := &testValidatorList{}
	vl.Validators = make([]*testValidator, len(validators))
	copy(vl.Validators, validators)
	return vl
}

func (vl *testValidatorList) Hash() []byte {
	return crypto.SHA3Sum256(vl.Bytes())
}

func (vl *testValidatorList) Bytes() []byte {
	return codec.MustMarshalToBytes(vl)
}

func (vl *testValidatorList) Flush() error {
	panic("not implemented")
}

func (vl *testValidatorList) IndexOf(addr module.Address) int {
	for i, v := range vl.Validators {
		if v.Address().Equal(addr) {
			return i
		}
	}
	return -1
}

func (vl *testValidatorList) Len() int {
	return len(vl.Validators)
}

func (vl *testValidatorList) Get(i int) (module.Validator, bool) {
	if i >= 0 && i < len(vl.Validators) {
		return vl.Validators[i], true
	}
	return nil, false
}

type testCommitVoteSet struct {
	zero bool
	Pass bool
}

func newCommitVoteSetFromBytes(bs []byte) module.CommitVoteSet {
	vs := &testCommitVoteSet{}
	if bs == nil {
		vs.zero = true
		return vs
	}
	_, err := codec.UnmarshalFromBytes(bs, vs)
	if err != nil {
		return nil
	}
	return vs
}

func newCommitVoteSet(pass bool) module.CommitVoteSet {
	return &testCommitVoteSet{Pass: pass}
}

func (vs *testCommitVoteSet) Verify(block module.Block, validators module.ValidatorList) error {
	if block.Height() == 0 && vs.zero {
		return nil
	}
	if vs.Pass {
		return nil
	}
	return errors.Errorf("verify error")
}

func (vs *testCommitVoteSet) Bytes() []byte {
	bs, _ := codec.MarshalToBytes(vs)
	return bs
}

func (vs *testCommitVoteSet) Hash() []byte {
	return crypto.SHA3Sum256(vs.Bytes())
}

func newRandomTestValidatorList(n int) *testValidatorList {
	wallets := newWallets(n)
	validators := make([]*testValidator, n)
	for i, w := range wallets {
		validators[i] = newTestValidator(w.Address())
	}
	return newTestValidatorList(validators)
}

func newGenesisTX(n int) *testTransaction {
	tx := newTestTransaction()
	wallets := newWallets(n)
	validators := make([]*testValidator, n)
	for i, w := range wallets {
		validators[i] = newTestValidator(w.Address())
	}
	tx.Data.Effect.NextValidators = newTestValidatorList(validators)
	return tx
}

func newTestChain(database db.Database, gtx *testTransaction) *testChain {
	if database == nil {
		database = newMapDB()
	}
	if gtx == nil {
		gtx = newGenesisTX(defaultValidators)
	}
	return &testChain{
		wallet:   common.NewWallet(),
		database: database,
		gtx:      gtx,
		vld:      newCommitVoteSetFromBytes,
	}
}

func newServiceManager(chain module.Chain) *testServiceManager {
	return newTestServiceManager(chain.Database())
}

func newWallets(n int) []module.Wallet {
	wallets := make([]module.Wallet, n)
	for i, _ := range wallets {
		wallets[i] = common.NewWallet()
	}
	return wallets
}

func newMapDB() db.Database {
	return db.Open("", "mapdb", "")
}