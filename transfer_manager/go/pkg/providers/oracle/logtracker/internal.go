package logtracker

import (
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/transfer_manager/go/pkg/abstract/coordinator"
	"github.com/doublecloud/transfer/transfer_manager/go/pkg/providers/oracle/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	OracleStateKey = "oracle_state"
)

type InternalLogTracker struct {
	controlPlaneClient coordinator.Coordinator
	transferID         string
}

func NewInternalLogTracker(controlPlaneClient coordinator.Coordinator, transferID string) (*InternalLogTracker, error) {
	tracker := &InternalLogTracker{
		controlPlaneClient: controlPlaneClient,
		transferID:         transferID,
	}
	return tracker, nil
}

func (tracker *InternalLogTracker) TransferID() string {
	return tracker.transferID
}

func (tracker *InternalLogTracker) Init() error {
	return nil
}

func (tracker *InternalLogTracker) ClearPosition() error {
	state := map[string]*coordinator.TransferStateData{
		OracleStateKey: nil,
	}

	if err := tracker.controlPlaneClient.SetTransferState(tracker.transferID, state); err != nil {
		return xerrors.Errorf("Unable to set transfer '%v' state: %w", tracker.transferID, err)
	}

	return nil
}

func (tracker *InternalLogTracker) ReadPosition() (*common.LogPosition, error) {
	state, err := tracker.controlPlaneClient.GetTransferState(tracker.transferID)
	if err != nil {
		return nil, xerrors.Errorf("Unable to get transfer '%v' state: %w", tracker.transferID, err)
	}

	if oracleState, ok := state[OracleStateKey]; ok && oracleState.GetOraclePosition() != nil {
		oracleStatePosition := oracleState.OraclePosition
		if oracleStatePosition.Scn == 0 {
			return nil, nil
		}

		ts := oracleStatePosition.Timestamp
		if oracleStatePosition.RsID == "" {
			//nolint:descriptiveerrors
			return common.NewLogPosition(
				oracleStatePosition.Scn,
				nil,
				nil,
				common.PositionType(oracleStatePosition.Type),
				ts,
			)
		} else {
			//nolint:descriptiveerrors
			return common.NewLogPosition(
				oracleStatePosition.Scn,
				&oracleStatePosition.RsID,
				&oracleStatePosition.SSN,
				common.PositionType(oracleStatePosition.Type),
				ts,
			)
		}
	}

	return nil, nil
}

func (tracker *InternalLogTracker) WritePosition(position *common.LogPosition) error {
	rsID := ""
	if position.RSID() != nil {
		rsID = *position.RSID()
	}

	ssn := uint64(0)
	if position.SSN() != nil {
		ssn = *position.SSN()
	}

	timestamp := timestamppb.New(position.Timestamp())
	state := map[string]*coordinator.TransferStateData{
		OracleStateKey: {
			OraclePosition: &coordinator.OraclePositionState{
				Scn:       position.SCN(),
				RsID:      rsID,
				SSN:       ssn,
				Type:      string(position.Type()),
				Timestamp: timestamp.AsTime(),
			},
		},
	}

	if err := tracker.controlPlaneClient.SetTransferState(tracker.transferID, state); err != nil {
		return xerrors.Errorf("Unable to set transfer '%v' state: %w", tracker.transferID, err)
	}

	return nil
}
