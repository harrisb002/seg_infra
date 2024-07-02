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
			Create Security Group Resources
		*/
		sgArgs := &ec2.SecurityGroupArgs{
			// Ingress is the specified allowed incomming traffic
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080), // Range of Ports
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				}, // Allowing ingress traffic to PORT:8080 from any addr
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22), //For SSH
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			// Allow all traffic to Leave the Server
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"), // Allow all traffic out
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			return err
		}

		/*
			Create Key Pair Resource
		*/
		// Used to SSH into Jenkins Server
		kp, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDU59A98tZnA82Bgil2vdi/XT+kR0Z3Itz36Y6pz3TAEu7RpcFamr433LlJ+C7E0cjVFzIaUkOBWdPhIsDdbcZTa+8FbeVoe3MwGHBTdsxvqEm4pQ7QzzPTYhfj1TZvt+Az4bFMkiIEM9OQhBi9M0gBPJC1TnCGdIRkwe9hsKpQWyOs9hRBXUwaugsxaOp5D9mQRZ8QC4G7jZpFPdyP/pAl+0CvGl/Qe1z094oKkkdA39gbWWBhpoFQwGLasyvoMkrAFbSMrkk32gNfBNMopuJw+438AHCpWOg0m/1lMJwn9sar2mdBw7NTx4pcDMVtVg0+G2CwlJAkp2gN3J29UTtICwssLiM5kHanwGDyQA4I4Y+1adHDta8DNh33J9jVmHk1B6Q3EcGQK8SMSLdKbA5qlUVtiY/HwbCMoTcZNxXsAlt3oXGHns7xt9Y4CJScwOhVNiYDFhoqyeAUpVEAgGp0kltp+mQGGwEx+b6C3x24bgOhHuOtPu8UrYe0fvUhn5s= harrisb@Benjamins-Mac-mini.local"),
		})
		if err != nil {
			return err
		}

		// Create New Server Instance
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
		ctx.Export("publicHostName", jenkinsServer.PublicDns)

		if err != nil {
			return nil
		}

		return nil
	})
}
