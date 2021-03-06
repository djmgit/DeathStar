# DeathStar

DeathStar is a tool to loadtest web based services and rest APIs in a easy, automated, cloud native and quick way without spending much time on infrastructure setup for load generation. DeathStar expects a configuration file with the details of the attack targets and it takes care of the infrastructure, spawns an ephemeral infrastreucture to carry out the loadtest on the target, generates the result and then tears down the infrastructure. Hence, the developer or the operations engineer does not have to think about the infrastructure for carrying out the attack/test. All that is exepcted from the user is to provide the config file.

## What is the need for such a tool

When working on or with web based services or APIs, performance testing of the service is extremely important. It gives us an idea some of the following unknowns

- How much traffic (qps) can your application serve at peak hour. What is the qps at which your application breaks.
- How does incoming load affects the resources consumed by your application. Is it memory intensive or processor intensive.
- What is the right resource consumption level at which you should consider autoscaling? Do you need autoscaling at all?
- What is the latency with which your application responds and how does it vary with load?
- How does your application react to sudden bursts of traffic

Above are some of the unknowns which a developer or an operations engineer should try to find out about the web service they are working on. Also not to
mention, with the advent of chaos engineering, we might want to artificially inject sudden large load into our web service to find out how our application
behaves in different situations like this.

There are several opensource tools present which allow us to loadtest applications by running them from command line or from UI. However the issue is, to carry
out such attacks or loadtests, we require the necessary infrastructure from where we can run a given loadtest tool. Lets say we have spawned a server in 
cloud to run a loadtest tool on our application, now once we are done with our testing, we would not want our loadtest servers to be lying around idle, that 
would be a waste of resource. So some one has to take the pain of spawning the server, orchestrating the tests and finally clean up the resources.
Also running such loadtest tools on shared servers which are already running some application can be risky as loadtest tools take up some resources and
that can affect the already running application.

The manual steps mentioned above come up not only as a time consuming problem but also pose issues when we are trying to carry out our loadtests in a automated way
from CI pipelines. Someone would have to provision the neccessary resurce before running the pipeline and then clean it up afterwards or write additional scripts to
automate the same.

Coming back to DeathStar, it is not only a tool which can carry out loadtest and provide necessary data but it also takes care of infrastructure provisioning can
cleanup by itself. All you have to do is provide config details so that it can orchestrate a successful test.

## What exactly does DeathStar do and how does it do

``` 
Before I countinue further, I would like to mention that although deathstar is usable, it is still under development and features are 
being added. Currently it only supports AWS lambda as a compute backend. 
```
DeathStar takes a configuration file, provisions a lambda function along with the handler, orchestrates the loadtest on the configured targets, displays the
recorded metrics and finally cleans up the lambda function, thereby allowing you to carry out your loadtest in a fully automated and serverless manner.

In the config file you can mention multiple targets and for each target you can provide the attack configurations like what should be the rate of the attack, that
is, what should be the traffic qps, what is the endpoint that DeathStar will be hitting with the mentioned qps, what should be the HTTP method (right now it only
supports GET, but will support other methods along with body and url parameters in future) and for how long will the test be continued. It rquires some other
details as well like the lambda memory size to be allocated, timeout, aws region etc, more on in subsequent sections.
DeathStar will parse the config, accordingly create a lambda fucntion, invoke the lambda function for carrying out the attack and finally will destroy the lambda
function. 
This would be a nice time to mention that DeathStar uses <a href="https://github.com/tsenart/vegeta"> Vegeta </a> under the hood to create the load. Vegeta
is a command line tool which allows you to generate load and hit a given web service endpoint. It also provides a library interface which is really awesome
and is being used by DeathStar.

DeathStar is written in golang and uses YAML configurations. In the next sections I will be talking more about how to setup DeathStar, configure it and carry out
your loadtest.

DeathStar has been tested on **Linux** environment. However it can also be trigerred from **macOS** in case if someone wants to kickstart a test from his or her 
local mac system. The procedure for running it on macOS is however different from that on Linux. I will be discussing about this as well in later sections for this
documentation.

### QuickStart

In this section we will see how to configure DeathStar and get started with it.

**Prerequisite**

You must have an active AWS account and the AWS CLI credentials set up in the environment from where you want to trigger DeathStar. This basically means the AWS
creds - ```aws_access_key_id``` and ```aws_secret_access_key``` should be setup either in ```.aws/credentials``` or as envronment variables. Right now DeathStar
does not provide the option to pass these credentials via config file. It might do in future.
Also, you will be requireing a AWS role for your lambda function. Nothing special is required, its just the usual IAM role that the lambda function will assume so
the it has access to AWS resources. You can find more about it <a href="https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-awscli.html"> here </a>.
In short you can give the following trust policy

```
{
  "Version": "",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```
and for permission you can use - ```arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole``` so that the lambda function can upload logs to CloudWatch.
In subsequent releases, the role creation will be automated so that DeathStar can create the role for itself and then destroy the role once test is done. However as
of now its manual.

**Setting up DeathStar**

DeathStar can be setup on your system/pipeline in either of the two ways. Either you can download pre-built binaries from github or you can build it yourself.
``` Please note : below instructions are for linux based systems and not for macOS. We have a separate section for running DeathStar on macOS ```

**Using the pre-built sources**

You can setup DeathStar from prebuilt binary using the following steps:

- Download the pre-built package suitabe for your system from the releases page.
- Decompress the package : ```unzip <package>```
- Optionally you can copy it to your path like ```/usr/local/bin```
- Check out the binary ```./deathstar -help```

**Building from source**

You can build DeathStar from source using the following steps. However before doing so make sure you have Golang 1.13.11 or above present in your system.
For building DeathStar uses GNU make.

- Download or clone this repository
- Open terminal and change directory to the project's directory
- Run ```make build```
- Your binary should be built and present in ```./dist``` directory present in project root.
- Check the built binary using ```./dist/deathstar -help```.
- To install DeathStar to your path, run ```make install```.
- This will install deathstar under /usr/bin and you can invoke ```deathstar``` from anywhere.

Now that DeathStar is setup we are ready for performing our first loadtest.

As mentioned already, DeathStar expects a single yaml configuration file in which the user will have to provide details about the targets to be attacked.
The user will basically have to provide a list of targets and some options related to them. Apart from target details user must provide some infrastructure
related config options. Lets see a sample config file for orchestrating an attack.

```
attacks:
  - attackName: attack-1
    attackDesc: Test attack
    vegetaAttackParams:
      httpMethod: GET
      url: "https://url-1.xyz"
      rate: 2
      duration: 2
  - attackName: attack-2
    attackDesc: Test attack2
    vegetaAttackParams:
      httpMethod: GET
      url: "https://url-2.abc"
      rate: 2
      duration: 2
  - attackName: attack-2
    attackDesc: Test attack2
    vegetaAttackParams:
      httpMethod: GET
      url: "https://url-3.def"
      rate: 2
      duration: 2

lambdaConfig:
  lambdaRole: "lambda_role_arn"
  lambdaMemorySize: 128
  lambdaTimeOut: 3
  lambdaRegion: us-east-1
```

In the above config, we are specifying three targets also known as attacks. In each attack definition we provide details like attack name, description,
attack url, rate (qps) and everything else required by vegeta to carry out the attack on the given target endpoint. I will be mentioning in details about the
individual attack options in details in the configuration reference section.
Once we are done with our target descriptions required to orchestrate our attack, next we need to provide some details about the lambda configuration. These are
required by DeathStar to spawn the on demand lambda function in the desired manner. We provide configs like the role arn - this is the role which we create in
our prerequisite section, we provide the lambda memory size, depending on the attack scale we might want to change this. We can also specify the lambda timeout and
the region where where the function will be spawned.

Once we are satisfied with our configuration, we can trigger DeathStar to actually orchestrate the attack.
In order to run DeathStar, use the following command:

```
./deathstar -deploy -conf config.yml
```

The above command will tell DeathStar to deploy the lambda function using the provided config, invoke the deployed lambda function with the attack details. The
lambda function will then invoke its handler which will inturn invoke vegeta to attack the given targets on after the order. Right now DeathStar only allows
sequential attack on the targets, in future it will provide attacking targets in parallel as well. Once the lambda execution completes, DeathStar will parse
back the result and display it to the user in json format and finally destroy the lambda function.

Example output:

```

Note: all duration like fields are in nanoseconds

[
	{
		"AttackResponseMetrics": {
			"latencies": {
				"total": 318595407,
				"mean": 79648851,
				"50th": 39478544,
				"90th": 201993451,
				"95th": 201993451,
				"99th": 201993451,
				"max": 201993451,
				"min": 37644867
			},
			"bytes_in": {
				"total": 51430,
				"mean": 12857.5
			},
			"bytes_out": {
				"total": 0,
				"mean": 0
			},
			"earliest": "2021-01-03T06:48:46.940882478Z",
			"latest": "2021-01-03T06:48:48.440780199Z",
			"end": "2021-01-03T06:48:48.479866919Z",
			"duration": 1499897721,
			"wait": 39086720,
			"requests": 4,
			"rate": 2.6668485083990605,
			"throughput": 2.599116594967577,
			"success": 1,
			"status_codes": {
				"200": 4
			},
			"errors": []
		},
		"attackDetails": {
			"attackName": "attack-1",
			"attackDesc": "Test attack",
			"vegetaAttackParams": {
				"httpMethod": "GET",
				"url": "https://url-1.xyz",
				"rate": 2,
				"duration": 2
			}
		}
	},
	{
		"AttackResponseMetrics": {
			"latencies": {
				"total": 153277127,
				"mean": 38319281,
				"50th": 17974179,
				"90th": 112314614,
				"95th": 112314614,
				"99th": 112314614,
				"max": 112314614,
				"min": 5014155
			},
			"bytes_in": {
				"total": 506808,
				"mean": 126702
			},
			"bytes_out": {
				"total": 0,
				"mean": 0
			},
			"earliest": "2021-01-03T06:48:49.205092672Z",
			"latest": "2021-01-03T06:48:50.705106735Z",
			"end": "2021-01-03T06:48:50.722270758Z",
			"duration": 1500014063,
			"wait": 17164023,
			"requests": 4,
			"rate": 2.666641666012167,
			"throughput": 2.636473619616992,
			"success": 1,
			"status_codes": {
				"200": 4
			},
			"errors": []
		},
		"attackDetails": {
			"attackName": "attack-2",
			"attackDesc": "Test attack2",
			"vegetaAttackParams": {
				"httpMethod": "GET",
				"url": "https://url-2.abc",
				"rate": 2,
				"duration": 2
			}
		}
	},
	{
		"AttackResponseMetrics": {
			"latencies": {
				"total": 852066357,
				"mean": 213016589,
				"50th": 221749682,
				"90th": 292663955,
				"95th": 292663955,
				"99th": 292663955,
				"max": 292663955,
				"min": 115903038
			},
			"bytes_in": {
				"total": 97152,
				"mean": 24288
			},
			"bytes_out": {
				"total": 0,
				"mean": 0
			},
			"earliest": "2021-01-03T06:48:51.445086615Z",
			"latest": "2021-01-03T06:48:52.945171147Z",
			"end": "2021-01-03T06:48:53.061074185Z",
			"duration": 1500084532,
			"wait": 115903038,
			"requests": 4,
			"rate": 2.6665163960240075,
			"throughput": 2.475266564086257,
			"success": 1,
			"status_codes": {
				"200": 4
			},
			"errors": []
		},
		"attackDetails": {
			"attackName": "attack-2",
			"attackDesc": "Test attack2",
			"vegetaAttackParams": {
				"httpMethod": "GET",
				"url": "https://url-3.def",
				"rate": 2,
				"duration": 2
			}
		}
	}
]
```

The output shows various metrics measured by vegeta like mean latency, success ratio, status codes etc. Thus, we carried out a succesful loadtest without spending
much time on the infrastructure side much. However as it was mentioned before, DeathStar does have some limitations and pitfalls, which I will be discussing about
in the future scope and limitations section.
Also you might be wondering where is the lambda function code and handler and why are we providing a CLI option like ```-deploy```. I will be explaining the same
in the subsequent sections.

## Configuration options : Orchestrating a loadtest

Lets use our previous sample config yaml file to go over the different config options

```
attacks:
  - attackName: attack-1
    attackDesc: Test attack
    vegetaAttackParams:
      httpMethod: GET
      url: "https://url-1.xyz"
      rate: 2
      duration: 2

lambdaConfig:
  lambdaRole: "lambda_role_arn"
  lambdaMemorySize: 128
  lambdaTimeOut: 3
  lambdaRegion: us-east-1

```

**attacks**

the first section or yaml object key is ```attacks```. As the name suggests, this takes a list of attack configurations. You can also refer to this as target
definitions. This list basically contains various targets to attack and details regarding how to attack, mostly required by vegeta.

- **attackName**: This is just an identifier for a givenm attack. It is suggested to use some meaningful identifier so that you can easily identify your attack target.

- **attackDesc**: A short description about the attack.

- **vegetaAttackParams**: These are parameters/options which describes your attack as well as are required by vegeta to carry out the attack.

- **url**: the web endpont that needs to be loadtested/attacked. Vegeta will bombard this endpoint with requests. This is the target.

- **httpMethod**: The http method to be used by vegeta while hitting the target. Right now, DeathStar only allows ```GET``` but more types will be allowed soon.

- **rate**: Number of requests to be made per second. This is basically the qps. This expects an integer.

- **duration**: Number of seconds for which the attack will be continued. This should not be more than the lambda function timeout.


**lambdaConfig**

The lambdaConfig yaml object keys takes parameters to configure the lambda function that will be created and invoked by DeathStar.

- **lambdaRole**: This is the ARN of the AWS role to be used by the lambda function. This role must be created before using DeathStar.

- **lambdaMemorySize**: The size of the lambda function. This takes integer and the unit is in Mega Byte. Default value considered by lambda is 128MB. We might
  want to tweak this value depending on the scale of our attack. The vCPUs allocated by AWS to a lambda function is proportional to the memory it is allocated.
  So if you want to generate higher load on your target then you will be requiring higher resources - memory and CPU. So you migth want to increase the memory
  limit of the lambda function.
  
- **lambdaTimeOut**: This specifies the timeout for the lambda fucntion. It expects an integer and the unit is second. This specifies for how much time the function
  will run. After the specified time period is over, the instance of the invoked function will be killed by AWS lambda. The value has a hard limit of 900, which is
  15 minute. Unfortunately this is a limitation. However if you want to attack a target for than 15 minute, you can create two attack definitions with the same
  target or invoke DeathStar again.
  
- **lambdaRegion**: This is simply the AWS region where you want DeathStar to spawn the function.

## How does DeathStar work internally : Prcess flow

Lets start with taking another look at the DeathStar execution command:

```./deathstar -deploy -conf config.yml```

when we invoke this comand, the first thing that DeathStar does is, it tries to open the config file and read the configuration. First it picks up the
config options provided for lambda and then it uses AWS Golang SDK to create a lambda function in the given region. Now as we know in order to create
a lambda function we need a zip file (if not S3) containing the function code. In our case, we do not have a separate
function handler project for DeathStar, but ```DeathStar itself is the function handler``` . If we do not provide the ```-deploy``` CLI option, DeathStar's
main method will invoke :
```
lambda.Start(HandleLambdaEvent)
```
which is basically the AWS Go SDK's way of invoking the lambda function handler to handle request. This handler function will then invoke vegeta library
to actually hit the target. So, DeathStar does both the jobs of orchestrating the attack by creating and invoking the lambda function as well as it
itself handles the function invocation as the the lambda function and carries out the actuall attack. The ```-deploy``` option simply tells DeathStar that
it is not running as a lambda function on AWS, but it is being run from some other system to orchestrate the attack. So to summarise, just before, DeathStar
creates the lambda function on AWS, it creates a zip package of itself (the binary being run) and provides the created zip file's path to the AWS SDK's create lambda function input.
As soon as the function is successfully created, it deletes the zip package.

Now that we know how DeathStar creates and works as the lambda function lets continue with the flow. Next, DeathStar will find out the attack configutations
provided in the config yaml file under the ```attacks``` section and start invoking the lambda function for each of the targets provided. Right now the attacks
are done sequentially one after the other but it will provide a way to do in parallel in future. After it was finished attacking all the targets it will display
the result and then clean up the lambda function.

## Running DeathStar from macOS

All the above steps are for linux based systems. However you can run DeathStar from macOS as well, but the steps will be different to some extend. The main issue
for running DeathStar from macOS is the fact that DeathStar deploys itself as a function handler on AWS lambda as well. So its basically the same binary that runs
on your system (local, jenkins, instance, wherever) as well as on the lambda. However, this is not possible on macOS. Because macOS wont be able to linux ELF
binary (unless you are using a VM or a container in which case the scenario changes completely and will be same as running on linux). So, when running on macOS
we have to provide the zip package which will contain the elf binary of DeathStar and run the other binary on macOS which is built for macOS. And additionaly
we will have to tell DeathStar that do not create a zip of yourself but use the zip package provided by us while creating the lambda function.

The above, might sound complicated, but it is not actually so. Lets see who we can do this easily :

**Using prebuilt packages**

- Setup DeathStar in the same way as shown above using prebuilt packages, however use the prebuild package from ```macOS``` from the realese page.
- Download the zip package for linux from here
- Now you can run DeathStar using the following

```
./deathstar -conf [conf.yml] -deploy -local -zip-file-path [zip-file-path-to-the-linux-zip-package]

```

**By building from source**

- Clone or download this repository and open it in your terminal.
- First we need to make the zip package
- Since we are building a Golang project binary on macOS, we need to tell Golang the target system is linux.
- Run ```GOOS=linux make build```
- This will build deathstar binary in the dist folder, you wont be able to run this binary as this is a linux elf.
- Next for creating the package run ```make lambda_package```.
- The above will create a zip file containing the linux elf in the root named ```deathstar.zip```
- Next we have to build DeathStar for macOS. This can be done using ```make build``` command as we had done earlier. This will replace the linux binary
  in dist folder with the binary meant to run on macOS.
- Now you are all set. You can now run DeathStar using the following

```
./dist/deathstar -conf [conf.yml] -deploy -local -zip-file-path deathstar.zip

```

The ```-local``` CLI option tells DeathStar that we want to use local zip file, hence DeathStar will not try to create a zip of its own. The ```-zip-file-path```
simply takes the path to the zip package which we want to use.

## Features to be added

As already mentioned, DeathStar is in development. DeathStar should support the following features:

- The AWS role creation for the lambda function is manual right now. This can be automated with the provision of providing custom role by user if required.
- Conducting attacks in parallel. This will result in calling the created lambda function in parallel. For this we might also have to introduce an option
  for configuring the lambda function concurrency.
- VPC support: This is required for testing services which are private to a VPC and not behind public IP. DeathStar should be able to spawn the lambda function
  in the desired subnet where the target service is running.
- Right now the only supported backend is serverless lambda functions. However lambda function might not always be suitable for generating huge traffic. Parallel
  running function (as mentioned in second point) may be usefull, still we may require alternate backends. DeathStar should support alternate backends like:
  
  -- EC2 instances: DeathStar should be able to spawn EC2 instances of desired type. We might be requiring higher number of vCPUs when generating huge
     SSL traffic.
     
  -- Provision for integrating with Kubernetes compute: DeathStar can spawn pods dynamically, run the attack and then destroy the pods. This way we can utilise a
     kubernetes cluster if already present.
     
     
