package cap

import (
	"reflect"
	"errors"
)

type CallbackInfo struct {
	CallbackType reflect.Type
}

type CallbackRegister struct {
	// [groupName][routeName] = CallbackInfo
	Routers				map[string]map[string]*CallbackInfo
}

func NewCallbackInfo(callback CallbackInterface) *CallbackInfo {
	reflectValue := reflect.ValueOf(callback)
	t := reflect.Indirect(reflectValue).Type()

	cb := CallbackInfo {
		CallbackType: t,
	}

	return &cb
}

func NewCallbackRegister() *CallbackRegister {
	return &CallbackRegister {
		Routers: make(map[string]map[string]*CallbackInfo),
	}
}

func (this *CallbackRegister) Add(group, name string, callback CallbackInterface) {
	if v, ok := this.Routers[group]; ok {
		if _, ok := v[name]; ok {
			panic("Duplicate group: " + group + " name: " + name)
		}
	} else {
		this.Routers[group] = make(map[string]*CallbackInfo)
	}

	this.Routers[group][name] = NewCallbackInfo(callback)
}

func (this *CallbackRegister) Get(group, name string) (CallbackInterface, error) {
	if v, ok := this.Routers[group]; ok {
		if cbInfo, ok := v[name]; ok {
			return active(cbInfo), nil
		}
	}

	return nil, errors.New("Not found " + name + " @ " + group)
}

func active(cbInfo *CallbackInfo) CallbackInterface {
	callback := reflect.New(cbInfo.CallbackType)
	execCallback, ok := callback.Interface().(CallbackInterface)
	if !ok {
		panic("callback is not CallbackInterface")
	}

	return execCallback
}