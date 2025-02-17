package suite_cancun

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/simulators/ethereum/engine/clmock"
	"github.com/ethereum/hive/simulators/ethereum/engine/globals"
)

// Timestamp delta between genesis and the withdrawals fork
func (bs *CancunBaseSpec) GetCancunGenesisTimeDelta() uint64 {
	return bs.CancunForkHeight * bs.GetBlockTimeIncrements()
}

// Calculates Shanghai fork timestamp given the amount of blocks that need to be
// produced beforehand.
func (bs *CancunBaseSpec) GetCancunForkTime() uint64 {
	return uint64(globals.GenesisTimestamp) + bs.GetCancunGenesisTimeDelta()
}

// Generates the fork config, including cancun fork timestamp.
func (bs *CancunBaseSpec) GetForkConfig() globals.ForkConfig {
	return globals.ForkConfig{
		ShanghaiTimestamp: big.NewInt(0), // No test starts before Shanghai
		CancunTimestamp:   new(big.Int).SetUint64(bs.GetCancunForkTime()),
	}
}

// Get the per-block timestamp increments configured for this test
func (bs *CancunBaseSpec) GetBlockTimeIncrements() uint64 {
	return 1
}

// Timestamp delta between genesis and the withdrawals fork
func (bs *CancunBaseSpec) GetBlobsGenesisTimeDelta() uint64 {
	return bs.CancunForkHeight * bs.GetBlockTimeIncrements()
}

// Calculates Shanghai fork timestamp given the amount of blocks that need to be
// produced beforehand.
func (bs *CancunBaseSpec) GetBlobsForkTime() uint64 {
	return uint64(globals.GenesisTimestamp) + bs.GetBlobsGenesisTimeDelta()
}

// Append the accounts we are going to withdraw to, which should also include
// bytecode for testing purposes.
func (bs *CancunBaseSpec) GetGenesis() *core.Genesis {
	genesis := bs.Spec.GetGenesis()

	// Remove PoW altogether
	genesis.Difficulty = common.Big0
	genesis.Config.TerminalTotalDifficulty = common.Big0
	genesis.Config.Clique = nil
	genesis.ExtraData = []byte{}

	if bs.CancunForkHeight == 0 {
		genesis.BlobGasUsed = new(uint64)
		genesis.ExcessBlobGas = new(uint64)
		genesis.BeaconRoot = new(common.Hash)
	}

	// Add accounts that use the DATAHASH opcode
	datahashCode := []byte{
		0x5F, // PUSH0
		0x80, // DUP1
		0x49, // DATAHASH
		0x55, // SSTORE
		0x60, // PUSH1(0x01)
		0x01,
		0x80, // DUP1
		0x49, // DATAHASH
		0x55, // SSTORE
		0x60, // PUSH1(0x02)
		0x02,
		0x80, // DUP1
		0x49, // DATAHASH
		0x55, // SSTORE
		0x60, // PUSH1(0x03)
		0x03,
		0x80, // DUP1
		0x49, // DATAHASH
		0x55, // SSTORE
	}

	for i := 0; i < DATAHASH_ADDRESS_COUNT; i++ {
		address := big.NewInt(0).Add(DATAHASH_START_ADDRESS, big.NewInt(int64(i)))
		genesis.Alloc[common.BigToAddress(address)] = core.GenesisAccount{
			Code:    datahashCode,
			Balance: common.Big0,
		}
	}

	// Add 1 wei to the 4788 stateful precompile so its account is not deleted at the end of
	// block execution.
	genesis.Alloc[common.BigToAddress(new(big.Int).SetUint64(uint64(HISTORY_STORAGE_ADDRESS)))] = core.GenesisAccount{
		Balance: big.NewInt(1),
	}

	return genesis
}

// Changes the CL Mocker default time increments of 1 to the value specified
// in the test spec.
func (bs *CancunBaseSpec) ConfigureCLMock(cl *clmock.CLMocker) {
	cl.BlockTimestampIncrement = big.NewInt(int64(bs.GetBlockTimeIncrements()))
}
