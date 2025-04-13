// Code generated by http://github.com/gojuno/minimock (v3.4.5). DO NOT EDIT.

package telegram

//go:generate minimock -i github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram.BotAPIInterface -o telegram_bot_mock.go -n BotAPIInterfaceMock -p telegram

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gojuno/minimock/v3"
)

// BotAPIInterfaceMock implements BotAPIInterface
type BotAPIInterfaceMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcSend          func(c TGBotAPI.Chattable) (m1 TGBotAPI.Message, err error)
	funcSendOrigin    string
	inspectFuncSend   func(c TGBotAPI.Chattable)
	afterSendCounter  uint64
	beforeSendCounter uint64
	SendMock          mBotAPIInterfaceMockSend
}

// NewBotAPIInterfaceMock returns a mock for BotAPIInterface
func NewBotAPIInterfaceMock(t minimock.Tester) *BotAPIInterfaceMock {
	m := &BotAPIInterfaceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SendMock = mBotAPIInterfaceMockSend{mock: m}
	m.SendMock.callArgs = []*BotAPIInterfaceMockSendParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mBotAPIInterfaceMockSend struct {
	optional           bool
	mock               *BotAPIInterfaceMock
	defaultExpectation *BotAPIInterfaceMockSendExpectation
	expectations       []*BotAPIInterfaceMockSendExpectation

	callArgs []*BotAPIInterfaceMockSendParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// BotAPIInterfaceMockSendExpectation specifies expectation struct of the BotAPIInterface.Send
type BotAPIInterfaceMockSendExpectation struct {
	mock               *BotAPIInterfaceMock
	params             *BotAPIInterfaceMockSendParams
	paramPtrs          *BotAPIInterfaceMockSendParamPtrs
	expectationOrigins BotAPIInterfaceMockSendExpectationOrigins
	results            *BotAPIInterfaceMockSendResults
	returnOrigin       string
	Counter            uint64
}

// BotAPIInterfaceMockSendParams contains parameters of the BotAPIInterface.Send
type BotAPIInterfaceMockSendParams struct {
	c TGBotAPI.Chattable
}

// BotAPIInterfaceMockSendParamPtrs contains pointers to parameters of the BotAPIInterface.Send
type BotAPIInterfaceMockSendParamPtrs struct {
	c *TGBotAPI.Chattable
}

// BotAPIInterfaceMockSendResults contains results of the BotAPIInterface.Send
type BotAPIInterfaceMockSendResults struct {
	m1  TGBotAPI.Message
	err error
}

// BotAPIInterfaceMockSendOrigins contains origins of expectations of the BotAPIInterface.Send
type BotAPIInterfaceMockSendExpectationOrigins struct {
	origin  string
	originC string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmSend *mBotAPIInterfaceMockSend) Optional() *mBotAPIInterfaceMockSend {
	mmSend.optional = true
	return mmSend
}

// Expect sets up expected params for BotAPIInterface.Send
func (mmSend *mBotAPIInterfaceMockSend) Expect(c TGBotAPI.Chattable) *mBotAPIInterfaceMockSend {
	if mmSend.mock.funcSend != nil {
		mmSend.mock.t.Fatalf("BotAPIInterfaceMock.Send mock is already set by Set")
	}

	if mmSend.defaultExpectation == nil {
		mmSend.defaultExpectation = &BotAPIInterfaceMockSendExpectation{}
	}

	if mmSend.defaultExpectation.paramPtrs != nil {
		mmSend.mock.t.Fatalf("BotAPIInterfaceMock.Send mock is already set by ExpectParams functions")
	}

	mmSend.defaultExpectation.params = &BotAPIInterfaceMockSendParams{c}
	mmSend.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmSend.expectations {
		if minimock.Equal(e.params, mmSend.defaultExpectation.params) {
			mmSend.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSend.defaultExpectation.params)
		}
	}

	return mmSend
}

// ExpectCParam1 sets up expected param c for BotAPIInterface.Send
func (mmSend *mBotAPIInterfaceMockSend) ExpectCParam1(c TGBotAPI.Chattable) *mBotAPIInterfaceMockSend {
	if mmSend.mock.funcSend != nil {
		mmSend.mock.t.Fatalf("BotAPIInterfaceMock.Send mock is already set by Set")
	}

	if mmSend.defaultExpectation == nil {
		mmSend.defaultExpectation = &BotAPIInterfaceMockSendExpectation{}
	}

	if mmSend.defaultExpectation.params != nil {
		mmSend.mock.t.Fatalf("BotAPIInterfaceMock.Send mock is already set by Expect")
	}

	if mmSend.defaultExpectation.paramPtrs == nil {
		mmSend.defaultExpectation.paramPtrs = &BotAPIInterfaceMockSendParamPtrs{}
	}
	mmSend.defaultExpectation.paramPtrs.c = &c
	mmSend.defaultExpectation.expectationOrigins.originC = minimock.CallerInfo(1)

	return mmSend
}

// Inspect accepts an inspector function that has same arguments as the BotAPIInterface.Send
func (mmSend *mBotAPIInterfaceMockSend) Inspect(f func(c TGBotAPI.Chattable)) *mBotAPIInterfaceMockSend {
	if mmSend.mock.inspectFuncSend != nil {
		mmSend.mock.t.Fatalf("Inspect function is already set for BotAPIInterfaceMock.Send")
	}

	mmSend.mock.inspectFuncSend = f

	return mmSend
}

// Return sets up results that will be returned by BotAPIInterface.Send
func (mmSend *mBotAPIInterfaceMockSend) Return(m1 TGBotAPI.Message, err error) *BotAPIInterfaceMock {
	if mmSend.mock.funcSend != nil {
		mmSend.mock.t.Fatalf("BotAPIInterfaceMock.Send mock is already set by Set")
	}

	if mmSend.defaultExpectation == nil {
		mmSend.defaultExpectation = &BotAPIInterfaceMockSendExpectation{mock: mmSend.mock}
	}
	mmSend.defaultExpectation.results = &BotAPIInterfaceMockSendResults{m1, err}
	mmSend.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmSend.mock
}

// Set uses given function f to mock the BotAPIInterface.Send method
func (mmSend *mBotAPIInterfaceMockSend) Set(f func(c TGBotAPI.Chattable) (m1 TGBotAPI.Message, err error)) *BotAPIInterfaceMock {
	if mmSend.defaultExpectation != nil {
		mmSend.mock.t.Fatalf("Default expectation is already set for the BotAPIInterface.Send method")
	}

	if len(mmSend.expectations) > 0 {
		mmSend.mock.t.Fatalf("Some expectations are already set for the BotAPIInterface.Send method")
	}

	mmSend.mock.funcSend = f
	mmSend.mock.funcSendOrigin = minimock.CallerInfo(1)
	return mmSend.mock
}

// When sets expectation for the BotAPIInterface.Send which will trigger the result defined by the following
// Then helper
func (mmSend *mBotAPIInterfaceMockSend) When(c TGBotAPI.Chattable) *BotAPIInterfaceMockSendExpectation {
	if mmSend.mock.funcSend != nil {
		mmSend.mock.t.Fatalf("BotAPIInterfaceMock.Send mock is already set by Set")
	}

	expectation := &BotAPIInterfaceMockSendExpectation{
		mock:               mmSend.mock,
		params:             &BotAPIInterfaceMockSendParams{c},
		expectationOrigins: BotAPIInterfaceMockSendExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmSend.expectations = append(mmSend.expectations, expectation)
	return expectation
}

// Then sets up BotAPIInterface.Send return parameters for the expectation previously defined by the When method
func (e *BotAPIInterfaceMockSendExpectation) Then(m1 TGBotAPI.Message, err error) *BotAPIInterfaceMock {
	e.results = &BotAPIInterfaceMockSendResults{m1, err}
	return e.mock
}

// Times sets number of times BotAPIInterface.Send should be invoked
func (mmSend *mBotAPIInterfaceMockSend) Times(n uint64) *mBotAPIInterfaceMockSend {
	if n == 0 {
		mmSend.mock.t.Fatalf("Times of BotAPIInterfaceMock.Send mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmSend.expectedInvocations, n)
	mmSend.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmSend
}

func (mmSend *mBotAPIInterfaceMockSend) invocationsDone() bool {
	if len(mmSend.expectations) == 0 && mmSend.defaultExpectation == nil && mmSend.mock.funcSend == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmSend.mock.afterSendCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmSend.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Send implements BotAPIInterface
func (mmSend *BotAPIInterfaceMock) Send(c TGBotAPI.Chattable) (m1 TGBotAPI.Message, err error) {
	mm_atomic.AddUint64(&mmSend.beforeSendCounter, 1)
	defer mm_atomic.AddUint64(&mmSend.afterSendCounter, 1)

	mmSend.t.Helper()

	if mmSend.inspectFuncSend != nil {
		mmSend.inspectFuncSend(c)
	}

	mm_params := BotAPIInterfaceMockSendParams{c}

	// Record call args
	mmSend.SendMock.mutex.Lock()
	mmSend.SendMock.callArgs = append(mmSend.SendMock.callArgs, &mm_params)
	mmSend.SendMock.mutex.Unlock()

	for _, e := range mmSend.SendMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.m1, e.results.err
		}
	}

	if mmSend.SendMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSend.SendMock.defaultExpectation.Counter, 1)
		mm_want := mmSend.SendMock.defaultExpectation.params
		mm_want_ptrs := mmSend.SendMock.defaultExpectation.paramPtrs

		mm_got := BotAPIInterfaceMockSendParams{c}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.c != nil && !minimock.Equal(*mm_want_ptrs.c, mm_got.c) {
				mmSend.t.Errorf("BotAPIInterfaceMock.Send got unexpected parameter c, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmSend.SendMock.defaultExpectation.expectationOrigins.originC, *mm_want_ptrs.c, mm_got.c, minimock.Diff(*mm_want_ptrs.c, mm_got.c))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSend.t.Errorf("BotAPIInterfaceMock.Send got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmSend.SendMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSend.SendMock.defaultExpectation.results
		if mm_results == nil {
			mmSend.t.Fatal("No results are set for the BotAPIInterfaceMock.Send")
		}
		return (*mm_results).m1, (*mm_results).err
	}
	if mmSend.funcSend != nil {
		return mmSend.funcSend(c)
	}
	mmSend.t.Fatalf("Unexpected call to BotAPIInterfaceMock.Send. %v", c)
	return
}

// SendAfterCounter returns a count of finished BotAPIInterfaceMock.Send invocations
func (mmSend *BotAPIInterfaceMock) SendAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSend.afterSendCounter)
}

// SendBeforeCounter returns a count of BotAPIInterfaceMock.Send invocations
func (mmSend *BotAPIInterfaceMock) SendBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSend.beforeSendCounter)
}

// Calls returns a list of arguments used in each call to BotAPIInterfaceMock.Send.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSend *mBotAPIInterfaceMockSend) Calls() []*BotAPIInterfaceMockSendParams {
	mmSend.mutex.RLock()

	argCopy := make([]*BotAPIInterfaceMockSendParams, len(mmSend.callArgs))
	copy(argCopy, mmSend.callArgs)

	mmSend.mutex.RUnlock()

	return argCopy
}

// MinimockSendDone returns true if the count of the Send invocations corresponds
// the number of defined expectations
func (m *BotAPIInterfaceMock) MinimockSendDone() bool {
	if m.SendMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.SendMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.SendMock.invocationsDone()
}

// MinimockSendInspect logs each unmet expectation
func (m *BotAPIInterfaceMock) MinimockSendInspect() {
	for _, e := range m.SendMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to BotAPIInterfaceMock.Send at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterSendCounter := mm_atomic.LoadUint64(&m.afterSendCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.SendMock.defaultExpectation != nil && afterSendCounter < 1 {
		if m.SendMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to BotAPIInterfaceMock.Send at\n%s", m.SendMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to BotAPIInterfaceMock.Send at\n%s with params: %#v", m.SendMock.defaultExpectation.expectationOrigins.origin, *m.SendMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSend != nil && afterSendCounter < 1 {
		m.t.Errorf("Expected call to BotAPIInterfaceMock.Send at\n%s", m.funcSendOrigin)
	}

	if !m.SendMock.invocationsDone() && afterSendCounter > 0 {
		m.t.Errorf("Expected %d calls to BotAPIInterfaceMock.Send at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.SendMock.expectedInvocations), m.SendMock.expectedInvocationsOrigin, afterSendCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *BotAPIInterfaceMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockSendInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *BotAPIInterfaceMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *BotAPIInterfaceMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockSendDone()
}
