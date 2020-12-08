# clusterip-webhook

validating and mutating webhook to manage the Kubernetes Services Cluster IPs


## How to use it

The admin defines the IP Ranges objects with the subnets, the Services ClusterIPs will be assigned
from this IP ranges.

Currently, it only supports one IPRange per IPFamily with the same value defined in the apiserver

TODO:

1. Move Service IP Range configuration out of the apiserver
2. Allow to use multiples IPranges
  - Modify Services to reference IPRange Object

## How it works

When an user creates a Kubernetes Service of Type ClusterIP, NodePort or LoadBalancer, it can
specify the ClusterIP or the apiserver will allocate one free.

### Create, Update

If the ClusterIP is set by the user, the webhook validates it belongs to the range and is free

If the ClusterIP is not set, the webhook assigns one free from the range

The ClusterIP is inmutable after creation, but when the Service Type changes

### Delete

When a Service is Deleted, the controller will deallocate the ClusterIP assigned from the IPRange
object once the Service Delete event is received (Not when the Delete request is seen)

