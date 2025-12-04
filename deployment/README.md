# Quizora Production Deployment

This directory contains all Docker configuration files for deploying Quizora to production.

## 📁 Directory Structure

```
deployment/
├── docker-compose.yml          # Main orchestration file
├── backend.Dockerfile          # Backend Go application
├── frontend.Dockerfile         # Frontend Next.js application
├── .env.example               # Environment variables template
├── mysql/
│   ├── my.cnf                 # MySQL configuration
│   └── init/                  # SQL initialization scripts
├── redis/
│   └── redis.conf             # Redis configuration
├── nginx/
│   ├── nginx.conf             # Nginx reverse proxy config
│   └── ssl/                   # SSL certificates directory
└── scripts/
    ├── deploy.sh              # Deployment automation script
    └── backup.sh              # Database backup script
```

## 🚀 Quick Start

### 1. Initial Setup on VPS

```bash
# Clone repository
git clone https://github.com/yourusername/quizora.git
cd quizora/deployment

# Copy and configure environment variables
cp .env.example .env
nano .env  # Edit with your production values

# Make scripts executable
chmod +x scripts/*.sh
```

### 2. Build and Start Services

```bash
# Build all images
docker compose build

# Start all services
docker compose up -d

# View logs
docker compose logs -f
```

### 3. Verify Deployment

```bash
# Check service status
docker compose ps

# Test backend health
curl http://localhost:8080/api/v1/health

# Test frontend
curl http://localhost:3000

# View specific service logs
docker compose logs backend
docker compose logs frontend
docker compose logs mysql
docker compose logs redis
```

## 🔧 Management Commands

### Service Control

```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# Restart specific service
docker compose restart backend

# View logs
docker compose logs -f [service_name]

# Execute commands in container
docker compose exec backend sh
docker compose exec mysql mysql -u root -p
```

### Database Operations

```bash
# Backup database
docker compose exec mysql mysqldump -u root -p quizora > backup.sql

# Restore database
docker compose exec -T mysql mysql -u root -p quizora < backup.sql

# Access MySQL shell
docker compose exec mysql mysql -u root -p quizora
```

### Redis Operations

```bash
# Access Redis CLI
docker compose exec redis redis-cli -a your_redis_password

# Check Redis info
docker compose exec redis redis-cli -a your_redis_password INFO

# Monitor Redis commands
docker compose exec redis redis-cli -a your_redis_password MONITOR
```

## 🔒 SSL Certificate Setup

### Using Let's Encrypt

```bash
# Install Certbot
apt install certbot

# Obtain certificates
certbot certonly --standalone -d medecole.com -d www.medecole.com
certbot certonly --standalone -d api.medecole.com

# Copy certificates to nginx directory
mkdir -p nginx/ssl/medecole.com
mkdir -p nginx/ssl/api.medecole.com

cp /etc/letsencrypt/live/medecole.com/fullchain.pem nginx/ssl/medecole.com/
cp /etc/letsencrypt/live/medecole.com/privkey.pem nginx/ssl/medecole.com/
cp /etc/letsencrypt/live/api.medecole.com/fullchain.pem nginx/ssl/api.medecole.com/
cp /etc/letsencrypt/live/api.medecole.com/privkey.pem nginx/ssl/api.medecole.com/

# Restart nginx
docker compose restart nginx
```

### Auto-renewal Setup

```bash
# Add cron job for auto-renewal
crontab -e

# Add this line:
0 3 * * * certbot renew --quiet --post-hook "docker compose -f /home/deployer/quizora/deployment/docker-compose.yml restart nginx"
```

## 📊 Monitoring

### Health Checks

```bash
# Backend health
curl https://api.medecole.com/api/v1/health

# Frontend health
curl https://medecole.com

# Check all container health
docker compose ps
```

### Resource Usage

```bash
# Container stats
docker stats

# Disk usage
docker system df

# Network usage
docker compose exec backend sh -c "netstat -tulpn"
```

## 🔄 Deployment Workflow

### Automated Deployment

```bash
# Run deployment script
cd /home/deployer/quizora/deployment
./scripts/deploy.sh
```

### Manual Deployment

```bash
# 1. Pull latest code
git pull origin main

# 2. Rebuild images
docker compose build

# 3. Stop old containers
docker compose down

# 4. Start new containers
docker compose up -d

# 5. Check logs
docker compose logs -f
```

## 💾 Backup & Restore

### Automated Backups

```bash
# Setup daily backup cron job
crontab -e

# Add this line (runs at 2 AM daily):
0 2 * * * /home/deployer/quizora/deployment/scripts/backup.sh >> /home/deployer/backups/backup.log 2>&1
```

### Manual Backup

```bash
./scripts/backup.sh
```

### Restore from Backup

```bash
# Restore MySQL
gunzip < /home/deployer/backups/mysql/quizora_20250104_020000.sql.gz | \
    docker compose exec -T mysql mysql -u root -p quizora

# Restore Redis
docker cp /home/deployer/backups/redis/redis_20250104_020000.rdb quizora-redis:/data/dump.rdb
docker compose restart redis
```

## 🐛 Troubleshooting

### Container Won't Start

```bash
# View detailed logs
docker compose logs [service_name]

# Check container status
docker compose ps

# Inspect container
docker inspect quizora-[service_name]
```

### Database Connection Issues

```bash
# Test MySQL connectivity
docker compose exec backend sh -c "nc -zv mysql 3306"

# Check MySQL logs
docker compose logs mysql

# Verify credentials
docker compose exec mysql mysql -u root -p -e "SHOW DATABASES;"
```

### Performance Issues

```bash
# Check resource usage
docker stats

# View MySQL slow queries
docker compose exec mysql cat /var/log/mysql/slow-query.log

# Monitor Redis
docker compose exec redis redis-cli -a password INFO stats
```

## 🔐 Security Checklist

- [ ] Change all default passwords in `.env`
- [ ] Use strong JWT secret (64+ characters)
- [ ] Configure firewall (UFW)
- [ ] Setup SSL certificates
- [ ] Enable Redis authentication
- [ ] Restrict MySQL to internal network
- [ ] Regular security updates
- [ ] Monitor logs for suspicious activity
- [ ] Backup encryption (optional)

## 📝 Environment Variables

See `.env.example` for all required environment variables. Key variables:

- `DB_PASSWORD`: MySQL root password
- `REDIS_PASSWORD`: Redis authentication password
- `JWT_SECRET`: JWT signing secret
- `GOOGLE_CLIENT_ID`: Google OAuth client ID
- `GOOGLE_CLIENT_SECRET`: Google OAuth secret

## 📞 Support

For issues or questions:

- Check logs: `docker compose logs -f`
- Review health checks: `docker compose ps`
- Contact: [your-email@example.com]

## 📜 License

[Your License Here]
