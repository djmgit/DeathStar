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
recoreded metrics and finally cleans up the lambda function, thereby allowing you to carry out your loadtest in a fully automated and serverless manner.









