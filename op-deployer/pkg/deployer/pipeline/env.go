package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/broadcaster"

	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/state"

	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"

	"github.com/ethereum-optimism/optimism/op-service/jsonutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

type Env struct {
	StateWriter  StateWriter
	L1ScriptHost *script.Host
	L1Client     *ethclient.Client
	Broadcaster  broadcaster.Broadcaster
	Deployer     common.Address
	Logger       log.Logger
}

type StateWriter interface {
	WriteState(st *state.State) error
}

type stateWriterFunc func(st *state.State) error

func (f stateWriterFunc) WriteState(st *state.State) error {
	return f(st)
}

func WorkdirStateWriter(workdir string) StateWriter {
	return stateWriterFunc(func(st *state.State) error {
		return WriteState(workdir, st)
	})
}

func NoopStateWriter() StateWriter {
	return stateWriterFunc(func(st *state.State) error {
		return nil
	})
}

func ReadIntent(workdir string) (*state.Intent, error) {
	intentPath := path.Join(workdir, "intent.toml")
	intent, err := jsonutil.LoadTOML[state.Intent](intentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read intent file: %w", err)
	}
	return intent, nil
}

func ReadState(workdir string) (*state.State, error) {
	statePath := path.Join(workdir, "state.json")
	st, err := jsonutil.LoadJSON[state.State](statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}
	return st, nil
}

func WriteState(workdir string, st *state.State) error {
	statePath := path.Join(workdir, "state.json")
	if st.PredeployedMap != nil {
		fmt.Println("Saving state to: ", statePath, "with predeployed map")
		for _, chain := range st.Chains {
			// allocsMap will host all the addresses
			var allocsMap = make(map[string]state.PredeployedEntry)

			// if the chain has allocs, we need to add them to the allocsMap
			if chain.Allocs != nil {
				allocsData, err := json.Marshal(chain.Allocs.Data)
				if err != nil {
					return fmt.Errorf("failed to marshal allocs data: %w", err)
				}
				var forgeAllocs foundry.ForgeAllocs
				err = json.Unmarshal(allocsData, &forgeAllocs)
				if err != nil {
					return fmt.Errorf("failed to unmarshal allocs: %w", err)
				}

				err = json.Unmarshal(allocsData, &allocsMap)
				if err != nil {
					return fmt.Errorf("failed to unmarshal allocs: %w", err)
				}
			}

			// add the predeployed addresses to the allocsMap
			for k, v := range st.PredeployedMap {
				// check the key does not exist in allocsMap
				if _, ok := allocsMap[k]; !ok {
					allocsMap[k] = v
					fmt.Println("Adding predeployed address:", k)
				}
			}

			new_allocs, err := json.Marshal(allocsMap)
			if err != nil {
				return fmt.Errorf("failed to marshal allocs: %w", err)
			}

			var fallocs foundry.ForgeAllocs
			if err := json.Unmarshal(new_allocs, &fallocs); err != nil {
				return fmt.Errorf("failed to unmarshal allocs data: %w", err)
			}
			chain.Allocs = &state.GzipData[foundry.ForgeAllocs]{
				Data: &fallocs,
			}
		}
	}
	return st.WriteToFile(statePath)
}

type ArtifactsBundle struct {
	L1 foundry.StatDirFs
	L2 foundry.StatDirFs
}

type Stage func(ctx context.Context, env *Env, bundle ArtifactsBundle, intent *state.Intent, st *state.State) error
