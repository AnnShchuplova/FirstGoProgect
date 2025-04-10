# Конфигурация
$API_URL = "http://localhost:8080"
$testEmail = "lerachapurina@mail.ru"
$testPassword = "12345678"

Write-Host "=== ПРОВЕРКА ПОЛЬЗОВАТЕЛЯ В БАЗЕ ===" -ForegroundColor Cyan
$userInDb = & 'C:\Program Files\PostgreSQL\17\bin\psql.exe' -U postgres -d furrytrack -c "SELECT id, email, password_hash, deleted_at FROM users WHERE email = '$testEmail';" -t
Write-Host $userInDb

if (-not $userInDb -or $userInDb -match "0 rows") {
    Write-Host "ОШИБКА: Пользователь $testEmail не найден!" -ForegroundColor Red
    exit
}

Write-Host "`n=== ПОПЫТКА АУТЕНТИФИКАЦИИ ===" -ForegroundColor Cyan
$authBody = @{
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

try {
    $loginresponse = Invoke-RestMethod -Uri "$API_URL/api/login" `
        -Method Post `
        -Body $authBody `
        -ContentType "application/json" `
        -ErrorAction Stop

    $loginresponse | Format-List | Out-Host
    
    $token = $loginresponse.data.token
    Write-Host "Токен получен" -ForegroundColor Green
}
catch {
    Write-Host "ОШИБКА 401: Неверный логин или пароль" -ForegroundColor Red
    
    Write-Host "ОШИБКА ЗАПРОСА:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode)" -ForegroundColor Red
    Write-Host "Сообщение: $($_.Exception.Message)" -ForegroundColor Red
    
    if ($_.ErrorDetails) {
        Write-Host "Тело ошибки:" -ForegroundColor Yellow
        $_.ErrorDetails.Message | Out-Host
    }
    exit
}


$petData = @{
    name = "Arty"
    type = "Cat"  
    breed = "Шотландский вислоухий"
    birth_date = "2015-06-03"  
    owner_id = $loginresponse.data.user.ID  
} | ConvertTo-Json -Depth 3

try {
    $petresponse = Invoke-RestMethod -Uri "$API_URL/api/pets" `
        -Method Post `
        -Body $petData `
        -Headers @{ 
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "`nПитомец успешно создан" -ForegroundColor Green
    $petresponse | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА СОЗДАНИЯ ПИТОМЦА:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
}