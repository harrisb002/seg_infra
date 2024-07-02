package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	// All the resources must be wrapped in this pulumi context being created
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create AWS resources for the EC2 Instance

		/*
			Security Group Resources
		*/
		sgArgs := &ec2.SecurityGroupArgs{
			// Ingress is the specified allowed incomming traffic
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080), // Range of Ports
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0/0")},
				}, // Allowing ingress traffic to PORT:8080 from any addr
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22), //For SSH
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0/0")},
				},
			},
			// Allow all traffic to Leave the Server
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"), // Allow all traffic out
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0/0")},
				},
			},
		}

		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			return err
		}

		/*
			Key Pair Resource
		*/
		pk := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDNtOhLX6zuDUGVS8U/zPsMvaoGibMqW9YBtdTZeFlo76BQhDeHLgq8q8IPAr8NIy6iX0XH+uKJF6S/CqBEewgneLomLBpRNDmkFEjcIwXcFyhBqU3hcVWvjdllw/pRteb1RVDp8m3zmMJfzIASfsQisb/ma3Y/Oo7+en+aSViGzYP7QL6w/vOZeiP9jXycJR9LxqeMwLtwL6/jh/be5OIH1idaOKAiFouc58OFGfgCS762mSEnEK5r1bcmHQuoQu3w7I4CvJMB95gdqeFhgBKRBRjmcfNQSWNhYveaCWX/UdmhTsWzo+6c+celUOpqc4A0c0CRkaxh/ogEFU59djPn8GUeWfPhxLttfRi2NnwFsBmApioJfLA2xLPcPE2VeJluWaXjqyVGztcraoMekcgznRf2mAVDnc44iOZtYm3YsI9W/FIcQb3IkcQxwWtzpnsbBTgDEueH4/uAF+gsT1ySUzvre0VJluQTwQMPGVNaiU6zWB6rng1B0J2NbUKsKf8"
		// Used to SSH into Jenkins Server
		kp, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{
			PublicKey: pulumi.String(pk),
		})

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),              // MAKE CERTAIN FREE TIER
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},            // Getting the security group created above
			Ami:                 pulumi.String("ami-066f98455b59ca1ee"), // OS instance AMI ID
			KeyName:             kp.KeyName,                             // Create the Instance
		})

		// For own curiosity
		fmt.Println(jenkinsServer.PublicIp)
		fmt.Println(jenkinsServer.PublicDns)

		// Register these Key Value Pairs with the current Context Stack
		// This stack is created via pulumi, named this "prod"
		ctx.Export("publicIp", jenkinsServer.PublicIp)
		ctx.Export("publicDns", jenkinsServer.PublicDns)

		if err != nil {
			return nil
		}

		return nil
	})
}
