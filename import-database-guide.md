# Database Import Guide for Railway MySQL

## Method 1: Export from cPanel phpMyAdmin

### Step 1: Export Database
1. Go to cPanel → phpMyAdmin
2. Select database `menarane_lilly`
3. Click "Export" tab
4. Choose "Quick" export method
5. Select "SQL" format
6. Click "Go" to download SQL file

### Step 2: Import to Railway
1. In Railway dashboard, find your MySQL service
2. Click "Connect" or "Query" button
3. Paste your SQL content or upload the SQL file
4. Execute the import

## Method 2: Using MySQL Command Line

### Step 1: Get Railway MySQL Connection Details
From Railway dashboard → MySQL service → Variables:
- MYSQL_HOST
- MYSQL_USER
- MYSQL_PASSWORD
- MYSQL_DATABASE
- MYSQL_PORT

### Step 2: Export from cPanel (Command Line)
```bash
mysqldump -h imogiri.idweb.host -u menarane_lilly -p menarane_lilly > backup.sql
```

### Step 3: Import to Railway
```bash
mysql -h $MYSQL_HOST -u $MYSQL_USER -p $MYSQL_DATABASE < backup.sql
```

## Method 3: Using MySQL Workbench

1. Download MySQL Workbench
2. Create new connection with Railway MySQL details
3. Open your exported SQL file
4. Execute the script

## Important Notes

- Make sure to export/import in the correct order (tables with foreign keys last)
- Check for any MySQL version compatibility issues
- Some cPanel-specific features might not work on Railway
- Test the import with a small subset first
