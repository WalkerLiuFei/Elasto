// Copyright (c) 2017-2019 The Elastos Foundation
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
// 

package state

import (
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/dpos/p2p/peer"
)

type Arbitrators interface {
	Start()
	CheckDPOSIllegalTx(block *types.Block) error
	ProcessSpecialTxPayload(p types.Payload, height uint32) error

	IsArbitrator(pk []byte) bool
	GetArbitrators() [][]byte
	GetCandidates() [][]byte
	GetNextArbitrators() [][]byte
	GetNextCandidates() [][]byte
	GetNeedConnectArbiters() []peer.PID
	GetDutyIndexByHeight(height uint32) int
	GetDutyIndex() int

	GetCurrentRewardData() RewardData
	GetNextRewardData() RewardData
	GetArbitersRoundReward() map[common.Uint168]common.Fixed64
	GetFinalRoundChange() common.Fixed64
	IsInactiveMode() bool
	IsUnderstaffedMode() bool

	GetCRCArbiters() [][]byte
	GetCRCProducer(publicKey []byte) *Producer
	GetCRCArbitrators() map[string]*Producer
	IsCRCArbitrator(pk []byte) bool
	IsActiveProducer(pk []byte) bool
	IsDisabledProducer(pk []byte) bool

	GetOnDutyArbitrator() []byte
	GetNextOnDutyArbitrator(offset uint32) []byte

	GetOnDutyCrossChainArbitrator() []byte
	GetCrossChainArbiters() [][]byte
	GetCrossChainArbitersCount() int
	GetCrossChainArbitersMajorityCount() int

	GetArbitersCount() int
	GetCRCArbitersCount() int
	GetArbitersMajorityCount() int
	HasArbitersMajorityCount(num int) bool
	HasArbitersMinorityCount(num int) bool

	GetSnapshot(height uint32) []*KeyFrame
	DumpInfo(height uint32)
}

type IArbitratorsRecord interface {
	GetHeightsDesc() ([]uint32, error)
	GetCheckPoint(height uint32) (*CheckPoint, error)
	SaveArbitersState(point *CheckPoint) error
}
