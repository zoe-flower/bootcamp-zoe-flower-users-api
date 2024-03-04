// Code generated by counterfeiter. DO NOT EDIT.
package eventbusfakes

import (
	"sync"

	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FakeConsumer struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	ListenStub        func()
	listenMutex       sync.RWMutex
	listenArgsForCall []struct {
	}
	OnStub        func(protoreflect.ProtoMessage, eventbus.EventHandler, ...eventbus.HandlerOption) error
	onMutex       sync.RWMutex
	onArgsForCall []struct {
		arg1 protoreflect.ProtoMessage
		arg2 eventbus.EventHandler
		arg3 []eventbus.HandlerOption
	}
	onReturns struct {
		result1 error
	}
	onReturnsOnCall map[int]struct {
		result1 error
	}
	OnAllStub        func(map[protoreflect.ProtoMessage]eventbus.EventHandler) error
	onAllMutex       sync.RWMutex
	onAllArgsForCall []struct {
		arg1 map[protoreflect.ProtoMessage]eventbus.EventHandler
	}
	onAllReturns struct {
		result1 error
	}
	onAllReturnsOnCall map[int]struct {
		result1 error
	}
	UseStub        func(eventbus.ConsumerMiddleware)
	useMutex       sync.RWMutex
	useArgsForCall []struct {
		arg1 eventbus.ConsumerMiddleware
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConsumer) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	stub := fake.CloseStub
	fakeReturns := fake.closeReturns
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConsumer) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeConsumer) CloseCalls(stub func() error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeConsumer) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsumer) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsumer) Listen() {
	fake.listenMutex.Lock()
	fake.listenArgsForCall = append(fake.listenArgsForCall, struct {
	}{})
	stub := fake.ListenStub
	fake.recordInvocation("Listen", []interface{}{})
	fake.listenMutex.Unlock()
	if stub != nil {
		fake.ListenStub()
	}
}

func (fake *FakeConsumer) ListenCallCount() int {
	fake.listenMutex.RLock()
	defer fake.listenMutex.RUnlock()
	return len(fake.listenArgsForCall)
}

func (fake *FakeConsumer) ListenCalls(stub func()) {
	fake.listenMutex.Lock()
	defer fake.listenMutex.Unlock()
	fake.ListenStub = stub
}

func (fake *FakeConsumer) On(arg1 protoreflect.ProtoMessage, arg2 eventbus.EventHandler, arg3 ...eventbus.HandlerOption) error {
	fake.onMutex.Lock()
	ret, specificReturn := fake.onReturnsOnCall[len(fake.onArgsForCall)]
	fake.onArgsForCall = append(fake.onArgsForCall, struct {
		arg1 protoreflect.ProtoMessage
		arg2 eventbus.EventHandler
		arg3 []eventbus.HandlerOption
	}{arg1, arg2, arg3})
	stub := fake.OnStub
	fakeReturns := fake.onReturns
	fake.recordInvocation("On", []interface{}{arg1, arg2, arg3})
	fake.onMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConsumer) OnCallCount() int {
	fake.onMutex.RLock()
	defer fake.onMutex.RUnlock()
	return len(fake.onArgsForCall)
}

func (fake *FakeConsumer) OnCalls(stub func(protoreflect.ProtoMessage, eventbus.EventHandler, ...eventbus.HandlerOption) error) {
	fake.onMutex.Lock()
	defer fake.onMutex.Unlock()
	fake.OnStub = stub
}

func (fake *FakeConsumer) OnArgsForCall(i int) (protoreflect.ProtoMessage, eventbus.EventHandler, []eventbus.HandlerOption) {
	fake.onMutex.RLock()
	defer fake.onMutex.RUnlock()
	argsForCall := fake.onArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeConsumer) OnReturns(result1 error) {
	fake.onMutex.Lock()
	defer fake.onMutex.Unlock()
	fake.OnStub = nil
	fake.onReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsumer) OnReturnsOnCall(i int, result1 error) {
	fake.onMutex.Lock()
	defer fake.onMutex.Unlock()
	fake.OnStub = nil
	if fake.onReturnsOnCall == nil {
		fake.onReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.onReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsumer) OnAll(arg1 map[protoreflect.ProtoMessage]eventbus.EventHandler) error {
	fake.onAllMutex.Lock()
	ret, specificReturn := fake.onAllReturnsOnCall[len(fake.onAllArgsForCall)]
	fake.onAllArgsForCall = append(fake.onAllArgsForCall, struct {
		arg1 map[protoreflect.ProtoMessage]eventbus.EventHandler
	}{arg1})
	stub := fake.OnAllStub
	fakeReturns := fake.onAllReturns
	fake.recordInvocation("OnAll", []interface{}{arg1})
	fake.onAllMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConsumer) OnAllCallCount() int {
	fake.onAllMutex.RLock()
	defer fake.onAllMutex.RUnlock()
	return len(fake.onAllArgsForCall)
}

func (fake *FakeConsumer) OnAllCalls(stub func(map[protoreflect.ProtoMessage]eventbus.EventHandler) error) {
	fake.onAllMutex.Lock()
	defer fake.onAllMutex.Unlock()
	fake.OnAllStub = stub
}

func (fake *FakeConsumer) OnAllArgsForCall(i int) map[protoreflect.ProtoMessage]eventbus.EventHandler {
	fake.onAllMutex.RLock()
	defer fake.onAllMutex.RUnlock()
	argsForCall := fake.onAllArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsumer) OnAllReturns(result1 error) {
	fake.onAllMutex.Lock()
	defer fake.onAllMutex.Unlock()
	fake.OnAllStub = nil
	fake.onAllReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsumer) OnAllReturnsOnCall(i int, result1 error) {
	fake.onAllMutex.Lock()
	defer fake.onAllMutex.Unlock()
	fake.OnAllStub = nil
	if fake.onAllReturnsOnCall == nil {
		fake.onAllReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.onAllReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsumer) Use(arg1 eventbus.ConsumerMiddleware) {
	fake.useMutex.Lock()
	fake.useArgsForCall = append(fake.useArgsForCall, struct {
		arg1 eventbus.ConsumerMiddleware
	}{arg1})
	stub := fake.UseStub
	fake.recordInvocation("Use", []interface{}{arg1})
	fake.useMutex.Unlock()
	if stub != nil {
		fake.UseStub(arg1)
	}
}

func (fake *FakeConsumer) UseCallCount() int {
	fake.useMutex.RLock()
	defer fake.useMutex.RUnlock()
	return len(fake.useArgsForCall)
}

func (fake *FakeConsumer) UseCalls(stub func(eventbus.ConsumerMiddleware)) {
	fake.useMutex.Lock()
	defer fake.useMutex.Unlock()
	fake.UseStub = stub
}

func (fake *FakeConsumer) UseArgsForCall(i int) eventbus.ConsumerMiddleware {
	fake.useMutex.RLock()
	defer fake.useMutex.RUnlock()
	argsForCall := fake.useArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsumer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.listenMutex.RLock()
	defer fake.listenMutex.RUnlock()
	fake.onMutex.RLock()
	defer fake.onMutex.RUnlock()
	fake.onAllMutex.RLock()
	defer fake.onAllMutex.RUnlock()
	fake.useMutex.RLock()
	defer fake.useMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeConsumer) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ eventbus.Consumer = new(FakeConsumer)