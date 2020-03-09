package main

import (
	"flag"
	"fmt"
	"github.com/morfien101/asg-health/asgmanager"
	"github.com/morfien101/asg-health/ec2metadatareader"
	"os"
	"strings"
)

var (
	version         = "0.0.3"
	actionAbandon   = "ABANDON"
	actionHeartBeat = "HEARTBEAT"
	actionContinue  = "CONTINUE"

	helpBlurb = `
	Use this tool to check if the EC2 instance is 'InService' in it's Autoscaling Group.
	It can also set the instance's custom health attribute to 'Unhealthy' which will cause
	the AutoScaling Group to start the termination process.

	Only a single action can be invoked in a single run.
	It will consume credentials from instance policies or ENV vars.
	There is no provision for manually feeding in credentials and never will be.
`

	versionFlag = flag.Bool("v", false, "Show the version")
	helpFlag    = flag.Bool("h", false, "Show the help menu")
	verboseFlag = flag.Bool("verbose", false, "Will log success statements as well as errors")

	instanceIDFlag   = flag.String("i", "", "instance_id for the EC2 instance. If - is passed the instance ID is determined automatically from the metadata if available")
	setUnhealthyFlag = flag.Bool("set-unhealthy", false, "Set the instance to unhealthy in it's AutoScaling Group")
	isInServiceFlag  = flag.Bool("in-service", false, "Checks to see if the instance is 'InService' in it's AutoScaling Group")
)

func main() {
	digestFlags()
	instanceID := ""
	if *instanceIDFlag == "-" {
		localInstanceID, err := ec2metadatareader.InstanceID()
		if err != nil {
			writeToStdErr(fmt.Sprintf("Could not determine instance id. Error: %s", err))
			os.Exit(1)
		}
		instanceID = localInstanceID
	} else {
		instanceID = *instanceIDFlag
	}

	if *isInServiceFlag {
		ok, currentState, err := asgmanager.IsInService(instanceID)
		if err != nil {
			writeToStdErr(err.Error())
			os.Exit(1)
		}
		if !ok {
			writeToStdErr(fmt.Sprintf("Instance %s in not 'InService'. State: %s", instanceID, currentState))
			os.Exit(1)
		}
		verboseLog(fmt.Sprintf("Instance %s is showing as '%s'", instanceID, currentState))
		return

	} else if *setUnhealthyFlag {
		err := asgmanager.SetUnhealthy(instanceID)
		if err != nil {
			writeToStdErr(fmt.Sprintf("Failed the set the health of the instance. Error: %s", err))
			os.Exit(1)
		}
		verboseLog("Successfully set the instance health to 'Unhealthy'")
		return
	} else {
		fmt.Println("No action specified")
		os.Exit(1)
	}
}

func digestFlags() {
	flag.Parse()
	// These 2 functions have the ability to exit the app
	showStopperFlags()
	validateActions()
}

func showStopperFlags() {
	if *helpFlag {
		fmt.Println(helpBlurb)
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
}

func validateActions() {
	errors := []string{}
	if err := validateRequiredVars(); err != nil {
		errors = append(errors, err.Error())
	}

	if *setUnhealthyFlag && *isInServiceFlag {
		errors = append(errors, "-set-unhealthy and -is-healthy can not be used together")
	}

	if len(errors) != 0 {
		fmt.Println(strings.Join(errors, "\n"))
		os.Exit(1)
	}
}

func validateRequiredVars() error {
	errors := []string{}
	if *setUnhealthyFlag || *isInServiceFlag {
		if *instanceIDFlag == "" {
			errors = append(errors, "-i instance_id must be specified")
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("%s", strings.Join(errors, ","))
	}
	return nil
}

func writeToStdErr(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func verboseLog(s string) {
	if *verboseFlag {
		fmt.Println(s)
	}
}
