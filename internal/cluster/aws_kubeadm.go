package cluster

import (
	"bufio"
	"github.com/kubernauts/tk8/pkg/common"
	"github.com/kubernauts/tk8/pkg/provisioner"
	"github.com/kubernauts/tk8/pkg/templates"
	"log"
	"os"
	"os/exec"
)

func Install() {
	err := common.CopyDir("./provisioner/aws-kubeadm", "./inventory/"+common.Name+"/provisioner")

	if err != nil {
		log.Printf("The copy operation failed %q\n", err)
	}

	kubeadmAWSPrepareConfigFiles(common.Name)

	// Check if a terraform state file already exists
	if _, err := os.Stat("./inventory/" + common.Name + "/provisioner/terraform.tfstate"); err == nil {
		log.Fatal("There is an existing cluster, please remove terraform.tfstate file or delete the installation before proceeding")
	} else {
		log.Println("starting terraform init")
		provisioner.ExecuteTerraform("init", "./inventory/"+common.Name+"/provisioner/")
	}
	terrSet := exec.Command("terraform", "apply", "-auto-approve")
	terrSet.Dir = "./inventory/" + common.Name + "/provisioner/"
	stdout, err := terrSet.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()
	if err := terrSet.Start(); err != nil {
		log.Fatal(err)
	}
	if err := terrSet.Wait(); err != nil {
		log.Fatal(err)
	}
	log.Println("Voila! Kubernetes cluster is provisioned in Rancher. Please check the further details about the cluster in Rancher GUI")
	os.Exit(0)
}

func kubeadmAWSPrepareConfigFiles(name string) {
	templates.ParseTemplate(templates.VariablesKubeadmAWS, "./inventory/"+common.Name+"/provisioner/variables.tf", GetConfig())
}
