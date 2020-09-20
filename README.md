# AWS EgressIP Operator

> What man is a man who does not make the world better.
>
> -- Balian, Kingdom of Heaven


## Abstract

This operator automates the assignment of egressIPs to namespaces. It is inspired from the [egressip-ipam-operator of the
GitHub Red Hat CoP](https://github.com/redhat-cop/egressip-ipam-operator) project.

It is incompatible. Instead of annotations to the namespace resource here EgressIPs are fully managed by 
CustomResources.



## Needed permissions for this operator
This operator needs some AWS permissions to do its work. These have to be handled via instance-profiles. The needed
permissions and the reasoning are:

AWS Permission | Reasoning
---------------|-----------------------------------
EC2:DescribeInstances | Getting information about the instances (tags, networking interfaces).
EC2:AssignPrivateIpAddresses | Manage the IP addresses of the instances.
EC2:UnassignPrivateIpAddresses | Manage the IP addresses of the instances.


## Deploying the Operator

This is a cluster-level operator that you can deploy in any namespace, `egress-ip-operator` is recommended.
If you need to pin the operator to special nodes (like the OCP infranodes), please use the namespace node-selector
annotation to do that. May be helpful in restricting the AWS permissions to only a few nodes.

**Note:** *Create the namespace with `openshift.io/node-selector: ''` in order to deploy to master nodes. Or select the
 nodes you gave the needed AWS permissions.*

## License
The license for the software is Apache License 2.0. 

## Note from the author
I started this operator since I developed a similar operator with the old  operator-sdk and while migrating I decided to
start from the scratch to improve the architecture.

If someone is interested in getting it faster, we may team up. I'm open for that. But be warned: I want to do it 
_right_. So no short cuts to get faster. And be prepared for some basic discussions about the architecture or software 
design :-).

---
Bensheim, 2020-09-19