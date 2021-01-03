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
Before I countinue further, I would like to mention that although deathstar is usable, it is still under development and features are being added. Currently it
only supports AWS lambda as a compute backend. 
```
DeathStar takes a configuration file, provisions a lambda function along with the handler, orchestrates the loadtest on the configured targets, displays the
recorded metrics and finally cleans up the lambda function, thereby allowing you to carry out your loadtest in a fully automated and serverless manner.

In the config file you can mention multiple targets and for each target you can provide the attack configurations like what should be the rate of the attack, that
is, what should be the traffic qps, what is the endpoint that DeathStar will be hitting with the mentioned qps, what should be the HTTP method (right now it only
supports GET, but will support other methods along with body and url parameters in future) and for how long will the test be continued. It rquires some other
details as well like the lambda memory size to be allocated, timeout, aws region etc, more on in subsequent sections.
DeathStar will parse the config, accordingly create a lambda fucntion, invoke the lambda function for carrying out the attack and finally will destroy the lambda
function. 
This would be a nice time to mention that DeathStar uses <a href="https://github.com/LazyWolves/"> Vegeta </a> under the hood to create the load. Vegeta
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








