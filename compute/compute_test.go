package compute

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/helpers"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/network"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/subosito/gotenv"
)

var (
	vmName           = "az-samples-go-" + helpers.GetRandomLetterSequence(10)
	nicName          = "nic1"
	username         = "az-samples-go-user"
	password         = "NoSoupForYou1!"
	sshPublicKeyPath = os.Getenv("HOME") + "/.ssh/id_rsa.pub"

	virtualNetworkName = "vnet1"
	subnet1Name        = "subnet1"
	subnet2Name        = "subnet2"
	nsgName            = "nsg1"
	ipName             = "ip1"
)

func init() {
	err := parseArgs()
	if err != nil {
		log.Fatalln("failed to parse args")
	}
}

func parseArgs() error {
	gotenv.Load()

	virtualNetworkName = os.Getenv("AZ_VNET_NAME")
	flag.StringVar(&virtualNetworkName, "vnetName", virtualNetworkName, "Specify a name for the vnet.")

	err := helpers.ParseArgs()
	if err != nil {
		return fmt.Errorf("cannot parse args: %v", err)
	}

	if !(len(virtualNetworkName) > 0) {
		virtualNetworkName = "vnet1"
	}

	return nil
}

func ExampleCreateVM() {
	ctx := context.Background()

	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, helpers.ResourceGroupName())
	if err != nil {
		helpers.PrintAndLog(err.Error())
	}
	helpers.PrintAndLog("resource group created")

	_, err = network.CreateVirtualNetworkAndSubnets(ctx, virtualNetworkName, subnet1Name, subnet2Name)
	if err != nil {
		helpers.PrintAndLog(err.Error())
	}
	helpers.PrintAndLog("created vnet and 2 subnets")

	_, err = network.CreateNetworkSecurityGroup(ctx, nsgName)
	if err != nil {
		helpers.PrintAndLog(err.Error())
	}
	helpers.PrintAndLog("created network security group")

	_, err = network.CreatePublicIP(ctx, ipName)
	if err != nil {
		helpers.PrintAndLog(err.Error())
	}
	helpers.PrintAndLog("created public IP")

	_, err = network.CreateNIC(ctx, virtualNetworkName, subnet1Name, nsgName, ipName, nicName)
	if err != nil {
		helpers.PrintAndLog(err.Error())
	}
	helpers.PrintAndLog("created nic")

	_, err = CreateVM(ctx, vmName, nicName, username, password, sshPublicKeyPath)
	if err != nil {
		helpers.PrintAndLog(err.Error())
	}
	helpers.PrintAndLog("created VM")

	// Output:
	// resource group created
	// created vnet and 2 subnets
	// created network security group
	// created public IP
	// created nic
	// created VM
}