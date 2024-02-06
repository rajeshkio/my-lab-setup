package ec2Instance

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateSecurityGroup(ctx *pulumi.Context, prefix, internetCIDR string, vpcID pulumi.IDOutput) (*ec2.SecurityGroup, error) {
	sgs, err := ec2.NewSecurityGroup(ctx, prefix+"-sg", &ec2.SecurityGroupArgs{
		VpcId: vpcID,
		Egress: ec2.SecurityGroupEgressArray{
			ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("-1"),
				FromPort:   pulumi.Int(0),
				ToPort:     pulumi.Int(0),
				CidrBlocks: pulumi.StringArray{pulumi.String(internetCIDR)},
			},
		},
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(80),
				ToPort:     pulumi.Int(80),
				CidrBlocks: pulumi.StringArray{pulumi.String(internetCIDR)},
			},
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(22),
				ToPort:     pulumi.Int(22),
				CidrBlocks: pulumi.StringArray{pulumi.String(internetCIDR)},
			},
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(443),
				ToPort:     pulumi.Int(443),
				CidrBlocks: pulumi.StringArray{pulumi.String(internetCIDR)},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return sgs, nil
}

func LocateAmi(ctx *pulumi.Context) (string, error) {
	ami, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
		MostRecent: pulumi.BoolRef(true),
		Filters: []ec2.GetAmiFilter{
			{
				Name: "name",
				Values: []string{
					"ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*",
				},
			},
			{
				Name: "virtualization-type",
				Values: []string{
					"hvm",
				},
			},
		},
		Owners: []string{
			"099720109477",
		},
	}, nil)
	if err != nil {
		return "", err
	}
	return ami.Id, nil
}

func NewInstance(ctx *pulumi.Context, amiId, instanceTypes, prefix, key string, azs []string, subnet, sgsInstance pulumi.IDOutput) (pulumi.StringOutput, error) {
	userData := `#!/bin/bash
		apt update && apt upgrade -y
        apt install -y nginx
        systemctl start nginx
        systemctl enable nginx
        echo 'Hello from Pulumi!' > /usr/share/nginx/html/index.html`

	instance, err := ec2.NewInstance(ctx, prefix, &ec2.InstanceArgs{
		Ami:                      pulumi.String(amiId),
		InstanceType:             pulumi.String(instanceTypes),
		AvailabilityZone:         pulumi.ToStringArray(azs)[0],
		KeyName:                  pulumi.String(key),
		AssociatePublicIpAddress: pulumi.Bool(true),
		SubnetId:                 subnet,
		VpcSecurityGroupIds: pulumi.StringArray{
			sgsInstance.ToStringOutput(),
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String("ReverseProxy"),
		},
		UserData: pulumi.String(userData),
	})
	if err != nil {
		return pulumi.StringOutput{}, err
	}
	return instance.ID().ToStringOutput(), nil
}

func CreateEIPs(ctx *pulumi.Context, prefix string) (*ec2.Eip, error) {
	eip, err := ec2.NewEip(ctx, prefix+"-eip", &ec2.EipArgs{
		Domain: pulumi.String("vpc"),
	})
	if err != nil {
		return nil, err
	}
	return eip, nil
}

func EipAssociation(ctx *pulumi.Context, prefix string, instanceId, allocationId pulumi.StringOutput) error {
	_, err := ec2.NewEipAssociation(ctx, prefix+"-eipAssoc", &ec2.EipAssociationArgs{
		InstanceId:   instanceId,
		AllocationId: allocationId,
	})
	if err != nil {
		return err
	}
	return nil
}
