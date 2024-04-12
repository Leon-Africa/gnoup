# DRIFTING THROUH THE COSMOS

Defacto setup for gnoland on local and cloud - next generation smart contract platform for social change.

# System Setup

Please ensure the following prequesites is installed on your system:

 1. Terraform: https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli
 2. Ansible: https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html
 3. aws-cli: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
 4. ssm manager plugin: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html

# Deployment

For setup on aws:

````aws configure````

Provide aws credentials that has the neccesary permissions.

i.e create an IAM user -> create an IAM group with Admin Policy attached -> Add the user to the group -> get the access and secret keys for the user

Launch the infrastucture:

````cd terraform/aws````

````terraform plan````

````terraform apply````

Apply the configuration to the infrastructure:

````cd ansible````

````ansible-galaxy install -r requirements.yml```

````AWS_PROFILE=default ansible-playbook -i inventory/aws_ec2.yml playbooks/gno-node.yml --flush-cache -vvv````

Connect to your Gnonode:

````aws ssm start-session --target instance-id````

Note: you can find your instance-id in the logs from running the ansible script

When your ssm session starts ensure that you:

````sudo su````

and navigate to:

````cd ../../../../home````

# Explore

Your keys can be found at: ````home/keys/mnemonic````

The logs for services running can be found under ````home/logs````

The directories in logs are arranged as follows:

node [node.log]

tx-indexer [tx-indexr.log]

supernova [supernova.log]

txtools [txtools.log]

So if you want to get log data for the running gnoland node: ````tail -f logs/node/node.log````


# Clean up

To clean up resources simply run ````terraform destroy```` from the ````terraform/aws```` directory on your system.
