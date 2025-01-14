// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/actions/shared"
	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFAppRepository struct {
	CreateAppStub        func(context.Context, authorization.Info, repositories.CreateAppMessage) (repositories.AppRecord, error)
	createAppMutex       sync.RWMutex
	createAppArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateAppMessage
	}
	createAppReturns struct {
		result1 repositories.AppRecord
		result2 error
	}
	createAppReturnsOnCall map[int]struct {
		result1 repositories.AppRecord
		result2 error
	}
	CreateOrPatchAppEnvVarsStub        func(context.Context, authorization.Info, repositories.CreateOrPatchAppEnvVarsMessage) (repositories.AppEnvVarsRecord, error)
	createOrPatchAppEnvVarsMutex       sync.RWMutex
	createOrPatchAppEnvVarsArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateOrPatchAppEnvVarsMessage
	}
	createOrPatchAppEnvVarsReturns struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}
	createOrPatchAppEnvVarsReturnsOnCall map[int]struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}
	GetAppStub        func(context.Context, authorization.Info, string) (repositories.AppRecord, error)
	getAppMutex       sync.RWMutex
	getAppArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	getAppReturns struct {
		result1 repositories.AppRecord
		result2 error
	}
	getAppReturnsOnCall map[int]struct {
		result1 repositories.AppRecord
		result2 error
	}
	GetAppByNameAndSpaceStub        func(context.Context, authorization.Info, string, string) (repositories.AppRecord, error)
	getAppByNameAndSpaceMutex       sync.RWMutex
	getAppByNameAndSpaceArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}
	getAppByNameAndSpaceReturns struct {
		result1 repositories.AppRecord
		result2 error
	}
	getAppByNameAndSpaceReturnsOnCall map[int]struct {
		result1 repositories.AppRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFAppRepository) CreateApp(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateAppMessage) (repositories.AppRecord, error) {
	fake.createAppMutex.Lock()
	ret, specificReturn := fake.createAppReturnsOnCall[len(fake.createAppArgsForCall)]
	fake.createAppArgsForCall = append(fake.createAppArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateAppMessage
	}{arg1, arg2, arg3})
	stub := fake.CreateAppStub
	fakeReturns := fake.createAppReturns
	fake.recordInvocation("CreateApp", []interface{}{arg1, arg2, arg3})
	fake.createAppMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) CreateAppCallCount() int {
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	return len(fake.createAppArgsForCall)
}

func (fake *CFAppRepository) CreateAppCalls(stub func(context.Context, authorization.Info, repositories.CreateAppMessage) (repositories.AppRecord, error)) {
	fake.createAppMutex.Lock()
	defer fake.createAppMutex.Unlock()
	fake.CreateAppStub = stub
}

func (fake *CFAppRepository) CreateAppArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateAppMessage) {
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	argsForCall := fake.createAppArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) CreateAppReturns(result1 repositories.AppRecord, result2 error) {
	fake.createAppMutex.Lock()
	defer fake.createAppMutex.Unlock()
	fake.CreateAppStub = nil
	fake.createAppReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) CreateAppReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
	fake.createAppMutex.Lock()
	defer fake.createAppMutex.Unlock()
	fake.CreateAppStub = nil
	if fake.createAppReturnsOnCall == nil {
		fake.createAppReturnsOnCall = make(map[int]struct {
			result1 repositories.AppRecord
			result2 error
		})
	}
	fake.createAppReturnsOnCall[i] = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVars(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateOrPatchAppEnvVarsMessage) (repositories.AppEnvVarsRecord, error) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	ret, specificReturn := fake.createOrPatchAppEnvVarsReturnsOnCall[len(fake.createOrPatchAppEnvVarsArgsForCall)]
	fake.createOrPatchAppEnvVarsArgsForCall = append(fake.createOrPatchAppEnvVarsArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateOrPatchAppEnvVarsMessage
	}{arg1, arg2, arg3})
	stub := fake.CreateOrPatchAppEnvVarsStub
	fakeReturns := fake.createOrPatchAppEnvVarsReturns
	fake.recordInvocation("CreateOrPatchAppEnvVars", []interface{}{arg1, arg2, arg3})
	fake.createOrPatchAppEnvVarsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsCallCount() int {
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	return len(fake.createOrPatchAppEnvVarsArgsForCall)
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsCalls(stub func(context.Context, authorization.Info, repositories.CreateOrPatchAppEnvVarsMessage) (repositories.AppEnvVarsRecord, error)) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	defer fake.createOrPatchAppEnvVarsMutex.Unlock()
	fake.CreateOrPatchAppEnvVarsStub = stub
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateOrPatchAppEnvVarsMessage) {
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	argsForCall := fake.createOrPatchAppEnvVarsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsReturns(result1 repositories.AppEnvVarsRecord, result2 error) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	defer fake.createOrPatchAppEnvVarsMutex.Unlock()
	fake.CreateOrPatchAppEnvVarsStub = nil
	fake.createOrPatchAppEnvVarsReturns = struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsReturnsOnCall(i int, result1 repositories.AppEnvVarsRecord, result2 error) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	defer fake.createOrPatchAppEnvVarsMutex.Unlock()
	fake.CreateOrPatchAppEnvVarsStub = nil
	if fake.createOrPatchAppEnvVarsReturnsOnCall == nil {
		fake.createOrPatchAppEnvVarsReturnsOnCall = make(map[int]struct {
			result1 repositories.AppEnvVarsRecord
			result2 error
		})
	}
	fake.createOrPatchAppEnvVarsReturnsOnCall[i] = struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) GetApp(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.AppRecord, error) {
	fake.getAppMutex.Lock()
	ret, specificReturn := fake.getAppReturnsOnCall[len(fake.getAppArgsForCall)]
	fake.getAppArgsForCall = append(fake.getAppArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetAppStub
	fakeReturns := fake.getAppReturns
	fake.recordInvocation("GetApp", []interface{}{arg1, arg2, arg3})
	fake.getAppMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) GetAppCallCount() int {
	fake.getAppMutex.RLock()
	defer fake.getAppMutex.RUnlock()
	return len(fake.getAppArgsForCall)
}

func (fake *CFAppRepository) GetAppCalls(stub func(context.Context, authorization.Info, string) (repositories.AppRecord, error)) {
	fake.getAppMutex.Lock()
	defer fake.getAppMutex.Unlock()
	fake.GetAppStub = stub
}

func (fake *CFAppRepository) GetAppArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.getAppMutex.RLock()
	defer fake.getAppMutex.RUnlock()
	argsForCall := fake.getAppArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) GetAppReturns(result1 repositories.AppRecord, result2 error) {
	fake.getAppMutex.Lock()
	defer fake.getAppMutex.Unlock()
	fake.GetAppStub = nil
	fake.getAppReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) GetAppReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
	fake.getAppMutex.Lock()
	defer fake.getAppMutex.Unlock()
	fake.GetAppStub = nil
	if fake.getAppReturnsOnCall == nil {
		fake.getAppReturnsOnCall = make(map[int]struct {
			result1 repositories.AppRecord
			result2 error
		})
	}
	fake.getAppReturnsOnCall[i] = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) GetAppByNameAndSpace(arg1 context.Context, arg2 authorization.Info, arg3 string, arg4 string) (repositories.AppRecord, error) {
	fake.getAppByNameAndSpaceMutex.Lock()
	ret, specificReturn := fake.getAppByNameAndSpaceReturnsOnCall[len(fake.getAppByNameAndSpaceArgsForCall)]
	fake.getAppByNameAndSpaceArgsForCall = append(fake.getAppByNameAndSpaceArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.GetAppByNameAndSpaceStub
	fakeReturns := fake.getAppByNameAndSpaceReturns
	fake.recordInvocation("GetAppByNameAndSpace", []interface{}{arg1, arg2, arg3, arg4})
	fake.getAppByNameAndSpaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) GetAppByNameAndSpaceCallCount() int {
	fake.getAppByNameAndSpaceMutex.RLock()
	defer fake.getAppByNameAndSpaceMutex.RUnlock()
	return len(fake.getAppByNameAndSpaceArgsForCall)
}

func (fake *CFAppRepository) GetAppByNameAndSpaceCalls(stub func(context.Context, authorization.Info, string, string) (repositories.AppRecord, error)) {
	fake.getAppByNameAndSpaceMutex.Lock()
	defer fake.getAppByNameAndSpaceMutex.Unlock()
	fake.GetAppByNameAndSpaceStub = stub
}

func (fake *CFAppRepository) GetAppByNameAndSpaceArgsForCall(i int) (context.Context, authorization.Info, string, string) {
	fake.getAppByNameAndSpaceMutex.RLock()
	defer fake.getAppByNameAndSpaceMutex.RUnlock()
	argsForCall := fake.getAppByNameAndSpaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *CFAppRepository) GetAppByNameAndSpaceReturns(result1 repositories.AppRecord, result2 error) {
	fake.getAppByNameAndSpaceMutex.Lock()
	defer fake.getAppByNameAndSpaceMutex.Unlock()
	fake.GetAppByNameAndSpaceStub = nil
	fake.getAppByNameAndSpaceReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) GetAppByNameAndSpaceReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
	fake.getAppByNameAndSpaceMutex.Lock()
	defer fake.getAppByNameAndSpaceMutex.Unlock()
	fake.GetAppByNameAndSpaceStub = nil
	if fake.getAppByNameAndSpaceReturnsOnCall == nil {
		fake.getAppByNameAndSpaceReturnsOnCall = make(map[int]struct {
			result1 repositories.AppRecord
			result2 error
		})
	}
	fake.getAppByNameAndSpaceReturnsOnCall[i] = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	fake.getAppMutex.RLock()
	defer fake.getAppMutex.RUnlock()
	fake.getAppByNameAndSpaceMutex.RLock()
	defer fake.getAppByNameAndSpaceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFAppRepository) recordInvocation(key string, args []interface{}) {
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

var _ shared.CFAppRepository = new(CFAppRepository)
