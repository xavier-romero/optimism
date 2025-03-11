package system

import (
	"fmt"
	"slices"

	"github.com/ethereum-optimism/optimism/devnet-sdk/descriptors"
	"github.com/ethereum-optimism/optimism/devnet-sdk/shell/env"
)

type system struct {
	identifier string
	l1         Chain
	l2s        []L2Chain
}

// system implements System
var _ System = (*system)(nil)

func NewSystemFromURL(url string) (System, error) {
	devnetEnv, err := env.LoadDevnetFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to load devnet from URL: %w", err)
	}

	sys, err := systemFromDevnet(devnetEnv.Config, devnetEnv.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create system from devnet: %w", err)
	}
	return sys, nil
}

func (s *system) L1() Chain {
	return s.l1
}

func (s *system) L2s() []L2Chain {
	return s.l2s
}

func (s *system) Identifier() string {
	return s.identifier
}

func systemFromDevnet(dn descriptors.DevnetEnvironment, identifier string) (System, error) {
	l1, err := newChainFromDescriptor(dn.L1)
	if err != nil {
		return nil, fmt.Errorf("failed to add L1 chain: %w", err)
	}

	l2s := make([]L2Chain, len(dn.L2))
	for i, l2 := range dn.L2 {
		l2s[i], err = newL2ChainFromDescriptor(l2)
		if err != nil {
			return nil, fmt.Errorf("failed to add L2 chain: %w", err)
		}
	}

	sys := &system{
		identifier: identifier,
		l1:         l1,
		l2s:        l2s,
	}

	if slices.Contains(dn.Features, "interop") {
		return &interopSystem{system: sys}, nil
	}

	return sys, nil
}

type interopSystem struct {
	*system
}

// interopSystem implements InteropSystem
var _ InteropSystem = (*interopSystem)(nil)

func (i *interopSystem) InteropSet() InteropSet {
	return i.system // TODO: the interop set might not contain all L2s
}
