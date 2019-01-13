## Network Policy Tests
In Kubernetes by default there is no network plugin is available[[1]](https://kubernetes.io/docs/concepts/services-networking/network-policies/#prerequisites), We need to have one of the network plugin (like calico, cilium, weave-net, flannel etc.. ) needed in the cluster.

A network policy is a specification of how groups of pods are allowed to communicate with each other and other network endpoints.

To demonstrate Network policies, we are going to use Sample App (service1) which we can test the connectivity across the cluster by adding/removing network policies to it. This sample app has below endpoints.

```
Endpoints: 

/ping - To check if health of service
        Ex - http://localhost:3000/ping 
/req - To call an external endpoint passed through header
        Ex - http://localhost:3000/req -H 'X-Req-URL: http://google.com/'
```

Later in the section we are going to run another service (service2) nginx, so we can test the pod to pod connectivity.

## Pre-Requisites
1. Setup EKS cluster. Refer [docs](https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html)
2. Setup Ingress controller and get the public ELB endpoint. Refer [docs](https://github.com/kubernetes/ingress-nginx) (for following example, we can use thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com)
3. Setup Cilium. Refer [docs](https://cilium.readthedocs.io/en/stable/kubernetes/install/eks/)
3. Deploy the sample app (service1)
```
kubectl apply -f sample-app.yaml
```
4. Deploy the nginx service (service2)
```
kubectl apply -f nginx.yaml
```

## Test 1 - Allow All
If we dont have any NetworkPolicies, the default behaviour is to allow all Ingress and Egress across cluster. 

1. Verify if service1 is reachable through Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/ping -H 'Host: helloworld.service'
```

2. Verify if service1 is reachable through another pod in cluster
```
kubectl run -it --rm --restart=Never centos --image=centos:7 bash
curl http://hw-service.hw.svc.cluster.local:3000/ping --connect-timeout 10
```
> So we have validated the Ingress rules of service1

3. Verify if service1 can reach Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H 'X-Req-URL: http://google.com/'
```
> So we have Validated the Egress rules of service1

4. Verify if service can reach any other service in cluster
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H "X-Req-URL: http://my-nginx.default.svc.cluster.local"
```
> So we have validated the inter-cluster rules of service1

## Test 2 - Deny All
To deny all ingress/egress of sample-app (service1) running in namespace `hw`. This policy can also be applied to all pods in that namespace by specifying `podSelector: {}`
```
kubectl apply -f deny-all.yaml
```

1. Verify if service1 is reachable through Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/ping -H 'Host: helloworld.service'
```
> You should be able to access it throuh Internet, Check here for explanation

2. Verify if service1 is reachable through another pod in cluster
```
kubectl run -it --rm --restart=Never centos --image=centos:7 bash
curl http://hw-service.hw.svc.cluster.local:3000/ping --connect-timeout 10
```
> So we have validated the Ingress rules of service1 is blocking traffic

3. Verify if service1 can reach Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H 'X-Req-URL: http://google.com/'
```
> So we have Validated the Egress rules of service1 is blocking traffic

4. Verify if service can reach any other service in cluster
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H "X-Req-URL: http://my-nginx.default.svc.cluster.local"
```
> So we have validated the inter-cluster rules of service1 is blocking traffic

## Test 3 - Allow Ingress, Deny Ingress for all
To allow ingress and deny all ingress of sample-app (service1) running in namespace `hw`
```
kubectl delete networkpolicy deny-all -n hw
kubectl apply -f deny-egress.yaml
```

1. Verify if service1 is reachable through Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/ping -H 'Host: helloworld.service'
```
> You should be able to access it throuh Internet, Check here for explanation

2. Verify if service1 is reachable through another pod in cluster
```
kubectl run -it --rm --restart=Never centos --image=centos:7 bash
curl http://hw-service.hw.svc.cluster.local:3000/ping --connect-timeout 10
```
> You should be able to get the PONG response. So we have validated the Ingress rules of service1 is allowing the traffic

3. Verify if service1 can reach Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H 'X-Req-URL: http://google.com/'
```
> You shouldn't be able to get the response. So we have Validated the Egress rules of service1 is blocking traffic

4. Verify if service can reach any other service in cluster
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H "X-Req-URL: http://my-nginx.default.svc.cluster.local"
```
> You shouldn't be able to get the response. So we have validated the inter-cluster rules of service1 is blocking traffic

## Test 4 - Allow Ingress, Deny Egress for all except Kube-dns
```
kubectl delete cnp deny-egress -n hw
kubectl apply -f deny-egress-except-dns.yaml
```

1. Verify if service1 is reachable through Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/ping -H 'Host: helloworld.service'
```
> You should be able to access it throuh Internet, Check here for explanation

2. Verify if service1 is reachable through another pod in cluster
```
kubectl run -it --rm --restart=Never centos --image=centos:7 bash
curl http://hw-service.hw.svc.cluster.local:3000/ping --connect-timeout 10
```
> You should be able to get the PONG response. So we have validated the Ingress rules of service1 is allowing the traffic

3. Verify if service1 can reach Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H 'X-Req-URL: http://google.com/'
```
> You shouldn't be able to get the response. So we have Validated the Egress rules of service1 is blocking traffic

4. Verify if service can reach any other service in cluster
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H "X-Req-URL: http://my-nginx.default.svc.cluster.local"
```
> You shouldn't be able to get the response. So we have validated the inter-cluster rules of service1 is blocking traffic

5. Verify DNS resolutions
```
# Dry run without any network policies, are we able to access the endpoints
# kubectl run -it --rm --restart=Never dnsutils --image=tutum/dnsutils bash
# dig google.com @8.8.8.8
# dig my-nginx.default.svc.cluster.local

# Instead of exec'ing into the actual pod, We are mimicing the service1 by adding the labels and namespace, which enforces the network policies
kubectl run -it --rm --restart=Never --labels="app=hw" -n hw tinytools --image=giantswarm/tiny-tools sh
dig google.com @8.8.8.8 #reaching out internet for DNS resolution
dig my-nginx.default.svc.cluster.local #reaching cluster DNS (coreDNS) for resolution
curl http://my-nginx.default.svc.cluster.local #Cant reach because there is no rule to allow
```

## Test 5 - Allow Ingress, Deny egress for all except kube-dns and nginx service
```
kubectl delete cnp deny-egress-except-kube-dns -n hw
kubectl apply -f deny-egress-except-nginx.yaml
```

1. Verify if service1 is reachable through Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/ping -H 'Host: helloworld.service'
```
> You should be able to access it throuh Internet, Check here for explanation

2. Verify if service1 is reachable through another pod in cluster
```
kubectl run -it --rm --restart=Never centos --image=centos:7 bash
curl http://hw-service.hw.svc.cluster.local:3000/ping --connect-timeout 10
```
> You should be able to get the PONG response. So we have validated the Ingress rules of service1 is allowing the traffic

3. Verify if service1 can reach Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H 'X-Req-URL: http://google.com/'
```
> You shouldn't be able to get the response. So we have Validated the Egress rules of service1 is blocking traffic

4. Verify if service can reach any other service in cluster
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H "X-Req-URL: http://my-nginx.default.svc.cluster.local"
```
> You should be able to get the response. So we have validated the inter-cluster rules of service1 is allowing traffic

5. Verify DNS resolutions
```
# Dry run without any network policies, are we able to access the endpoints
# kubectl run -it --rm --restart=Never tinytools --image=giantswarm/tiny-tools sh
# dig google.com @8.8.8.8
# dig my-nginx.default.svc.cluster.local
# curl http://my-nginx.default.svc.cluster.local

# Instead of exec'ing into the actual pod, We are mimicing the service1 by adding the labels and namespace, which enforces the network policies
kubectl run -it --rm --restart=Never --labels="app=hw" -n hw tinytools --image=giantswarm/tiny-tools sh
dig google.com @8.8.8.8 
dig my-nginx.default.svc.cluster.local
curl http://my-nginx.default.svc.cluster.local
```

## Test 6 - Allow Ingress, Allow egress except cluster endpoints
```
kubectl delete cnp deny-egress-except-nginx -n hw
kubectl apply -f deny-cluster-access.yaml
```

1. Verify if service1 is reachable through Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/ping -H 'Host: helloworld.service'
```
> You should be able to access it throuh Internet, Check here for explanation

2. Verify if service1 is reachable through another pod in cluster
```
kubectl run -it --rm --restart=Never centos --image=centos:7 bash
curl http://hw-service.hw.svc.cluster.local:3000/ping --connect-timeout 10
```
> You should be able to get the PONG response. So we have validated the Ingress rules of service1 is allowing the traffic

3. Verify if service1 can reach Internet
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H 'X-Req-URL: http://google.com/'
```
> You should be able to get the response. So we have Validated the Egress rules of service1 is allowing traffic

4. Verify if service can reach any other service in cluster
```
curl http://thrash98-sbx-or2-k8s-pub-lb-0-672520883.us-west-2.elb.amazonaws.com/req -H 'Host: helloworld.service' -H "X-Req-URL: http://my-nginx.default.svc.cluster.local"
```
> You should be able to get the response. So we have validated the inter-cluster rules of service1 is allowing traffic

5. Verify Network policy enforcement
```
# Dry run without any network policies, are we able to access the endpoints
# kubectl run -it --rm --restart=Never tinytools --image=giantswarm/tiny-tools sh
# dig google.com @8.8.8.8
# dig my-nginx.default.svc.cluster.local
# curl http://my-nginx.default.svc.cluster.local
# curl http://100.77.82.80 # IP address of Nginx service
# curl http://10.10.13.82:6000 # Accessing the contour running on worker host, although it returns 404 network connection is established
# curl https://172.20.94.97 # Accessing the kubernetes dashboard service endpoint

# Instead of exec'ing into the actual pod, We are mimicing the service1 by adding the labels and namespace, which enforces the network policies
kubectl run -it --rm --restart=Never --labels="app=hw" -n hw tinytools --image=giantswarm/tiny-tools sh
dig google.com @8.8.8.8 # working because of 0.0.0.0/0 whitelisting
dig my-nginx.default.svc.cluster.local
curl http://my-nginx.default.svc.cluster.local
curl http://100.77.82.80
curl http://10.10.13.82:6000 --connect-timeout 10 # timesout because of 10.10.0.0/16 cidr rule
curl https://172.20.94.97 --connect-timeout 10 # timesout because of 172.20.0.0/16 cidr rule
```