$API_URL = "http://localhost:8080"
# $testUsername = "admin"
$testEmail = "admin@furrytrack.ru"
$testPassword = "SecureAdmin123!"



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
    Write-Host "ОШИБКА" -ForegroundColor Red
    
    if ($_.ErrorDetails) {
        Write-Host "Тело ошибки:" -ForegroundColor Yellow
        $_.ErrorDetails.Message | Out-Host
    }
    exit
}



Write-Host "`n=== СОЗДАНИЕ ВАКЦИНЫ ===" -ForegroundColor Cyan
$vaccineData = @{
    name = "New_2"
    description = "new vaccine"
    duration_days = 365  # Срок действия в днях
    #is_mandatory = $true
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$API_URL/api/vaccines" `
        -Method Post `
        -Body $vaccineData `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "Вакцина успешно создана!" -ForegroundColor Green
    $response | Format-List | Out-Host
    #$vaccineId = $response.data.id
}
catch {
    Write-Host "ОШИБКА СОЗДАНИЯ ВАКЦИНЫ:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    Write-Host $_.Exception.Message
    
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
    exit
}

# Получение списка вакцин
Write-Host "`n=== ПОЛУЧЕНИЕ СПИСКА ВАКЦИН ===" -ForegroundColor Cyan
try {
    $vaccinesResponse = Invoke-RestMethod -Uri "$API_URL/api/vaccines" `
        -Method Get `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "Список вакцин получен (количество: $($vaccinesResponse.data.count))" -ForegroundColor Green
    
    # Выводим таблицу с вакцинами
    $vaccinesResponse.data.vaccines | Format-Table -Property id, name, duration_days, description -AutoSize | Out-Host
    
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ СПИСКА ВАКЦИН:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    Write-Host $_.Exception.Message
    
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
    exit
}


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


Write-Host "`n=== УДАЛЕНИЕ ВАКЦИНЫ ===" -ForegroundColor Cyan
try {
    $deleteResponse = Invoke-RestMethod -Uri "$API_URL/api/vaccines/$vaccineId" `
        -Method Delete `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "Вакцина успешно удалена. Ответ сервера:" -ForegroundColor Green
    $deleteResponse | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА УДАЛЕНИЯ ВАКЦИНЫ:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    Write-Host $_.Exception.Message
    exit
}

# Проверка удаления
Write-Host "`n=== ПРОВЕРКА УДАЛЕНИЯ ===" -ForegroundColor Cyan
try {
    $checkResponse = Invoke-RestMethod -Uri "$API_URL/api/vaccines/$vaccineId" `
        -Method Get `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "ОШИБКА: Вакцина все еще существует" -ForegroundColor Red
    $checkResponse | Format-List | Out-Host
    exit
}
catch {
    if ($_.Exception.Response.StatusCode -eq 404) {
        Write-Host "Вакцина успешно удалена (404 Not Found)" -ForegroundColor Green
    }
    else {
        Write-Host "Ошибка:" -ForegroundColor Red
        Write-Host $_.Exception.Message
        exit
    }
}

Write-Host "`n=== ТЕСТ УСПЕШНО ЗАВЕРШЕН ===" -ForegroundColor Green