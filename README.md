# Balembala

# WARNING: переименование rebalancer -> retry-processor

### [Дизайн Web-интерфейса](https://www.figma.com/design/Bg2BRmB2xAySUbpqfYpdBS/balembala)

### POSTMAN: https://warped-sunset-469734.postman.co/workspace/Team-Workspace~0790e497-4c26-4ddc-8790-15e8f1d3240f/collection/43133821-da74b742-d17c-4e35-a322-3dc7a04ff8b6?action=share&creator=43133821
### Структура

```
notification-system/
├── api-gateway/
│   ├── Dockerfile
│   └── cmd/
├── auth-service/
│   ├── Dockerfile
│   └── cmd/
├── sender-service/
│   ├── Dockerfile
│   └── cmd/
├── notification-service/
│   ├── Dockerfile
│   └── cmd/
├── contacts-templates-service/
│   ├── Dockerfile
│   └── cmd/
├── retry-processor-service/
│   ├── Dockerfile
│   └── cmd/
├── web/
│   ├── Dockerfile
│   └── cmd/
├── docker-compose.yml   
└── README.md      
```

### Ветки

- **main** — основная ветка, содержащая стабильную версию кода, которая идет на прод.

### Процедура добавления кода

1. Перед началом работы над сервисом и его функционалом создать ветку `dev/`.
2. После завершения работы над веткой создайте Merge Reques в ветку `main`.
3. Название MR должно соответствовать названию ветки, в описании хорошо бы писать что сделано.
