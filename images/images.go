package images

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"sync"

	"github.com/gordonianj/blacksite/execHelpers"
)

// Images represents a list of compute images
type Images struct {
	Images []string
}

// Search searches for images in the list
func (il *Images) Search(image string) (bool, int) {
	sort.Slice(il.Images, func(i, j int) bool {
		return il.Images[i] < il.Images[j]
	})

	// TODO: Test for image at cloud provider before returning true
	i := sort.SearchStrings(il.Images, image)
	if i < len(il.Images) && il.Images[i] == image {
		return true, i
	}

	return false, -1
}

// Add adds an image to the list
func (il *Images) Add(image string) error {

	// TODO: Test for AWS access key and secret key in env
	if found, _ := il.Search(image); found {
		return fmt.Errorf("image %s already in the list", image)
	}

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

// Remove removes an image from the list
func (il *Images) Remove(image string) error {
	if found, i := il.Search(image); found {
		// TODO: Remove image from cloud provider if exists there
		il.Images = append(il.Images[:i], il.Images[i+1:]...)
		return nil
	}

	return fmt.Errorf("image %s is not in the list", image)
}

// Load obtains images from an images file
func (il *Images) Load(imagesFile string) error {
	// TODO: Load file from cloud provider
	f, err := os.Open(imagesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// TODO: Test for image at cloud provider before appending it
		il.Images = append(il.Images, scanner.Text())
	}

	return nil
}

// Save saves images to an images file
func (il *Images) Save(imagesFile string) error {
	output := ""

	for _, i := range il.Images {
		// TODO: Test for image at cloud provider before concatenating it
		output += fmt.Sprintln(i)
	}

	// TODO: Put file to cloud provider
	return ioutil.WriteFile(imagesFile, []byte(output), 0644)
}
