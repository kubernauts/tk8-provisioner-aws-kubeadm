# TK8  AWS KUBEADM Provisioner

Using aws-kubeadm provisioner with TK8

## Introduction:

The TK8’s new aws-kubeadm provisioner uses the kubeadm terraform implementation for creating a Kubernetes cluster on AWS.
This integrates the work done by Jakub Scholz in https://github.com/scholzj/aws-kubernetes

## Prerequisites:

To begin using the provisioner, there are a couple of things required to setup.

* AWS access key and secret key.
* Terraform 0.12
* Hosted zone created in AWS region specified [Public or Private]
* Existing VPC, Internet Gateway subnet ids for the region via AWS cli or the console.
* Route added to divert incoming traffic from internet into your VPC.

## Setting Environment Variables:

The following environment variables need to be set up:

* `AWS_ACCESS_KEY_ID` - AWS access key
* `AWS_SECRET_ACCESS_KEY` - AWS secret key

## Getting Started:

The provisioner requires that a config.yaml be created with specifications for your cluster

Example `config.yaml`:

```bash
aws-kubeadm:
  # AWS region where should the AWS Kubernetes be deployed
  aws_region: "eu-west-1"
  # Name for AWS resources
  cluster_name: "elrond"
  # Instance types for mster and worker nodes
  master_instance_type: "t2.medium"
  worker_instance_type: "t2.medium"
  # SSH key for the machines
  ssh_public_key: "~/.ssh/id_rsa.pub"

  # Subnet IDs where the cluster should run (should belong to the same VPC)
  # - Master can be only in single subnet
  # - Workers can be in multiple subnets
  # - Worker subnets can contain also the master subnet
  master_subnet_id: "subnet-0451691dd1232d"
  worker_subnet_ids:
    - "subnet-0451691dd1232d"
  # Number of worker nodes
  min_worker_count: 1
  max_worker_count: 3
  # DNS zone where the domain is placed
  hosted_zone: "yourdomain.com"
  hosted_zone_private: true
  # Tags
  tags:
    Application: "AWS-Kubernetes"

#  Tags in a different format for Auto Scaling Group
  tags2:
    - key: "Application"
      value: "AWS-Kubernetes"
      propagate_at_launch: false

  # Kubernetes Addons
  # Supported addons:
  # https://raw.githubusercontent.com/scholzj/terraform-aws-kubernetes/master/addons/storage-class.yaml
  # https://raw.githubusercontent.com/scholzj/terraform-aws-kubernetes/master/addons/heapster.yaml
  # https://raw.githubusercontent.com/scholzj/terraform-aws-kubernetes/master/addons/dashboard.yaml
  # https://raw.githubusercontent.com/scholzj/terraform-aws-kubernetes/master/addons/external-dns.yaml
  # https://raw.githubusercontent.com/scholzj/terraform-aws-kubernetes/master/addons/autoscaler.yaml
  # https://raw.githubusercontent.com/scholzj/terraform-aws-kubernetes/master/addons/ingress.yaml
  # https://raw.githubusercontent.com/scholzj/terraform-aws-kubernetes/master/addons/fluentd-es-kibana-logging.yaml

  addons:
    - "https://raw.githubusercontent.com/scholzj/terraform-aws-kubernetes/master/addons/dashboard.yaml"

  # List of CIDRs from which SSH access is allowed
  ssh_access_cidr:
    - "0.0.0.0/0"

  # List of  CIDRs from which API access is allowed
  api_access_cidr:
    - "0.0.0.0/0"
```
## AWS kubeadm Deployment

1. Download the latest binary based on your platform from here - https://github.com/kubernauts/tk8/releases
2. Set environment variables.
3. Use a `config.yaml` from the above example.
4. Run `tk8 cluster install aws-kubeadm`.


## Access Kubeconfig

Once the cluster is created , wait for a few minutes
You can get the kubeconfig file for the cluster by executing the below command

* `To copy the kubectl config file using DNS record`
Change the values for `cluster_name and` `hosted_zone`

```bash
scp centos@<cluster_name>.<hosted_zone>:/home/centos/kubeconfig .kubeconfig'
```

* `To copy the kubectl config file using IP address`
Change the value for `elastic_ip_address_of_master_instance`

```bash
scp centos@<elastic_ip_address_of_master_instance>:/home/centos/kubeconfig_ip .kubeconfig'
```

## Troubleshooting

If you are not able to scp / ssh into the nodes , check if you have added a route that channels 
incoming traffic from the internet to your VPC via the internet gateway.

Also make sure the necessary resources are created as mentioned in the pre-requisites.


### To add a route to the internet gateway

1. Go to VPC and click "Internet Gateways" from the left menu.
2. Click "Create internet gateway" button and provide Name tag (any name - optional) and click create.
3. By default, it is detached. So click the Actions drop-down and select "Attach to VPC" and attach it with default VPC
4. Now go to "Route Table" and select default route table and edit the route by clicking "Edit routes" button under Routes tab
5. Then in the Destination text box provide "0.0.0.0/0" and in target select the newly created Internet gateway (starts with igw-alphanumeric) and save the route.
6. Now you should be able to SSH / SCP the EC2 instance.