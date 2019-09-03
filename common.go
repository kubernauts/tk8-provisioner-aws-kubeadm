package provisioner

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/kubernauts/tk8-provisioner-aws-kubeadm/internal/cluster"
	"github.com/kubernauts/tk8/pkg/provisioner"
	"log"
	"os/exec"
	"strings"
)

type AWSKubeadm struct {
}

func (a AWSKubeadm) Init(args []string) {
}

func (a AWSKubeadm) Setup(args []string) {

	// check for pre-requisites
	// kubectl should be present
	// terraform should be present
	binaryName := "kubectl"
	err := checkPresenceOfBinary(binaryName)
	if err != nil {
		log.Fatalf("kubectl not found ::: %s", err.Error())
	}
	version := "1.10.0"
	err = checkKubectlVersionGreaterThan(version)
	if err != nil {
		log.Fatalf("Error with the kubectl version :: %s", err.Error())
	}

	binaryName = "terraform"
	err = checkPresenceOfBinary(binaryName)
	if err != nil {
		log.Fatalf("terraform not found ::: %s", err.Error())
	}

	// Proceed with installing
	cluster.Install()
}

func (a AWSKubeadm) Scale(args []string) {
	//	cluster.AWSScale()

}

func (a AWSKubeadm) Reset(args []string) {
	//	cluster.AWSReset()

}

func (a AWSKubeadm) Remove(args []string) {
	//	cluster.AWSRemove()

}

func (a AWSKubeadm) Upgrade(args []string) {
	provisioner.NotImplemented()
}

func (a AWSKubeadm) Destroy(args []string) {
	//	cluster.AWSDestroy()
}

func NewAWSKubeadm() provisioner.Provisioner {
	cluster.SetClusterName()
	provisioner := new(AWSKubeadm)
	return provisioner
}

func checkPresenceOfBinary(binaryName string) error {
	_, err := exec.LookPath(binaryName)
	return err
}

func checkKubectlVersionGreaterThan(version string) error {
	cmd := exec.Command("kubectl", "version", "--client", "--short")
	otpt, err := cmd.Output()
	if err != nil {
		return err
	}

	//Check if kubectl version is greater or equal to 1.10
	parts := strings.Split(string(otpt), " ")

	kctlVersion := strings.Replace((parts[2]), "v", "", -1)

	v1, err := semver.Make(version)
	v2, err := semver.Make(strings.TrimSpace(kctlVersion))

	if v2.LT(v1) {
		log.Fatalln("kubectl client version on this system is less than the required version 1.10.0")
		return fmt.Errorf("kubectl client version on this system is less than the required version %s", version)
	}

	return nil
}
