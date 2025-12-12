#!/bin/bash

# Medecole Backup Script
# Usage: ./backup.sh
# Cron: 0 2 * * * /home/deployer/medecole/deployment/scripts/backup.sh

set -e

# Configuration
BACKUP_DIR="/home/deployer/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Load environment variables
source /home/deployer/medecole/deployment/.env

echo "Starting backup at $(date)"

# Create backup directory
mkdir -p "$BACKUP_DIR"/{mysql,redis,logs}

# Backup MySQL database
echo "Backing up MySQL database..."
docker exec medecole-mysql mysqldump \
    -u root \
    -p"$DB_PASSWORD" \
    --single-transaction \
    --routines \
    --triggers \
    --events \
    medecole | gzip > "$BACKUP_DIR/mysql/medecole$DATE.sql.gz"

echo "MySQL backup completed: medecole$DATE.sql.gz"

# Backup Redis data
echo "Backing up Redis data..."
docker exec medecole-redis redis-cli -a "$REDIS_PASSWORD" BGSAVE
sleep 5
docker cp medecole-redis:/data/dump.rdb "$BACKUP_DIR/redis/redis_$DATE.rdb"

echo "Redis backup completed: redis_$DATE.rdb"

# Backup application logs
echo "Backing up application logs..."
docker cp medecole-backend:/app/logs "$BACKUP_DIR/logs/backend_$DATE"
docker cp medecole-nginx:/var/log/nginx "$BACKUP_DIR/logs/nginx_$DATE"

echo "Logs backup completed"

# Clean up old backups (older than retention days)
echo "Cleaning up old backups..."
find "$BACKUP_DIR/mysql" -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
find "$BACKUP_DIR/redis" -name "*.rdb" -mtime +$RETENTION_DAYS -delete
find "$BACKUP_DIR/logs" -type d -mtime +$RETENTION_DAYS -exec rm -rf {} +

echo "Backup cleanup completed"

# Calculate backup size
BACKUP_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
echo "Total backup size: $BACKUP_SIZE"

echo "Backup completed successfully at $(date)"
