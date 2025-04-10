# Тест регистрации пользователя
$API_URL = "http://localhost:8080"

# Генерация тестовых данных
$timestamp = Get-Date -Format "yyyyMMddHHmmss"
$testUsername = "New_User_$timestamp"
$testEmail = "test_$timestamp@example.com"
$testPassword = "12345678Password"



Write-Host "=== ТЕСТ РЕГИСТРАЦИИ ===" -ForegroundColor Cyan
Write-Host "Имя пользователя: $testUsername"
Write-Host "Email: $testEmail"
Write-Host "Пароль: $testPassword"

$regBody = @{
    username = $testUsername
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$API_URL/api/register" `
        -Method Post `
        -Body $regBody `
        -ContentType "application/json"
    
    Write-Host "УСПЕШНАЯ РЕГИСТРАЦИЯ" -ForegroundColor Green
    $response | Format-List | Out-Host
    
    # Проверка в БД
    Write-Host "`nПроверка в базе данных:" -ForegroundColor Yellow
    & 'C:\Program Files\PostgreSQL\17\bin\psql.exe' -U postgres -d furrytrack -c "SELECT id, username, email FROM users WHERE email = '$testEmail';"
}
catch {
    Write-Host "ОШИБКА РЕГИСТРАЦИИ:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    
}

Write-Host "`nТест завершен. Сервер продолжает работать." -ForegroundColor Cyan