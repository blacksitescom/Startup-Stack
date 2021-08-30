package images

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gordonianj/blacksite/execHelpers"
)

type Images struct {
	Images []string
}

// Add creates an AMI at EC2 by calling Packer
// TODO: check for image with same name before adding
// TODO: create name blacksite-%random% rather than from arg
func (il *Images) Add(image string) error {

	prg := "packer"

	initPrgSub := "init"
	initPrgDir := "images/packer-files/"
	initCmd := exec.Command(prg, initPrgSub, initPrgDir)

	var initStdout, initStderr []byte
	var errInitStdout, errInitStderr error
	initStdoutIn, _ := initCmd.StdoutPipe()
	initStderrIn, _ := initCmd.StderrPipe()
	initErr := initCmd.Start()
	if initErr != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", initErr)
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// initWg ensures that we finish
	var initWg sync.WaitGroup
	initWg.Add(1)
	go func() {
		initStdout, errInitStdout = execHelpers.CopyAndCapture(os.Stdout, initStdoutIn)
		initWg.Done()
	}()

	initStderr, errInitStderr = execHelpers.CopyAndCapture(os.Stderr, initStderrIn)

	initWg.Wait()

	initErr = initCmd.Wait()
	if initErr != nil {
		log.Fatalf("initCmd.Run() failed with %s\n", initErr)
	}
	if errInitStdout != nil || errInitStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	initOutStr, initErrStr := string(initStdout), string(initStderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", initOutStr, initErrStr)

	buildPrgSub := "build"
	buildTemplate := "images/packer-files/aws-ubuntu-xenial.pkr.hcl"
	buildCmd := exec.Command(prg, buildPrgSub, buildTemplate)

	env := os.Environ()
	env = append(env, fmt.Sprintf("PKR_VAR_ami_name=%s", image))
	buildCmd.Env = env

	var buildStdout, buildStderr []byte
	var errBuildStdout, errBuildStderr error
	buildStdoutIn, _ := buildCmd.StdoutPipe()
	buildStderrIn, _ := buildCmd.StderrPipe()
	buildErr := buildCmd.Start()
	if buildErr != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", buildErr)
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// buildWg ensures that we finish
	var buildWg sync.WaitGroup
	buildWg.Add(1)
	go func() {
		buildStdout, errBuildStdout = execHelpers.CopyAndCapture(os.Stdout, buildStdoutIn)
		buildWg.Done()
	}()

	buildStderr, errBuildStderr = execHelpers.CopyAndCapture(os.Stderr, buildStderrIn)

	buildWg.Wait()

	buildErr = buildCmd.Wait()
	if buildErr != nil {
		log.Fatalf("buildCmd.Run() failed with %s\n", buildErr)
	}
	if errBuildStdout != nil || errBuildStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	buildOutStr, buildErrStr := string(buildStdout), string(buildStderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", buildOutStr, buildErrStr)

	il.Images = append(il.Images, image)
	return fmt.Errorf("packer succeeded: %s", buildStdout)
}

// Load gets image descriptions from AWS EC2
func (il *Images) Describe(awsRegion string) error {

	self := "self"
	ownerSelf := []*string{&self}
	filterBlacksite := []*ec2.Filter{
		{
			Name:   aws.String("tag:BuiltBy"),
			Values: []*string{aws.String("Blacksite")},
		},
	}
	ownBlacksiteImages := &ec2.DescribeImagesInput{
		Owners:  ownerSelf,
		Filters: filterBlacksite,
	}
	svc := ec2.New(session.New(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	result, err := svc.DescribeImages(ownBlacksiteImages)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	fmt.Println(result)

	//il.Images = append(il.Images, scanner.Text())

	return nil
}

// Describes the latest AMI
func (il *Images) Latest() error {
	output := ""

	for _, i := range il.Images {
		// TODO: Test for image at cloud provider before concatenating it
		output += fmt.Sprintln(i)
	}

	// TODO: Put file to cloud provider
	return nil
}
