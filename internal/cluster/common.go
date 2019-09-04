package cluster

import (
	"github.com/kubernauts/tk8/pkg/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type additionalTags map[string]string
type AWSKubeadmConfig struct {
	Config `yaml:"aws-kubeadm"`
}
type Config struct {
	AWSRegion                   string            `yaml:"aws_region"`
	ClusterName                 string            `yaml:"cluster_name"`
	MasterInstanceType          string            `yaml:"master_instance_type"`
	WorkerInstanceType          string            `yaml:"worker_instance_type"`
	SSHPublicKey                string            `yaml:"ssh_public_key"`
	MasterSubnetID              string            `yaml:"master_subnet_id"`
	WorkerSubnetIDS             []string          `yaml:"worker_subnet_ids"`
	MinWorkerCount              int               `yaml:"min_worker_count"`
	MaxWorkerCount              int               `yaml:"max_worker_count"`
	HostedZone                  string            `yaml:"hosted_zone"`
	HostedZonePrivate           bool              `yaml:"hosted_zone_private"`
	Tags                        map[string]string `yaml:"tags"`
	Tags2                       []additionalTags  `yaml:"tags2"`
	SSHAccessCIDR               []string          `yaml:"ssh_access_cidr"`
	APIAccessCIDR               []string          `yaml:"api_access_cidr"`
	Addons                      []string          `yaml:"addons"`
	TagsInStringForm            string
	APIAccessCIDRInStringForm   string
	AddonsInStringForm          string
	WorkerSubnetIDSInStringForm string
	SSHAccessCIDRInStringForm   string
	Tags2InStringForm           string
}

func SetClusterName() {
	if len(common.Name) < 1 {
		config := GetConfig()
		common.Name = config.Config.ClusterName
	}
}

func GetConfig() *AWSKubeadmConfig {

	awskubeadm := &AWSKubeadmConfig{}
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, awskubeadm)
	if err != nil {
		log.Printf("marshal err   #%v ", err)
	}

	awskubeadm.TagsInStringForm = convertTagsToString(awskubeadm.Config.Tags)
	awskubeadm.Tags2InStringForm = convertTags2ToString(awskubeadm.Config.Tags2)
	awskubeadm.AddonsInStringForm = converListToString(awskubeadm.Config.Addons)
	awskubeadm.SSHAccessCIDRInStringForm = converListToString(awskubeadm.Config.SSHAccessCIDR)
	awskubeadm.APIAccessCIDRInStringForm = converListToString(awskubeadm.Config.APIAccessCIDR)
	awskubeadm.WorkerSubnetIDSInStringForm = converListToString(awskubeadm.Config.WorkerSubnetIDS)

	log.Println("Config is ----\n ", awskubeadm.Config.Tags2InStringForm)

	return awskubeadm
}

func convertTagsToString(m map[string]string) string {
	var str strings.Builder
	str.WriteString("{")
	for key, value := range m {
		str.WriteString(key)
		str.WriteString("=")
		str.WriteString("\"" + value + "\"")
		str.WriteString("\n")
		//	str.WriteString("\n")
	}
	str.WriteString("}")

	return str.String()
}

func convertTags2ToString(listMapOfTag2 []additionalTags) string {
	var str strings.Builder
	str.WriteString("[")
	for i, m := range listMapOfTag2 {
		str.WriteString("\n{")
		index := 0
		for key, value := range m {
			str.WriteString(key + "= ")
			str.WriteString("\"" + value + "\"")
			if index != len(m)-1 {
				str.WriteString(",")
			}
			index++
		}
		str.WriteString("}")
		if i != len(listMapOfTag2)-1 {
			str.WriteString(",")
		}
	}
	str.WriteString("]")
	return str.String()
}

func converListToString(input []string) string {
	var str strings.Builder
	str.WriteString("[")
	for i, s := range input {
		str.WriteString("\"" + s + "\"")
		if i != len(input)-1 {
			str.WriteString("\n")
		}
	}
	str.WriteString("]")
	return str.String()
}
