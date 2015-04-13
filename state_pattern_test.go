package main

import (
	"fmt"
	"testing"
)

type Machine interface {
	Function()
	On()
	Off()
}

type MachineStateContext struct {
	state Machine
}

func (instance *MachineStateContext) Function() {
	instance.state.Function()
}

func (instance *MachineStateContext) On() {
	instance.state.On()
}

func (instance *MachineStateContext) Off() {
	instance.state.Off()
}

func NewMachineInstance() Machine {
	context := &MachineStateContext{}
	initialState := &OffState{context}
	context.state = initialState
	return context
}

type OffState struct {
	context *MachineStateContext
}

func (instance *OffState) Function() {
	fmt.Println("Cannot function when turned off")
}

func (instance *OffState) On() {
	fmt.Println("coming online...")
	instance.context.state = &OnState{instance.context}
}

func (instance *OffState) Off() {
	fmt.Println("Cannot turn off as I am already off - late night?")
}

type OnState struct {
	context *MachineStateContext
}

func (instance *OnState) Function() {
	fmt.Println("Function completed ok!")
}

func (instance *OnState) On() {
	fmt.Println("Already on and functioning")
}

func (instance *OnState) Off() {
	fmt.Println("going offline...")
	instance.context.state = &OffState{instance.context}
}

func Test_StatePattern(t *testing.T) {
	machine := NewMachineInstance()
	machine.Function()
	machine.On()
	machine.Function()
	machine.On()
	machine.Off()
	machine.Off()
	machine.Function()
}
