package cluster

import (
	"log"

	"github.com/kubernauts/tk8/pkg/common"
	"github.com/spf13/viper"
)

type AWSKubeadmConfig struct {
	AWSRegion          string
	ClusterName        string
	MasterInstanceType string
	WorkerInstanceType string
	SSHPublicKey       string
	MasterSubnetID     string
	WorkerSubnetIDS    []string
	MinWorkerCount     int
	MaxWorkerCount     int
	HostedZone         string
	HostedZonePrivate  bool
	Tags               []string
	Tags2              []string
	SSHAccessCIDR      []string
	APIAccessCIDR      []string
	Addons             []string
}

func SetClusterName() {
	if len(common.Name) < 1 {
		config := GetConfig()
		common.Name = config.ClusterName
	}
}

// ReadViperConfigFile is define the config paths and read the configuration file.
func ReadViperConfigFile(configName string) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/tk8")
	verr := viper.ReadInConfig() // Find and read the config file.
	if verr != nil {             // Handle errors reading the config file.
		log.Fatalln(verr)
	}
}

func GetConfig() AWSKubeadmConfig {
	ReadViperConfigFile("config")
	return AWSKubeadmConfig{
		AWSRegion:          viper.GetString("aws-kubeadm.aws_region"),
		ClusterName:        viper.GetString("aws-kubeadm.cluster_name"),
		MasterInstanceType: viper.GetString("aws-kubeadm.master_instance_type"),
		WorkerInstanceType: viper.GetString("aws-kubeadm.worker_instance_type"),
		SSHPublicKey:       viper.GetString("aws-kubeadm.ssh_public_key"),
		MasterSubnetID:     viper.GetString("aws-kubeadm.master_subnet_id"),
		WorkerSubnetIDS:    viper.GetStringSlice("aws-kubeadm.woker_subnet_ids"),
		MinWorkerCount:     viper.GetInt("aws-kubeadm.min_worker_count"),
		MaxWorkerCount:     viper.GetInt("aws-kubeadm.max_worker_count"),
		HostedZone:         viper.GetString("aws-kubeadm.hosted_zone"),
		HostedZonePrivate:  viper.GetBool("aws-kubeadm.hosted_private_zone"),
		Tags:               viper.GetStringSlice("aws-kubeadm.tags"),
		Tags2:              viper.GetStringSlice("aws-kubeadm.tags2"),
		SSHAccessCIDR:      viper.GetStringSlice("aws-kubeadm.ssh_access_cidr"),
		APIAccessCIDR:      viper.GetStringSlice("aws-kubeadm.api_access_cidr"),
		Addons:             viper.GetStringSlice("aws-kubeadm.addons"),
	}
}
