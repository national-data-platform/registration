# NDP

## Docker

Build+Start
```
docker compose up --build -d
```

## Kubernetes
Create secrets from postgress:
```
kubectl create secret generic osdf-env --from-env-file=.env-k8s
```

Create postgres pvc:
```
kubectl apply -f kubernetes/postgres-pvc.yaml
```

Deploy postgres:
```
kubectl apply -f kubernetes/postgres.yaml
```

Deploy api service:
```
kubectl apply -f kubernetes/api.yaml
```

## API Usage
### Local
Download test.txt file from osdf/pelican caches:
```
curl -o test.txt -X POST http://127.0.0.1:8080/download -d '{"Name": "/ospool/uc-shared/public/OSG-Staff/validation/test.txt"}'
```

Upload a test.txt file to osdf/pelican caches:
```
curl -X POST http://127.0.0.1:8080/upload -H "Content-Type: multipart/form-data" -F 'file=@test.txt' -F 'name=/ndp/burnpro3d/test.txt' -F 'token=token'
```

### NRP
Download test.txt file from osdf/pelican caches:
```
curl -o test.txt -X POST https://osdf-api.nrp-nautilus.io/download -d '{"Name": "/ospool/uc-shared/public/OSG-Staff/validation/test.txt"}'
```