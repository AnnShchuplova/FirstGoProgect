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
    Write-Host "ОШИБКА ЗАПРОСА:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode)" -ForegroundColor Red
    Write-Host "Сообщение: $($_.Exception.Message)" -ForegroundColor Red
    
    if ($_.ErrorDetails) {
        Write-Host "Тело ошибки:" -ForegroundColor Yellow
        $_.ErrorDetails.Message | Out-Host
    }
    exit
}

# Тест получения питомцев
Write-Host "`n=== ТЕСТ ПОЛУЧЕНИЯ ПИТОМЦЕВ ===" -ForegroundColor Cyan
try {
    $petsResponse = Invoke-RestMethod -Uri "$API_URL/api/pets" `
        -Method Get `
        -Headers @{ 
            "Authorization" = "Bearer $token"
        } `
        -ErrorAction Stop

    Write-Host "Найдено питомцев: $($petsResponse.Count)" -ForegroundColor Green
    $petsResponse | ConvertTo-Json -Depth 3 | Out-Host
    $petToUpdate = $petsResponse[0] 
    # Write-Host $petToUpdate.data
    Write-Host "Будем обновлять питомца: $($petToUpdate.data.Name[0]) (ID: $($petToUpdate.data.pet_id[0]))" -ForegroundColor Cyan
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ ПИТОМЦЕВ:" -ForegroundColor Red
    Write-Host $_.Exception.Message
}

# Тест обновления питомца
if ($petToUpdate.data.pet_id) {
    Write-Host "`n=== ТЕСТ ОБНОВЛЕНИЯ ПИТОМЦА ===" -ForegroundColor Cyan
    $updateData = @{
        name = "Arty_new_name"
        breed = "Шотландская вислоухая"
    } | ConvertTo-Json

    try {
        $updateResponse = Invoke-RestMethod -Uri "$API_URL/api/pets/$($petToUpdate.data.pet_id[0])" `
            -Method Put `
            -Body $updateData `
            -Headers @{ 
                "Authorization" = "Bearer $token"
                "Content-Type" = "application/json"
            } `
            -ErrorAction Stop

        Write-Host "Питомец обновлен" -ForegroundColor Green
        $updateResponse | Format-List | Out-Host
    }
    catch {
        Write-Host "ОШИБКА ОБНОВЛЕНИЯ:" -ForegroundColor Red
        Write-Host $_.Exception.Message
    }
}