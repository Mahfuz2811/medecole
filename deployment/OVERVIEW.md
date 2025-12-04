# Deployment Directory

This folder structure is now set up for production deployment:

```
deployment/
├── docker-compose.yml          # Main orchestration (MySQL, Redis, Backend, Frontend, Nginx)
├── backend.Dockerfile          # Go backend multi-stage build
├── frontend.Dockerfile         # Next.js frontend optimized build
├── .env.example               # Environment variables template
│
├── mysql/
│   ├── my.cnf                 # MySQL performance tuning
│   └── init/                  # Place .sql initialization scripts here
│
├── redis/
│   └── redis.conf             # Redis configuration with AOF persistence
│
├── nginx/
│   ├── nginx.conf             # Reverse proxy for medecole.com & api.medecole.com
│   └── ssl/                   # Place SSL certificates here
│
└── scripts/
    ├── deploy.sh              # Automated deployment script
    └── backup.sh              # Daily backup script (MySQL + Redis)
```

## 🎯 Complete Stack

**All services run in Docker containers:**

- ✅ **MySQL 8.4** - Database with optimized configuration
- ✅ **Redis 7** - Cache with password authentication & persistence
- ✅ **Go Backend** - API server (multi-stage Alpine build)
- ✅ **Next.js Frontend** - Web app (standalone mode)
- ✅ **Nginx** - Reverse proxy with SSL termination

## 📦 What's Included

1. **Production-ready Dockerfiles**

   - Multi-stage builds for minimal image size
   - Security: non-root users
   - Health checks for all services
   - Optimized layer caching

2. **Docker Compose**

   - Service orchestration
   - Network isolation (quizora-network)
   - Volume persistence (data survives container restart)
   - Health check dependencies
   - Environment variable injection

3. **MySQL Configuration**

   - InnoDB buffer pool tuning
   - Connection pooling (200 max connections)
   - Slow query logging
   - Binary logging for backup
   - UTF8MB4 character set

4. **Redis Configuration**

   - Password authentication
   - AOF + RDB persistence
   - Memory limit (512MB)
   - LRU eviction policy
   - Optimized for caching

5. **Nginx Configuration**

   - SSL/TLS termination
   - HTTP to HTTPS redirect
   - Gzip compression
   - Rate limiting (API & login endpoints)
   - Security headers (HSTS, XSS, etc.)
   - Caching for static assets

6. **Automation Scripts**
   - `deploy.sh`: Pull code → Build → Deploy with health checks
   - `backup.sh`: Backup MySQL + Redis + Logs daily

## 🚀 Next Steps

1. **Update .env file** with production credentials
2. **Obtain SSL certificates** from Let's Encrypt
3. **Push to GitHub** repository
4. **Deploy to VPS** following the deployment plan

Refer to `README.md` for detailed deployment instructions.
