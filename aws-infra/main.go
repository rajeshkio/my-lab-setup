package main

import (
	"aws-infra/ec2Instance"
	"aws-infra/vpc"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	prefix := "nginx-proxy"
	vpcCIDR := "172.16.0.0/16"
	pubCIDRs := []string{"172.16.1.0/24"}
	azs := []string{"ap-south-1a"}
	internetCIDR := "0.0.0.0/0"
	instanceTypes := "t2.micro"
	key := "nginx"

	pulumi.Run(func(ctx *pulumi.Context) error {
		vpcInstance, err := vpc.CreateVPC(ctx, prefix, vpcCIDR)
		if err != nil {
			return err
		}

		SubnetGen := vpc.NewSubnetGenerator(ctx, pubCIDRs, azs, vpcInstance.ID(), prefix)

		pubSubnetIDs, err := SubnetGen.CreatePubSubnet()
		if err != nil {
			return err
		}

		sgsInstance, err := ec2Instance.CreateSecurityGroup(ctx, prefix, internetCIDR, vpcInstance.ID())
		if err != nil {
			return err
		}

		eip, err := ec2Instance.CreateEIPs(ctx, prefix)
		if err != nil {
			return err
		}

		igw, err := vpc.CreateInternetGateway(ctx, prefix, vpcInstance.ID())
		if err != nil {
			return err
		}

		publicRouteTable, err := vpc.CreatePubRouteTable(ctx, prefix, internetCIDR, vpcInstance.ID(), igw.ID())
		if err != nil {
			return err
		}

		err = vpc.CreatePubRouteTableAssoc(ctx, prefix, pubSubnetIDs, publicRouteTable.ID())
		if err != nil {
			return err
		}

		amiId, err := ec2Instance.LocateAmi(ctx)
		if err != nil {
			return err
		}

		instanceId, err := ec2Instance.NewInstance(ctx, amiId, instanceTypes, prefix, key, azs, pubSubnetIDs[0], sgsInstance.ID())
		if err != nil {
			return err
		}

		err = ec2Instance.EipAssociation(ctx, prefix, instanceId, eip.AllocationId)
		if err != nil {
			return err
		}
		ctx.Export("url", pulumi.Sprintf("ssh -i %v.pem ubuntu@%v", key, eip.PublicIp))
		return nil
	})

}
