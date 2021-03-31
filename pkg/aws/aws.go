package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/submariner-io/cloud-prepare/pkg/api"

	"k8s.io/client-go/rest"
)

type awsCloud struct {
	client         ec2iface.EC2API
	k8sConfig      *rest.Config
	gwInstanceType string
	infraID        string
	region         string
}

// NewCloud creates a new api.Cloud instance which can prepare AWS for Submariner to be deployed on it
func NewCloud(k8sConfig *rest.Config, client ec2iface.EC2API, infraID, region, gwInstanceType string) api.Cloud {
	return &awsCloud{
		client:         client,
		k8sConfig:      k8sConfig,
		gwInstanceType: gwInstanceType,
		infraID:        infraID,
		region:         region,
	}
}

func (ac *awsCloud) PrepareForSubmariner(input api.PrepareForSubmarinerInput, reporter api.Reporter) error {
	reporter.Started("Retrieving VPC ID")
	vpcID, err := ac.getVpcID()
	if err != nil {
		reporter.Failed(err)
		return err
	}
	reporter.Succeeded("Retrieved VPC ID %s", vpcID)

	for _, port := range input.InternalPorts {
		reporter.Started("Opening port %v protocol %s for intra-cluster communications", port.Port, port.Protocol)
		err = ac.allowPortInCluster(vpcID, port.Port, port.Protocol)
		if err != nil {
			reporter.Failed(err)
			return err
		}
		reporter.Succeeded("Opened port %v protocol %s for intra-cluster communications", port.Port, port.Protocol)
	}

	reporter.Started("Creating Submariner gateway security group")
	gatewaySG, err := ac.createGatewaySG(vpcID, input.PublicPorts)
	if err != nil {
		reporter.Failed(err)
		return err
	}
	reporter.Succeeded("Created Submariner gateway security group %s", gatewaySG)

	publicSubnet, err := ac.getPublicSubnet(vpcID)
	if err != nil {
		return err
	}

	publicSubnetName := extractName(publicSubnet.Tags)
	reporter.Started("Adjusting public subnet %s to support Submariner", publicSubnetName)
	err = ac.tagPublicSubnet(publicSubnet.SubnetId)
	if err != nil {
		reporter.Failed(err)
		return err
	}
	reporter.Succeeded("Adjusted public subnet %s to support Submariner", publicSubnetName)

	reporter.Started("Deploying gateway node")
	err = ac.deployGateway(vpcID, gatewaySG, publicSubnet)
	if err != nil {
		reporter.Failed(err)
		return err
	}
	reporter.Succeeded("Deployed gateway node")
	return nil
}

func (ac *awsCloud) CleanupAfterSubmariner(reporter api.Reporter) error {
	return nil
}
