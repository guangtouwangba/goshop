# GoShop - Independent E-commerce Engine

[English](README.md) | [中文](README-CN.md)


GoShop is an open-source e-commerce platform based on microservices architecture, developed in Go. This project aims to provide a comprehensive, flexible, and extensible e-commerce solution.

## Project Architecture

GoShop adopts a Monorepo structure with the following microservices:

- **User Service (user-service)**: Handles user registration, login, profile management, etc.
- **Product Service (product-service)**: Manages products, categories, SKUs, etc.
- **Inventory Service (inventory-service)**: Manages inventory, stock movement, etc.
- **Order Service (order-service)**: Handles order creation, management, and shopping cart
- **Payment Service (payment-service)**: Handles payment integration and transaction records
- **Marketing Service (marketing-service)**: Manages coupons, promotions, membership, etc.
- **CMS Service (cms-service)**: Handles blogs, static pages, etc.
- **Shipping Service (shipping-service)**: Manages shipping methods, shipping fee calculation, etc.
- **API Gateway (gateway-service)**: Handles request routing, authentication, etc.
- **Authentication Service (auth-service)**: Handles JWT authentication, permissions, etc.
- **Admin Service (admin-service)**: Handles backend management functions

## Technology Stack

- **Programming Language**: Go
- **Web Framework**: Gin/Echo
- **RPC Framework**: gRPC
- **ORM**: GORM
- **Database**: PostgreSQL, Redis, Elasticsearch/Meilisearch
- **Message Queue**: NATS/RabbitMQ
- **Configuration Management**: Viper
- **Logging**: Zap
- **Authentication**: JWT
- **Containerization**: Docker, Kubernetes

## Project Features

- Microservices architecture, service decoupling, independent scaling
- Efficient inter-service communication based on gRPC
- RESTful API and/or GraphQL interface
- Comprehensive user authentication and permission control
- Rich product management features, supporting multiple SKUs
- Complete order management and shopping cart functionality
- Integration with mainstream payment gateways
- Flexible promotion and marketing capabilities
- Content management functionality
- Shipping and delivery management
- High performance, high availability design
- Comprehensive test coverage
- Observability support (Prometheus, Grafana, Jaeger)

## Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/goshop.git
cd goshop

# Start dependency services with Docker Compose
docker-compose up -d

# Build and run
make build
make run
```

## Contribution Guidelines

Contributions to code, reporting issues, or suggesting new features are welcome. Please refer to the [Contribution Guidelines](CONTRIBUTING.md) for details.

## License

This project is licensed under the [MIT License](LICENSE).
