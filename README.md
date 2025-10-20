# Transline Test - Microservices with Distributed Tracing

Микросервисная архитектура с gRPC, REST API и распределённой трассировкой через OpenTelemetry/Jaeger.

## Архитектура

```
Client → Envoy (HTTP:8080) → shipment-service → customer-service (gRPC:9090) → PostgreSQL
                                      ↓                    ↓
                              PostgreSQL            PostgreSQL
```

**Envoy** — точка входа для внешних REST-запросов. Межсервисное gRPC-взаимодействие происходит напрямую внутри Docker-сети.

## Быстрый старт

```bash
docker-compose up
```

## API примеры

### Создать отгрузку
```bash
curl -X POST http://localhost:8080/api/v1/shipments \
  -H "Content-Type: application/json" \
  -d '{"route":"ALMATY→ASTANA","price":120000,"customer":{"idn":"990101123456"}}'
```

### Получить отгрузку
```bash
curl http://localhost:8080/api/v1/shipments/<id>
```

## Трассировка

Открыть Jaeger UI: **http://localhost:16686**

Выберите сервис `shipment-service` → **Find Traces**

### Структура трейса
```
📍 envoy-proxy (ingress)
  └─ 📍 shipment-service (HTTP handler)
     └─ 📍 shipment-service (gRPC client)
        └─ 📍 customer-service (gRPC server)
           └─ 🗄️ Database operations
```

Полная цепочка: **REST → Envoy → shipment-service → gRPC → customer-service → DB**

## Сервисы

- **shipment-service** (HTTP:8080) — REST API для управления отгрузками
- **customer-service** (gRPC:9090) — gRPC сервис для работы с клиентами
- **envoy** (HTTP:8080) — API Gateway и прокси
- **jaeger** (UI:16686) — визуализация распределённых трейсов
- **otel-collector** — сбор и экспорт телеметрии
- **postgres** — две базы данных (для каждого сервиса)

## Технологии

- Go 1.24
- gRPC + Protocol Buffers
- OpenTelemetry (otelhttp, otelgrpc)
- Jaeger Tracing
- Envoy Proxy
- PostgreSQL
- Docker Compose
