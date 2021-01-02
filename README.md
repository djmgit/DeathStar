# DeathStar

DeathStar is a tool to loadtest web based services and rest APIs in a easy, automated and quick way without spending much time on infrastructure setup for load generation. DeathStar expects a configuration file with the details of the attack targets and it takes care of the infrastructure, spawns an ephemeral infrastreucture to carry out the loadtest on the target, generates the result and then tears down the infrastructure. Hence, the developer or the operations engineer does not have to think about the infrastructure for carrying out the attack/test. All that is exepcted from the user is to provide the config file.

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
