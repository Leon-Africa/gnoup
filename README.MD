[NOTE: THIS IS AN EC2 FOCUSSED DEPLOYMENT]

# gnoup

Defacto setup for gno.land on AWS EC2

# System Setup

````git clone https://github.com/Leon-Africa/gnoup.git && cd gnoup````

## Configure Credentials
We will be using an automated role to allow you to securely deploy your gno.land node to your AWS account:

- Go to AWS Console
- In the console top right hand corner change region to eu-west-1 [Ireland]
- Select Service > CloudFormation > Create stack > With new resources (standard).
- Select upload template and upload the file gno_node_role.yml from this repo and follow the prompts to create the stack.
- Once deployed select the newly created stack.
- Go to the Outputs tab.
- Note the value of the AccessKeyId and SecretAccessKey.

## Docker

Please ensure that you have [Docker](https://www.docker.com/products/docker-desktop/) installed to use the preconfigure image which will setup all the dependancies and you can deploy from the container:


````docker build -t gno-node-deployer .````

````docker run -it -w /home/gnoup --name gno-deployer gno-node-deployer```` 

NOTE: You will deploy and access your gno.land node from the Docker container.

# Deployment

Automatically deploy the infrastructure and configuration as code for a fully running gno.land node with metrics and dashboards to your account from the Docker Container.

````./deploy-gno-node.sh````

NOTE:

- when prompted for the AccessKey and SecretKey enter the values recieved from the Cloudformation deployment.
- enter Default region name as eu-west-1 and Default output format as json

That's it you will now have a fully running gno.land node in AWS with monitoring, observability and logging.


### Connect to your gno.land node:

You can connect to your node in 2 ways.

1. From AWS Console [Make sure you change the region to eu-west-1 in the top right hand corner]

Connect -> Session Manager -> Connect

2. From the docker container:

````aws ssm start-session --target instance-id````

Note: you can find your instance-id in the aws console

In both cases when your ssm session starts ensure that you:

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

So if you want to get log data for the running gno.land node: 

````tail -f logs/node/node.log````

Txdata processed from tx-indexer endpoint can be found at: 

````tail -f logs/txtools/txtools.log````  corresponding code is in ````home/txtools/````

# Metrics

Prometheus:  ````your_aws_instance_ip:9090````

Grafana: ````your_aws_instance_ip:3000```` [username: admin, password: admin] [see Dashboards]

Note: Only accesible from the host machine you deployed from. 

# Clean up

To clean up resources simply run ````terraform destroy```` from the ````terraform/aws directory````.

Note: Do not delete your container before terraform destroy as your [state](https://developer.hashicorp.com/terraform/language/state) is in the container. If for some reason you have exited your container to access it start and access it with:

````docker start gno-deployer````

````docker exec -it gno-deployer /bin/bash````

# Improvements
Check the issues/project section for coming updates. 

Please feel free to create an issue/s for any ideas/improvements/fixes/enhancements. 
