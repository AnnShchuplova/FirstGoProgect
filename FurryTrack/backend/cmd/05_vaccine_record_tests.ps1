$API_URL = "http://localhost:8080"


# Аутентификация тестового пользователя
Write-Host "АУТЕНТИФИКАЦИЯ ПОЛЬЗОВАТЕЛЯ"
$loginBody = @{
    email = "lerachapurina@mail.ru"
    password = "12345678"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$API_URL/api/login" `
        -Method Post `
        -Body $loginBody `
        -ContentType "application/json"
    
    $token = $loginresponse.data.token
    Write-Host "Успешная аутентификация" -ForegroundColor Green
    $headers = @{
        "Authorization" = "Bearer $token"
    }
}
catch {
    Write-Host "ОШИБКА АВТОРИЗАЦИИ:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit
}

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

# Тест добавления записи о вакцинации
Write-Host "ТЕСТ ДОБАВЛЕНИЯ ЗАПИСИ О ВАКЦИНАЦИИ"
$vaccineRecordBody = @{
    pet_id = $petToUpdate.data.pet_id[0]
    vaccine_name = "New"
    date = (Get-Date -Format "yyyy-MM-dd")
    clinic = "Тест"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$API_URL/api/pets/$($petToUpdate.data.pet_id[0])/vaccine-records" `
        -Method Post `
        -Body $vaccineRecordBody `
        -Headers $headers `
        -ContentType "application/json"
    
    Write-Host "ЗАПИСЬ УСПЕШНО ДОБАВЛЕНА:" -ForegroundColor Green
    $response | Format-List | Out-Host
    
    $testRecordId = $response.id
}
catch {
    Write-Host "ОШИБКА ДОБАВЛЕНИЯ ЗАПИСИ:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
    exit
}

# Тест получения истории вакцинации
Write-Host "ТЕСТ ПОЛУЧЕНИЯ ИСТОРИИ ВАКЦИНАЦИИ"
try {
    $historyResponse = Invoke-RestMethod -Uri "$API_URL/api/pets/$($petToUpdate.data.pet_id[0])/vaccine-records" `
        -Method Get `
        -Headers $headers
    
    Write-Host "ИСТОРИЯ ВАКЦИНАЦИИ:" -ForegroundColor Green
    $historyResponse | Format-Table -AutoSize | Out-Host
    
    # Проверка что новая запись есть в истории
    $newRecordFound = $historyResponse | Where-Object { $_.id -eq $testRecordId }
    if (-not $newRecordFound) {
        Write-Host "ОШИБКА: Новая запись не найдена в истории" -ForegroundColor Red
    }
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ ИСТОРИИ:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}