# GoShop - 独立电商引擎

GoShop 是一个基于微服务架构的开源电商平台，使用 Go 语言开发。本项目旨在提供一个全面、灵活且可扩展的电商解决方案。

## 项目架构

GoShop 采用 Monorepo 结构，包含以下微服务：

- **用户服务 (user-service)**：处理用户注册、登录、个人资料管理等
- **商品服务 (product-service)**：管理商品、分类、SKU 等
- **库存服务 (inventory-service)**：管理库存、库存变动等
- **订单服务 (order-service)**：处理订单创建、管理和购物车功能
- **支付服务 (payment-service)**：处理支付集成和交易记录
- **营销服务 (marketing-service)**：管理优惠券、促销活动、会员等
- **内容管理服务 (cms-service)**：处理博客、静态页面等
- **物流服务 (shipping-service)**：管理配送方式、运费计算等
- **API 网关 (gateway-service)**：处理请求路由、认证等
- **认证服务 (auth-service)**：处理 JWT 认证、权限等
- **管理后台服务 (admin-service)**：处理后台管理功能

## 技术栈

- **编程语言**：Go
- **Web 框架**：Gin/Echo
- **RPC 框架**：gRPC
- **ORM 框架**：GORM
- **数据库**：PostgreSQL、Redis、Elasticsearch/Meilisearch
- **消息队列**：NATS/RabbitMQ
- **配置管理**：Viper
- **日志管理**：Zap
- **认证**：JWT
- **容器化**：Docker、Kubernetes

## 项目特性

- 微服务架构，服务解耦，独立扩展
- 基于 gRPC 的高效服务间通信
- RESTful API 和/或 GraphQL 接口
- 完善的用户认证和权限控制
- 丰富的商品管理功能，支持多 SKU
- 完整的订单管理和购物车功能
- 集成主流支付网关
- 灵活的促销和营销能力
- 内容管理功能
- 物流配送管理
- 高性能、高可用设计
- 全面的测试覆盖
- 可观测性支持 (Prometheus、Grafana、Jaeger)

## 快速开始

```bash
# 克隆仓库
git clone https://github.com/yourusername/goshop.git
cd goshop

# 使用 Docker Compose 启动依赖服务
docker-compose up -d

# 构建并运行
make build
make run
```

## 贡献指南

欢迎对代码贡献、问题报告或新功能建议。详情请参阅 [贡献指南](CONTRIBUTING.md)。

## 许可证

本项目采用 [MIT 许可证](LICENSE)。 