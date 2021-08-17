package images

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
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
	initStdout, initErr := initCmd.Output()
	if initErr != nil {
		return fmt.Errorf("packer failed: %s %s", initErr, initStdout)
	}

	buildPrgSub := "build"
	buildTemplate := "images/packer-files/aws-ubuntu-xenial.pkr.hcl"
	buildCmd := exec.Command(prg, buildPrgSub, buildTemplate)
	env := os.Environ()
	env = append(env, fmt.Sprintf("PKR_VAR_ami_name=%s", image))
	buildCmd.Env = env
	buildStdout, buildErr := buildCmd.Output()
	if buildErr != nil {
		return fmt.Errorf("packer failed: %s", buildErr)
	}

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
