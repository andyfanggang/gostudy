kubectl create secret tls wm-motor-secret -n istio-system  --cert=1_wm-motor.com_bundle.crt  --key=2_wm-motor.com.key
kubectl apply -f httpserver.yaml
kubectl apply -f spec.yaml
[root@master01 gateway]# kubectl get svc -n istio-system -o wide
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                      AGE     SELECTOR
istio-egressgateway    ClusterIP      192.89.178.15   <none>        80/TCP,443/TCP                                                               3d19h   app=istio-egressgateway,istio=egressgateway
istio-ingressgateway   LoadBalancer   192.89.91.15    <pending>     15021:30548/TCP,80:32025/TCP,443:32413/TCP,31400:30494/TCP,15443:32029/TCP   3d19h   app=istio-ingressgateway,istio=ingressgateway
istiod                 ClusterIP      192.89.1.131    <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP                                        3d19h   app=istiod,istio=pilot

[root@master01 gateway]# curl --resolve httpsserver.wm-motor.com:443:192.89.91.15 https://httpsserver.wm-motor.com/healthz
<h1>200</h1>
