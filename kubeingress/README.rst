========================
Yeti Ingress Controllers
========================

Quick summary: you probably only care about
api-ingress-ing-template.yaml.

Hello and welcome. Ingress controllers are a feature of Kubernetes that,
effectively, provide a standardized mechanism for passing runtime
kubernetes configuration information down into a load balancer /
webserver / router container.

In our case we use an ingress controller to perform a mapping of
services to paths in HTTP requests. We can easily write rules that say
things like: Map requests to /api/v1/network to the network service
running here in k8s.

The mapping is made by the ingress controller, in our case a bit of code
sitting somewhere down near nginx. We don't really care about the
details here.

The ingress-template file contained in this directory specifies the
mappings we care about. Define paths in there, like the existing ones,
and the ingress controller will take care of routing requests made to
the API at the path you specify to your pod.

Requests from the web now more or less travel like so:

- User browser
- An ELB (hosted yeti only)
- The ingress-controller (nginx) k8s service listening on port 30948
- An api-ingress-lb pod
- A pod providing the proper service (yes, direct from nginx to the
  pod).

Sorta funky is that the nginx controllers today appear to poll for
changes to service membership and update their config once, what appears
to be, every 5 seconds. So a little slow.

Components
==========

Ingress controllers are actually made up of several components.

We talk a lot about an ingress resource in k8s, this is a resource of
type Ingress. An Ingress is just a spec, a definition saying "forward
requests matching these rules to these services." When you look at our
ingress, api-ingress-ing-template.yaml, those are the rules used to
define ingressing.

The question then to ask is how the hell do I get requests to a thing
that is parsing those rules? This is handled by a service.

Here we have api-ingress-svc-template.yaml. This service picks up a
NodePort of 30948 (our api node port) and waits for traffic. Traffic is
forwarded to pods matching the selector of `k8s-app: api-ingress-lb`.

The next logical component is then whatever provides k8s-app of
api-ingress-lb. That is in api-ingress-lb-dp-template.yaml. This
specifies a deployment that maintains an nginx pod that understands
ingress. This the crux of ingress. This pod is simply nginx + a little
program that talks to the k8s API and scans the Ingress resources for
rules. Based on the rules it finds it generates an nginx configuration
file to implement the ingress.

Ingress rules specify "if you match this condition" forward to "this
service."

Presently our rules only reference the yetiapi service. The yetiapi
service is defined over in yeti/services/api/k8s/api-svc-template.yaml.

Debugging
=========

If things start going sideways the nginx pod is reasonably
introspectable.

Access logs from nginx are viewable via Sumo or kubectl logs. In
addition, mixed in there, are logs from the controller that is updating
config based on k8s changes.

If that's not good enough you can `kubectl exec -it <nginx pod name>
bash` and look at /etc/nginx/nginx.conf. You can also `apt update` and
then `apt get <pkg>` to do things like tcpdump. Be super careful doing
this in prod since tcpdump will show real JWTs among other sensitive
things.

