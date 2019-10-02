# Cloudflare Host Mapper

This command will pull all the hosts for every load balancer in the namespace and then write the dns mapping to Cloudflare to allow us to run under their proxy whilst also not having to pay billions of pounds for wildcard proxy support.

This client directly interacts with both the Kubernetes API and the Cloudflare API and as such requires access to the kubernetes namespace that the container is located within and internet access.

Make sure to set the Cloudflare settings or else this client won't work 😉

**This does not work outside of a K8s cluster**

## Setup

This presumes you are setting up the role to be used in the `platform` namespace. If you are using `default` then this is not required.

Create the role to read ingresses in the default namespace\
`kubectl create clusterrole ingress-reader --verb=get --verb=list --resource=ingress --namespace=default`

Create the role binding in the cluster\
`kubectl create clusterrolebinding platform-default-ingress-binding --clusterrole=ingress-reader --user=system:serviceaccount:platform:default`

## Running

To run this command you will need to pass the config in via either the `--config` or the environment variable `CONFIG`.
