package status

import (
	"context"

	"github.com/TeamFoxx2025/LadyFoxx/command"
	"github.com/TeamFoxx2025/LadyFoxx/command/helper"
	"github.com/TeamFoxx2025/LadyFoxx/server/proto"
)

var (
	params = &statusParams{}
)

const (
	peerIDFlag = "peer-id"
)

type statusParams struct {
	peerID string

	peerStatus *proto.Peer
}

func (p *statusParams) getRequiredFlags() []string {
	return []string{
		peerIDFlag,
	}
}

func (p *statusParams) initPeerInfo(grpcAddress string) error {
	systemClient, err := helper.GetSystemClientConnection(grpcAddress)
	if err != nil {
		return err
	}

	peerStatus, err := systemClient.PeersStatus(
		context.Background(),
		&proto.PeersStatusRequest{
			Id: p.peerID,
		},
	)
	if err != nil {
		return err
	}

	p.peerStatus = peerStatus

	return nil
}

func (p *statusParams) getResult() command.CommandResult {
	return &PeersStatusResult{
		ID:        p.peerStatus.Id,
		Protocols: p.peerStatus.Protocols,
		Addresses: p.peerStatus.Addrs,
	}
}
