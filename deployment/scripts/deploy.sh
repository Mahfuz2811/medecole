#!/bin/bash

# Quizora Production Deployment Script
# Usage: ./deploy.sh

set -e  # Exit on error

echo "🚀 Starting Quizora Deployment..."

# Configuration
DEPLOYMENT_DIR="/home/deployer/quizora"
BACKUP_DIR="/home/deployer/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored messages
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Check if running as correct user
if [ "$USER" != "deployer" ]; then
    print_error "This script must be run as 'deployer' user"
    exit 1
fi

# Navigate to deployment directory
cd "$DEPLOYMENT_DIR" || exit 1

# Step 1: Pull latest code from GitHub
print_warning "Pulling latest code from GitHub..."
git pull origin main || {
    print_error "Failed to pull latest code"
    exit 1
}
print_success "Code updated successfully"

# Step 2: Backup current database
print_warning "Creating database backup..."
mkdir -p "$BACKUP_DIR"
docker exec quizora-mysql mysqldump -u root -p"$DB_PASSWORD" quizora > "$BACKUP_DIR/quizora_$DATE.sql" || {
    print_error "Database backup failed"
    exit 1
}
print_success "Database backed up to $BACKUP_DIR/quizora_$DATE.sql"

# Step 3: Build new Docker images
print_warning "Building Docker images..."
cd deployment
docker compose build --no-cache || {
    print_error "Docker build failed"
    exit 1
}
print_success "Docker images built successfully"

# Step 4: Stop old containers gracefully
print_warning "Stopping old containers..."
docker compose down || {
    print_error "Failed to stop containers"
    exit 1
}
print_success "Old containers stopped"

# Step 5: Start new containers
print_warning "Starting new containers..."
docker compose up -d || {
    print_error "Failed to start containers"
    print_warning "Rolling back..."
    docker compose down
    exit 1
}
print_success "New containers started"

# Step 6: Wait for services to be healthy
print_warning "Waiting for services to become healthy..."
sleep 10

# Check backend health
for i in {1..30}; do
    if docker exec quizora-backend wget --spider -q http://localhost:8080/api/v1/health; then
        print_success "Backend is healthy"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "Backend health check failed"
        docker compose logs backend
        exit 1
    fi
    sleep 2
done

# Check frontend health
for i in {1..30}; do
    if docker exec quizora-frontend wget --spider -q http://localhost:3000; then
        print_success "Frontend is healthy"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "Frontend health check failed"
        docker compose logs frontend
        exit 1
    fi
    sleep 2
done

# Step 7: Clean up old Docker images
print_warning "Cleaning up old Docker images..."
docker image prune -f
print_success "Cleanup completed"

# Step 8: Display running containers
print_warning "Current running containers:"
docker compose ps

# Step 9: Display logs (last 50 lines)
print_warning "Recent logs:"
docker compose logs --tail=50

print_success "🎉 Deployment completed successfully!"
print_warning "Monitor logs with: cd $DEPLOYMENT_DIR/deployment && docker compose logs -f"
