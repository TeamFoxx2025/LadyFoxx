package server

import (
	"github.com/TeamFoxx2025/LadyFoxx/chain"
	"github.com/TeamFoxx2025/LadyFoxx/consensus"
	consensusDev "github.com/TeamFoxx2025/LadyFoxx/consensus/dev"
	consensusDummy "github.com/TeamFoxx2025/LadyFoxx/consensus/dummy"
	consensusIBFT "github.com/TeamFoxx2025/LadyFoxx/consensus/ibft"
	consensusPolyBFT "github.com/TeamFoxx2025/LadyFoxx/consensus/polybft"
	"github.com/TeamFoxx2025/LadyFoxx/forkmanager"
	"github.com/TeamFoxx2025/LadyFoxx/secrets"
	"github.com/TeamFoxx2025/LadyFoxx/secrets/awsssm"
	"github.com/TeamFoxx2025/LadyFoxx/secrets/gcpssm"
	"github.com/TeamFoxx2025/LadyFoxx/secrets/hashicorpvault"
	"github.com/TeamFoxx2025/LadyFoxx/secrets/local"
	"github.com/TeamFoxx2025/LadyFoxx/state"
)

type GenesisFactoryHook func(config *chain.Chain, engineName string) func(*state.Transition) error

type ConsensusType string

type ForkManagerFactory func(forks *chain.Forks) error

type ForkManagerInitialParamsFactory func(config *chain.Chain) (*forkmanager.ForkParams, error)

const (
	DevConsensus     ConsensusType = "dev"
	IBFTConsensus    ConsensusType = "ibft"
	PolyBFTConsensus ConsensusType = consensusPolyBFT.ConsensusName
	DummyConsensus   ConsensusType = "dummy"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	DevConsensus:     consensusDev.Factory,
	IBFTConsensus:    consensusIBFT.Factory,
	PolyBFTConsensus: consensusPolyBFT.Factory,
	DummyConsensus:   consensusDummy.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
	secrets.GCPSSM:         gcpssm.SecretsManagerFactory,
}

var genesisCreationFactory = map[ConsensusType]GenesisFactoryHook{
	PolyBFTConsensus: consensusPolyBFT.GenesisPostHookFactory,
}

var forkManagerFactory = map[ConsensusType]ForkManagerFactory{
	PolyBFTConsensus: consensusPolyBFT.ForkManagerFactory,
}

var forkManagerInitialParamsFactory = map[ConsensusType]ForkManagerInitialParamsFactory{
	PolyBFTConsensus: consensusPolyBFT.ForkManagerInitialParamsFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
