package frontend

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/types"
)

type FrontendMetrics interface {
	opmetrics.RPCServerMetricer
}

type Backend interface {
	sources.SupervisorAdminAPI
	sources.SupervisorQueryAPI
}

type QueryFrontend struct {
	Supervisor sources.SupervisorQueryAPI
	Metrics    FrontendMetrics
}

var _ sources.SupervisorQueryAPI = (*QueryFrontend)(nil)

// CheckMessage checks the safety-level of an individual message.
// The payloadHash references the hash of the message-payload of the message.
func (q *QueryFrontend) CheckMessage(ctx context.Context, identifier types.Identifier, payloadHash common.Hash, executingDescriptor types.ExecutingDescriptor) (types.SafetyLevel, error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_checkMessage")
		defer done()
	}
	return q.Supervisor.CheckMessage(ctx, identifier, payloadHash, executingDescriptor)
}

// CheckMessagesV2 checks the safety-level of a collection of messages,
// and returns if the minimum safety-level is met for all messages.
func (q *QueryFrontend) CheckMessagesV2(
	ctx context.Context,
	messages []types.Message,
	minSafety types.SafetyLevel,
	executingDescriptor types.ExecutingDescriptor) error {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_checkMessagesV2")
		defer done()
	}
	return q.Supervisor.CheckMessagesV2(ctx, messages, minSafety, executingDescriptor)
}

// CheckMessages checks the safety-level of a collection of messages,
// and returns if the minimum safety-level is met for all messages.
// Deprecated: This method does not check for message expiry.
func (q *QueryFrontend) CheckMessages(
	ctx context.Context,
	messages []types.Message,
	minSafety types.SafetyLevel) error {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_checkMessages")
		defer done()
	}
	return q.Supervisor.CheckMessages(ctx, messages, minSafety)
}

func (q *QueryFrontend) LocalUnsafe(ctx context.Context, chainID eth.ChainID) (eth.BlockID, error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_localUnsafe")
		defer done()
	}
	return q.Supervisor.LocalUnsafe(ctx, chainID)
}

func (q *QueryFrontend) CrossSafe(ctx context.Context, chainID eth.ChainID) (types.DerivedIDPair, error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_crossSafe")
		defer done()
	}
	return q.Supervisor.CrossSafe(ctx, chainID)
}

func (q *QueryFrontend) Finalized(ctx context.Context, chainID eth.ChainID) (eth.BlockID, error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_finalized")
		defer done()
	}
	return q.Supervisor.Finalized(ctx, chainID)
}

func (q *QueryFrontend) FinalizedL1(ctx context.Context) (eth.BlockRef, error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_finalizedL1")
		defer done()
	}
	return q.Supervisor.FinalizedL1(ctx)
}

// CrossDerivedFrom is deprecated, but remains for backwards compatibility to callers
// it is equivalent to CrossDerivedToSource
func (q *QueryFrontend) CrossDerivedFrom(ctx context.Context, chainID eth.ChainID, derived eth.BlockID) (derivedFrom eth.BlockRef, err error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_crossDerivedFrom")
		defer done()
	}
	return q.Supervisor.CrossDerivedToSource(ctx, chainID, derived)
}

func (q *QueryFrontend) CrossDerivedToSource(ctx context.Context, chainID eth.ChainID, derived eth.BlockID) (derivedFrom eth.BlockRef, err error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_crossDerivedToSource")
		defer done()
	}
	return q.Supervisor.CrossDerivedToSource(ctx, chainID, derived)
}

func (q *QueryFrontend) SuperRootAtTimestamp(ctx context.Context, timestamp hexutil.Uint64) (eth.SuperRootResponse, error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_superRootAtTimestamp")
		defer done()
	}
	return q.Supervisor.SuperRootAtTimestamp(ctx, timestamp)
}

func (q *QueryFrontend) AllSafeDerivedAt(ctx context.Context, derivedFrom eth.BlockID) (derived map[eth.ChainID]eth.BlockID, err error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_allSafeDerivedAt")
		defer done()
	}
	return q.Supervisor.AllSafeDerivedAt(ctx, derivedFrom)
}

func (q *QueryFrontend) SyncStatus(ctx context.Context) (eth.SupervisorSyncStatus, error) {
	if q.Metrics != nil {
		done := q.Metrics.RecordRPCServerRequest("supervisor_syncStatus")
		defer done()
	}
	return q.Supervisor.SyncStatus(ctx)
}

type AdminFrontend struct {
	Supervisor Backend
	Metrics    FrontendMetrics
}

var _ sources.SupervisorAdminAPI = (*AdminFrontend)(nil)

// Start starts the service, if it was previously stopped.
func (a *AdminFrontend) Start(ctx context.Context) error {
	if a.Metrics != nil {
		done := a.Metrics.RecordRPCServerRequest("admin_start")
		defer done()
	}
	return a.Supervisor.Start(ctx)
}

// Stop stops the service, if it was previously started.
func (a *AdminFrontend) Stop(ctx context.Context) error {
	if a.Metrics != nil {
		done := a.Metrics.RecordRPCServerRequest("admin_stop")
		defer done()
	}
	return a.Supervisor.Stop(ctx)
}

// AddL2RPC adds a new L2 chain to the supervisor backend
func (a *AdminFrontend) AddL2RPC(ctx context.Context, rpc string, jwtSecret eth.Bytes32) error {
	if a.Metrics != nil {
		done := a.Metrics.RecordRPCServerRequest("admin_addL2RPC")
		defer done()
	}
	return a.Supervisor.AddL2RPC(ctx, rpc, jwtSecret)
}
