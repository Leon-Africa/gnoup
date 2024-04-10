terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "eu-west-1"
}

#Create VPC
resource "aws_vpc" "txindexer-vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "txindexer-vpc"
  }
}

#Create Public Subnet
resource "aws_subnet" "txindexer-public" {
  vpc_id            = aws_vpc.txindexer-vpc.id
  availability_zone = "eu-west-1a"
  cidr_block        = "10.0.1.0/24"

  tags = {
    Name = "txindexer-public"
  }
}

#Create Private Subnet
resource "aws_subnet" "txindexer-private" {
  vpc_id            = aws_vpc.txindexer-vpc.id
  availability_zone = "eu-west-1a"
  cidr_block        = "10.0.2.0/24"

  tags = {
    Name = "txindexer-private"
  }
}

#Create internet Gateway
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.txindexer-vpc.id

  tags = {
    Name = "txindexer-igw"
  }
}

#Create Route Table
resource "aws_route_table" "txindexer-public-rt" {
  vpc_id = aws_vpc.txindexer-vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

#Associate Route Table to Subnet
resource "aws_route_table_association" "txindexer" {
  subnet_id      = aws_subnet.txindexer-public.id
  route_table_id = aws_route_table.txindexer-public-rt.id

}


#Create Security Group
resource "aws_security_group" "txindexer-sg" {
  vpc_id = aws_vpc.txindexer-vpc.id

  name        = "txindexer-sg"
  description = "Security Group for Gno Node"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

#AWS SSM Role
resource "aws_iam_instance_profile" "ssm-profile" {
  name = "EC2SSM"
  role = aws_iam_role.ssm-role.name
}

resource "aws_iam_role" "ssm-role" {
  name               = "EC2SSM"
  description        = "EC2 SSM Role"
  assume_role_policy = <<EOF
{
"Version": "2012-10-17",
"Statement": {
"Effect": "Allow",
"Principal": {"Service": "ec2.amazonaws.com"},
"Action": "sts:AssumeRole"
}
}
EOF

  tags = {
    Name = "txindexer"
  }
}

resource "aws_iam_role_policy_attachment" "ssm-policy" {
  role       = aws_iam_role.ssm-role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "s3-policy" {
  role       = aws_iam_role.ssm-role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_role_policy_attachment" "ec2-policy" {
  role       = aws_iam_role.ssm-role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
}


#Create the EC2 Instance
resource "aws_instance" "txindexer" {
  ami                         = "ami-0f007bf1d5c770c6e"
  instance_type               = "t2.medium"
  subnet_id                   = aws_subnet.txindexer-public.id
  availability_zone           = "eu-west-1a"
  iam_instance_profile        = aws_iam_instance_profile.ssm-profile.name
  associate_public_ip_address = true

  root_block_device {
    volume_size = 7770
    volume_type = "gp2"
    delete_on_termination = true
  }

  vpc_security_group_ids = [
    aws_security_group.txindexer-sg.id,
  ]

  tags = {
    Terraform = "true"
    Name      = "txindexer"
  }

}

#Create EBS Volume
resource "aws_ebs_volume" "txindexer" {
  availability_zone = "eu-west-1a"
  size              = "7770"
  type              = "gp2"
  
  tags = {
    Name = "txindexer-volume"
  }

  lifecycle {
    prevent_destroy = false
    # ignore_changes  = [txindexer]
  }
}

#Attach EBS Volume
resource "aws_volume_attachment" "txindexer" {
  device_name  = "/dev/sdh"
  volume_id    = aws_ebs_volume.txindexer.id
  instance_id  = aws_instance.txindexer.id
  force_detach = false
}

resource "aws_s3_bucket" "ssm-bucket" {
  bucket = "txindexer-aws-ssm-connection-playbook"

  tags = {
    Name = "SSM Connection Bucket"
  }
}
