#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const environments = {
  local: {
    // MySQL Database Configuration
    MYSQL_HOST: 'localhost',
    MYSQL_PORT: '3306',
    MYSQL_USER: 'root',
    MYSQL_PASSWORD: '',
    MYSQL_DATABASE: 'iqgncnzy_skripsi',
    
    // Alternative variable names (fallback)
    DB_HOST: 'localhost',
    DB_PORT: '3306',
    DB_USER: 'root',
    DB_PASSWORD: '',
    DB_NAME: 'iqgncnzy_skripsi',
    
    // Mikrotik Configuration - Local
    MIKROTIK_HOST: '10.10.9.203',
    MIKROTIK_PORT: '22',
    MIKROTIK_USERNAME: 'rnd',
    MIKROTIK_PASSWORD: 'rnd@123',
    
    // JWT Secret
    JWT_SECRET: 'your-local-secret-key-here',
    JWT_SECRET_KEY: 'your-local-secret-key-here',
    
    // Server Configuration
    PORT: '3001'
  },
  production: {
    // MySQL Database Configuration - Production (Railway/cPanel)
    MYSQL_HOST: 'your-production-host',
    MYSQL_PORT: '3306',
    MYSQL_USER: 'your-production-user',
    MYSQL_PASSWORD: 'your-production-password',
    MYSQL_DATABASE: 'iqgncnzy_skripsi',
    
    // Mikrotik Configuration - Production
    MIKROTIK_HOST: '103.148.18.52',
    MIKROTIK_PORT: '22',
    MIKROTIK_USERNAME: 'your-production-mikrotik-user',
    MIKROTIK_PASSWORD: 'your-production-mikrotik-password',
    
    // JWT Secret
    JWT_SECRET: 'your-production-secret-key-here',
    JWT_SECRET_KEY: 'your-production-secret-key-here',
    
    // Server Configuration
    PORT: '3001'
  }
};

const env = process.argv[2];

if (!env || !environments[env]) {
  console.log('Usage: node switch-env.js <local|production>');
  console.log('Available environments:', Object.keys(environments).join(', '));
  process.exit(1);
}

const envContent = Object.entries(environments[env])
  .map(([key, value]) => `${key}=${value}`)
  .join('\n') + '\n';

fs.writeFileSync('.env', envContent);
console.log(`✅ Switched to ${env} environment`);
console.log(`📝 Created .env file for backend`);
console.log('\nConfiguration:');
console.log(`  - Database: ${environments[env].MYSQL_HOST}:${environments[env].MYSQL_PORT}/${environments[env].MYSQL_DATABASE}`);
console.log(`  - Mikrotik: ${environments[env].MIKROTIK_HOST}:${environments[env].MIKROTIK_PORT}`);
console.log(`  - Server Port: ${environments[env].PORT}`);

