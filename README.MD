# GNOUP

Defacto setup for gnoland on cloud - next generation smart contract platform for social change.

# System Setup

## Docker

Please ensure that you have [Docker](https://www.docker.com/products/docker-desktop/) installed to use the preconfigure image which will setup all the dependancies and you can deploy from the container:


````docker build --no-cache -t ubuntu_launcher .````

````docker run -it --name launcher ubuntu_launcher /bin/bash```` 

then go to the #deployment section in this readme.

## Machine

If you decide to use your machine instead of Docker - please ensure the following prequesites are installed on your system:

 1. [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
 2. [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
 3. [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
 4. [ssm manager plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

# Deployment

For setup on aws:

````aws configure````

Provide aws credentials that has the neccesary permissions.

i.e AWS console search IAM -> create an IAM user -> create an IAM group with Admin Policy attached -> Add the user to the group -> get the access and secret keys for the user [select user -> security credentials -> create access key]

Note: set region as ````eu-west-1```` and output as ````json````

## Auto

Automatically deploy the infrastructure and configuration as code for a fully running gno node with metrics, dashboards and txtools to your account:

Run ```` cd gnoup/````
    ````./gnogetem.sh````

After running the script go straight to #connect_to_your_gno_node and #explore in this readme.

If you prefer to do things manually go to the #manual section in this readme.

## Manual

### Launch the infrastucture:

````cd terraform/aws````

````terraform init````

````terraform plan````

````terraform apply````

### Apply the configuration to the infrastructure:

````cd ansible````

````ansible-galaxy install -r requirements.yml````

````AWS_PROFILE=default ansible-playbook -i inventory/aws_ec2.yml playbooks/gno-node.yml --flush-cache -vvv````

### Connect to your Gnonode:

````aws ssm start-session --target instance-id````

Note: you can find your instance-id in the logs from running the ansible script [the blue color] or in the aws console.

When your ssm session starts ensure that you:

````sudo su````

and navigate to:

````cd ../../../../home````

# Explore

Your keys can be found at: ````home/keys/mnemonic````

The logs for services running can be found under ````home/logs````

The directories in logs are arranged as follows:

````node [node.log]````

````tx-indexer [tx-indexr.log]````

````supernova [supernova.log]````

````txtools [txtools.log]````

So if you want to get log data for the running gnoland node: 

````tail -f logs/node/node.log````

Txdata processed from tx-indexer endpoint can be found at: 

````tail -f logs/txtools/txtools.log````  corresponding code is in ````home/txtools/````

# Metrics

Prometheus:  ````your_aws_instance_ip:9090````

Grafana: ````your_aws_instance_ip:3000```` [username: admin, password: admin] [see Dashboards]

Note: Only accesible from the host machine you deployed from. 

# Clean up

To clean up resources simply run ````terraform destroy```` from the ````terraform/aws```` directory.
