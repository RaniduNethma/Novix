# Novix
Where stories come alive.
## Project Architecture

```
/Novix
  │
  ├── /services                          		# Backend Microservices
  │     ├── /user-service                		(Java/Spring Boot)
  │     │     └── Auth, profiles, watch history
  │     ├── /content-service             		(Java/Spring Boot)
  │     │     └── Metadata, catalog, CQRS, search
  │     ├── /notification-service        		(Java/Spring Boot)
  │     │     └── Email, push, in-app
  │     ├── /recommendation-service      		(Java/Spring Boot)
  │     │     └── ML recommendations, trending
  │     ├── /streaming-service           		(Go)
  │     │     └── HLS/DASH delivery, CDN, Redis
  │     └── /processing-service          		(Go)
  │           └── FFmpeg, transcoding, S3, job workers
  │
  ├── /frontend                          		# Frontend Apps
  │     ├── /web                         		(Next.js — viewer facing)
  │     │     └── SSR, video player, HLS.js, watch UI
  │     └── /admin                       		(Next.js — internal dashboard)
  │           └── Content mgmt, analytics, moderation
  │
  ├── /packages                          		# Shared Packages (Frontend + Backend)
  │     ├── /ui                          		(Shared component library)
  │     ├── /types                       		(Shared TypeScript interfaces)
  │     ├── /api-client                  		(Shared backend communication)
  │     └── /proto                       		(Protobuf/gRPC contracts Java+Go)
  │
  ├── /infrastructure                    		# Infra & DevOps Configs
  │     ├── /kong                        		(API Gateway — routes, plugins, auth)
  │     ├── /consul                      		(Service discovery + health checks)
  │     ├── /kafka                       		(Message broker + topic configs)
  │     ├── /prometheus-grafana          		(Metrics + dashboards)
  │     ├── /docker                      		(Dockerfiles per service)
  │     └── /k8s                         		(Kubernetes manifests)
  │
  ├── docker-compose.yml                 		# Local dev full stack
  └── Makefile                           		# Unified commands
```
