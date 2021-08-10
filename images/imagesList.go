package images

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
)

// ImagesList represents a list of compute images
type ImagesList struct {
	Images []string
}

// Search searches for images in the list
func (il *ImagesList) Search(image string) (bool, int) {
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
func (il *ImagesList) Add(image string) error {
	if found, _ := il.Search(image); found {
		return fmt.Errorf("image %s already in the list", image)
	}

	// TODO: Create an image at a cloud provider unless exists there
	prg := "packer"
	arg1 := "build"
	arg2 := "-var-file=../packer-files/variables.json"
	arg3 := "../packer-files/template.json"
	cmd := exec.Command(prg, arg1, arg2, arg3)
	stdout, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("packer failed: %s", err.Error())
	}

	//il.Images = append(il.Images, image)
	fmt.Println(stdout)
	return nil
}

// Remove removes an image from the list
func (il *ImagesList) Remove(image string) error {
	if found, i := il.Search(image); found {
		// TODO: Remove image from cloud provider if exists there
		il.Images = append(il.Images[:i], il.Images[i+1:]...)
		return nil
	}

	return fmt.Errorf("image %s is not in the list", image)
}

// Load obtains images from an images file
func (il *ImagesList) Load(imagesFile string) error {
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
func (il *ImagesList) Save(imagesFile string) error {
	output := ""

	for _, i := range il.Images {
		// TODO: Test for image at cloud provider before concatenating it
		output += fmt.Sprintln(i)
	}

	// TODO: Put file to cloud provider
	return ioutil.WriteFile(imagesFile, []byte(output), 0644)
}
