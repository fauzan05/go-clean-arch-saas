#!/bin/bash

echo "🚀 Go SaaS Starter - Setup Script"
echo "=================================="

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    exit 1
fi

# Load .env
source .env

# Create database if not exists
echo "📦 Creating database..."
mysql -u${DB_USERNAME} -p${DB_PASSWORD} -h${DB_HOST} -e "CREATE DATABASE IF NOT EXISTS ${DB_NAME};" 2>/dev/null || {
    echo "⚠️  Could not create database. Please create manually: ${DB_NAME}"
}

# Run migrations
echo "📦 Running migrations..."
make migrate-up

# Seed data
echo "🌱 Seeding initial data..."
make seed

echo ""
echo "✅ Setup complete!"
echo "📖 Run 'make dev' to start development server"
