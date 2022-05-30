kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.2.0/deploy/static/provider/cloud/deploy.yaml

kubectl get pods --namespace=ingress-nginx

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

kubectl expose deployment ${serviceName}

kubectl create ingress demo-localhost --class=nginx \
  --rule=${url}/*=${serviceName}:${exposePort}

kubectl port-forward --namespace=ingress-nginx service/ingress-nginx-controller 8080:80
