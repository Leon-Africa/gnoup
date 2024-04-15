#!/bin/bash

echo "Automatically deploying the infrastructe and configuration as code for a fully running gno node with metrics, dashboards and txtools to your account."

cd terraform/aws

# Initialize Terraform
terraform init

# Generate and review Terraform plan
terraform plan

echo "Intiating Infrastructure"
# Deploy the node infrasructure automatically answering "yes" to any prompts
yes yes | terraform apply

echo "wait for ssm agent"
sleep 30

cd ../../ansible

#Ansible dependancies
ansible-galaxy install -r requirements.yml

echo "Intiating Configuration"

# Configure the node
AWS_PROFILE=default ansible-playbook -i inventory/aws_ec2.yml playbooks/gno-node.yml --flush-cache -vvv

echo "Deployment complete gnogetem!"