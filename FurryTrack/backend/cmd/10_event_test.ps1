# Конфигурация
$API_URL = "http://localhost:8080"
$testEmail = "lerachapurina@mail.ru"
$testPassword = "12345678"



Write-Host "=== АУТЕНТИФИКАЦИЯ ===" -ForegroundColor Cyan
$authBody = @{
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$API_URL/api/login" `
        -Method Post `
        -Body $authBody `
        -ContentType "application/json" `
        -ErrorAction Stop

    $token = $loginResponse.data.token
    # $userId = $loginResponse.data.user.ID
    Write-Host "Успешная аутентификация" -ForegroundColor Green
}
catch {
    Write-Host "ОШИБКА АУТЕНТИФИКАЦИИ:" -ForegroundColor Red
    Write-Host $_.Exception.Message
    exit
}

# Получаем список питомцев пользователя
Write-Host "`n=== ПОЛУЧЕНИЕ ПИТОМЦЕВ ПОЛЬЗОВАТЕЛЯ ===" -ForegroundColor Cyan
try {
    $petsResponse = Invoke-RestMethod -Uri "$API_URL/api/pets" `
        -Method Get `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    $petId = $null
    if ($petsResponse.data.Count -gt 0) {
        $petId = $petsResponse.data.pet_id[0]
        Write-Host "Найден питомец (ID: $petId)" -ForegroundColor Green
    } else {
        Write-Host "У пользователя нет питомцев" -ForegroundColor Yellow
        exit
    }
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ ПИТОМЦЕВ:" -ForegroundColor Red
    Write-Host $_.Exception.Message
    exit
}

# Создание события
Write-Host "`n=== СОЗДАНИЕ СОБЫТИЯ ===" -ForegroundColor Cyan
$eventBody = @{
    pet_id = $petId
    type = "vet_visit"
    title = "Плановый осмотр"
    description = "Ежегодный осмотр у ветеринара"
    date = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
    location = "Ветеринарная клиника 'Доктор Айболит'"
    cost = 1500
} | ConvertTo-Json

try {
    $eventResponse = Invoke-RestMethod -Uri "$API_URL/api/events" `
        -Method Post `
        -Body $eventBody `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop


    $eventId = $eventResponse.event.ID
    Write-Host "Событие создано (ID: $eventId)" -ForegroundColor Green
    Write-Host "Детали события:" -ForegroundColor Cyan
    $eventResponse.data | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА СОЗДАНИЯ СОБЫТИЯ:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
    exit
}

# Получение событий питомца
Write-Host "`n=== ПОЛУЧЕНИЕ СОБЫТИЙ ПИТОМЦА ===" -ForegroundColor Cyan
try {
    $eventsResponse = Invoke-RestMethod -Uri "$API_URL/api/pets/$petId/events" `
        -Method Get `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    $count = $eventsResponse.data.Count
    Write-Host "УСПЕХ: Получено $count событий" -ForegroundColor Green
    if ($count -gt 0) {
        Write-Host "Последнее событие питомца:" -ForegroundColor Cyan
        $eventsResponse.data[0] | Format-List | Out-Host
    }
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ СОБЫТИЙ:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
}

Write-Host "`n=== ТЕСТИРОВАНИЕ ЗАВЕРШЕНО ===" -ForegroundColor Green