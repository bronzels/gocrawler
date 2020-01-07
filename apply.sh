kubectl delete deployment gocrawlernews -n default
kubectl apply -f kubernetes/crawlnews/deployment.yaml
kubectl get deployment -n default
kubectl delete svc gocrawlernews -n default
kubectl apply -f kubernetes/crawlnews/service.yaml
kubectl get svc -n default
