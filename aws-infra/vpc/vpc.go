package vpc

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateVPC(ctx *pulumi.Context, prefix, vpcCIDR string) (*ec2.Vpc, error) {
	vpc, err := ec2.NewVpc(ctx, prefix+"-vpc", &ec2.VpcArgs{
		CidrBlock: pulumi.String(vpcCIDR),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("eks-vpc"),
		},
	})
	if err != nil {
		return nil, err
	}
	return vpc, nil
}

type SubnetGen struct {
	pubCIDRs []string
	azs      []string
	vpcId    pulumi.StringInput
	prefix   string
	ctx      *pulumi.Context
}

func NewSubnetGenerator(ctx *pulumi.Context, pubCIDRs, azs []string, vpcId pulumi.StringInput, prefix string) *SubnetGen {
	return &SubnetGen{
		pubCIDRs: pubCIDRs,
		azs:      azs,
		vpcId:    vpcId,
		prefix:   prefix,
		ctx:      ctx,
	}
}

func (sg *SubnetGen) CreatePubSubnet() ([]pulumi.IDOutput, error) {
	var pubSubnetID []pulumi.IDOutput
	pubSubnetsMap := make(map[string]string)

	for i, key := range sg.pubCIDRs {
		value := sg.azs[i%len(sg.azs)]
		pubSubnetsMap[key] = value
	}

	i := 0
	for cidr, az := range pubSubnetsMap {
		subnet, err := ec2.NewSubnet(sg.ctx, fmt.Sprintf(sg.prefix+"-pub-sub-%d", i), &ec2.SubnetArgs{
			VpcId:            sg.vpcId,
			CidrBlock:        pulumi.String(cidr),
			AvailabilityZone: pulumi.String(az),
		})
		if err != nil {
			return nil, err
		}
		pubSubnetID = append(pubSubnetID, subnet.ID())
		i++
	}
	return pubSubnetID, nil
}

func CreateInternetGateway(ctx *pulumi.Context, prefix string, vpcID pulumi.IDOutput) (*ec2.InternetGateway, error) {
	igw, err := ec2.NewInternetGateway(ctx, prefix+"-igw", &ec2.InternetGatewayArgs{
		VpcId: vpcID,
	})
	if err != nil {
		return nil, err
	}
	return igw, nil
}

func CreatePubRouteTable(ctx *pulumi.Context, prefix, internetCIDR string, vpcID, igwID pulumi.StringInput) (*ec2.RouteTable, error) {
	publicRouteTable, err := ec2.NewRouteTable(ctx, prefix+"-pub-rt", &ec2.RouteTableArgs{
		VpcId: vpcID,
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String(internetCIDR),
				GatewayId: igwID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return publicRouteTable, nil
}

func CreatePubRouteTableAssoc(ctx *pulumi.Context, prefix string, pubSubnetIDs []pulumi.IDOutput, rtID pulumi.IDOutput) error {

	for i, id := range pubSubnetIDs {
		_, err := ec2.NewRouteTableAssociation(ctx, fmt.Sprintf(prefix+"-pub-rta-%d", i), &ec2.RouteTableAssociationArgs{
			SubnetId:     id,
			RouteTableId: rtID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
