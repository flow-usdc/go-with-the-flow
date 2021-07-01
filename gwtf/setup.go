package gwtf

import (
	"log"

	"github.com/enescakir/emoji"
	"github.com/spf13/afero"

	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/config"
	"github.com/onflow/flow-cli/pkg/flowkit/gateway"
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
)

// DiscordWebhook stores information about a webhook
type DiscordWebhook struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Wait  bool   `json:"wait"`
}

// GoWithTheFlow Entire configuration to work with Go With the Flow
type GoWithTheFlow struct {
	State    *flowkit.State
	Services *services.Services
	Network  string
	Logger   output.Logger
}


//NewGoWithTheFlowInMemoryEmulator this method is used to create an in memory emulator, deploy all contracts for the emulator and create all accounts
func NewGoWithTheFlowInMemoryEmulator() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "emulator", true).InitializeContracts().CreateAccounts("emulator-account")
}

//NewGoWithTheFlowEmulator create a new client
func NewGoWithTheFlowEmulator() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "emulator", false)
}

func NewGoWithTheFlowDevNet() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "testnet", false)
}

func NewGoWithTheFlowMainNet() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "mainnet", false)
}

// NewGoWithTheFlow with custom file panic on error
func NewGoWithTheFlow(filenames []string, network string, inMemory bool) *GoWithTheFlow {
	gwtf, err := NewGoWithTheFlowError(filenames, network, inMemory)
	if err != nil {
		log.Fatalf("%v error %+v", emoji.PileOfPoo, err)
	}
	return gwtf
}


// NewGoWithTheFlowError creates a new local go with the flow client
func NewGoWithTheFlowError(paths []string, network string, inMemory bool) (*GoWithTheFlow, error) {

	loader := &afero.Afero{Fs: afero.NewOsFs()}
	state, err := flowkit.Load(paths, loader)
	if err != nil {
		return nil, err
	}

	logger := output.NewStdoutLogger(output.InfoLog)


	var service *services.Services
	if inMemory {
		//YAY we can run it inline in memory!
		acc, _ := state.EmulatorServiceAccount()
		//TODO: How can i get the log output here? And enable verbose logging?
		gw := gateway.NewEmulatorGateway(acc)
		service = services.NewServices(gw, state, logger)
	} else {
		host := state.Networks().ByName(network).Host
		gw, err := gateway.NewGrpcGateway(host)
		if err != nil {
			log.Fatal(err)
		}
		service = services.NewServices(gw, state, logger)
	}
	return &GoWithTheFlow{
		State:    state,
		Services: service,
		Network:  network,
		Logger: logger,
	}, nil

}
