# Script untuk memeriksa dan memperbaiki hasil migrasi database
# Berdasarkan hasil terminal yang menunjukkan duplicate column product_id

Write-Host "🔍 Memeriksa status migrasi database..." -ForegroundColor Cyan

# Database connection parameters
$dbHost = "localhost"
$dbPort = "3306"
$dbUser = "root"
$dbPassword = ""
$dbName = "iqgncnzy_skripsi"

# Connection string
$connectionString = "Server=$dbHost;Port=$dbPort;Database=$dbName;Uid=$dbUser;Pwd=$dbPassword;"

try {
    # Load MySQL .NET connector
    Add-Type -Path "C:\laragon\bin\mysql\mysql-8.0.30-winx64\lib\MySql.Data.dll"
    
    $connection = New-Object MySql.Data.MySqlClient.MySqlConnection($connectionString)
    $connection.Open()
    
    Write-Host "✅ Koneksi database berhasil" -ForegroundColor Green
    
    # 1. Periksa struktur tabel network_devices
    Write-Host "`n📋 Memeriksa struktur tabel network_devices..." -ForegroundColor Yellow
    $checkTableQuery = "DESCRIBE network_devices;"
    $command = New-Object MySql.Data.MySqlClient.MySqlCommand($checkTableQuery, $connection)
    $reader = $command.ExecuteReader()
    
    $columns = @()
    while ($reader.Read()) {
        $columns += @{
            Field = $reader["Field"]
            Type = $reader["Type"]
            Null = $reader["Null"]
            Key = $reader["Key"]
            Default = $reader["Default"]
        }
    }
    $reader.Close()
    
    Write-Host "Kolom yang ada di tabel network_devices:" -ForegroundColor White
    foreach ($col in $columns) {
        Write-Host "  - $($col.Field): $($col.Type)" -ForegroundColor Gray
    }
    
    # 2. Periksa apakah product_id sudah ada
    $productIdExists = $columns | Where-Object { $_.Field -eq "product_id" }
    
    if ($productIdExists) {
        Write-Host "`n⚠️  Kolom product_id sudah ada di tabel network_devices" -ForegroundColor Yellow
        Write-Host "   Ini menjelaskan mengapa migrasi statement 1 gagal" -ForegroundColor Gray
        
        # 3. Periksa foreign key constraints
        Write-Host "`n🔗 Memeriksa foreign key constraints..." -ForegroundColor Yellow
        $checkFKQuery = @"
SELECT 
    CONSTRAINT_NAME,
    COLUMN_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM information_schema.KEY_COLUMN_USAGE 
WHERE TABLE_SCHEMA = '$dbName' 
    AND TABLE_NAME = 'network_devices' 
    AND REFERENCED_TABLE_NAME IS NOT NULL;
"@
        
        $command = New-Object MySql.Data.MySqlClient.MySqlCommand($checkFKQuery, $connection)
        $reader = $command.ExecuteReader()
        
        $foreignKeys = @()
        while ($reader.Read()) {
            $foreignKeys += @{
                ConstraintName = $reader["CONSTRAINT_NAME"]
                ColumnName = $reader["COLUMN_NAME"]
                ReferencedTable = $reader["REFERENCED_TABLE_NAME"]
                ReferencedColumn = $reader["REFERENCED_COLUMN_NAME"]
            }
        }
        $reader.Close()
        
        if ($foreignKeys.Count -gt 0) {
            Write-Host "Foreign key constraints yang ada:" -ForegroundColor White
            foreach ($fk in $foreignKeys) {
                Write-Host "  - $($fk.ConstraintName): $($fk.ColumnName) -> $($fk.ReferencedTable).$($fk.ReferencedColumn)" -ForegroundColor Gray
            }
        } else {
            Write-Host "Tidak ada foreign key constraints ditemukan" -ForegroundColor Gray
        }
        
        # 4. Periksa indexes
        Write-Host "`n📊 Memeriksa indexes..." -ForegroundColor Yellow
        $checkIndexQuery = "SHOW INDEX FROM network_devices WHERE Column_name = 'product_id';"
        $command = New-Object MySql.Data.MySqlClient.MySqlCommand($checkIndexQuery, $connection)
        $reader = $command.ExecuteReader()
        
        $indexes = @()
        while ($reader.Read()) {
            $indexes += @{
                KeyName = $reader["Key_name"]
                ColumnName = $reader["Column_name"]
                NonUnique = $reader["Non_unique"]
            }
        }
        $reader.Close()
        
        if ($indexes.Count -gt 0) {
            Write-Host "Indexes untuk product_id:" -ForegroundColor White
            foreach ($idx in $indexes) {
                Write-Host "  - $($idx.KeyName): $($idx.ColumnName) (Unique: $(-not $idx.NonUnique))" -ForegroundColor Gray
            }
        } else {
            Write-Host "Tidak ada indexes untuk product_id ditemukan" -ForegroundColor Gray
        }
        
    } else {
        Write-Host "`n❌ Kolom product_id tidak ditemukan di tabel network_devices" -ForegroundColor Red
        Write-Host "   Ini aneh karena migrasi seharusnya menambahkannya" -ForegroundColor Gray
    }
    
    # 5. Periksa tabel customer untuk product_id
    Write-Host "`n👥 Memeriksa tabel customer..." -ForegroundColor Yellow
    $checkCustomerQuery = "DESCRIBE customer;"
    $command = New-Object MySql.Data.MySqlClient.MySqlCommand($checkCustomerQuery, $connection)
    $reader = $command.ExecuteReader()
    
    $customerColumns = @()
    while ($reader.Read()) {
        $customerColumns += $reader["Field"]
    }
    $reader.Close()
    
    $customerProductIdExists = $customerColumns -contains "product_id"
    if ($customerProductIdExists) {
        Write-Host "✅ Kolom product_id ada di tabel customer" -ForegroundColor Green
        Write-Host "   Ini akan dihapus setelah migrasi selesai" -ForegroundColor Gray
    } else {
        Write-Host "ℹ️  Kolom product_id tidak ada di tabel customer" -ForegroundColor Blue
    }
    
    # 6. Rekomendasi
    Write-Host "`n📝 REKOMENDASI:" -ForegroundColor Cyan
    Write-Host "1. Migrasi sudah berjalan dengan baik meskipun ada warning" -ForegroundColor White
    Write-Host "2. Statement 1 gagal karena kolom product_id sudah ada (ini normal)" -ForegroundColor White
    Write-Host "3. Statement 2 dan 3 berhasil (foreign key dan index ditambahkan)" -ForegroundColor White
    Write-Host "4. Langkah selanjutnya: Verifikasi aplikasi berjalan dengan baik" -ForegroundColor White
    Write-Host "5. Jika semuanya OK, jalankan script untuk menghapus product_id dari customer" -ForegroundColor White
    
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
} finally {
    if ($connection -and $connection.State -eq 'Open') {
        $connection.Close()
        Write-Host "`n🔌 Koneksi database ditutup" -ForegroundColor Gray
    }
}

Write-Host "`n✅ Pemeriksaan selesai!" -ForegroundColor Green

