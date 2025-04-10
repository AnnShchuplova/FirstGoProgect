# Конфигурация
$API_URL = "http://localhost:8080"
$testEmail = "lerachapurina@mail.ru"
$testPassword = "12345678"
$targetUserEmail = "mogemoge911@mail.ru" 

# Аутентификация основного пользователя
Write-Host "=== АУТЕНТИФИКАЦИЯ ОСНОВНОГО ПОЛЬЗОВАТЕЛЯ ===" -ForegroundColor Cyan
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
    $userId = $loginResponse.data.user.ID
    Write-Host "Токен получен для пользователя $testEmail" -ForegroundColor Green
}
catch {
    Write-Host "ОШИБКА АУТЕНТИФИКАЦИИ:" -ForegroundColor Red
    Write-Host $_.Exception.Message
    exit
}

# Получение ID целевого пользователя
Write-Host "`n=== ПОЛУЧЕНИЕ ID ЦЕЛЕВОГО ПОЛЬЗОВАТЕЛЯ ===" -ForegroundColor Cyan
try {
    $targetUserResponse = Invoke-RestMethod -Uri "$API_URL/api/users/email/$targetUserEmail" `
        -Method Get `
        -Headers @{ 
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    $targetUserId = $targetUserResponse.data.ID
    Write-Host "ID целевого пользователя: $targetUserId" -ForegroundColor Green
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ ЦЕЛЕВОГО ПОЛЬЗОВАТЕЛЯ:" -ForegroundColor Red
    Write-Host $_.Exception.Message
    exit
}

# Подписаться на пользователя
Write-Host "`n=== ПОДПИСКА НА ПОЛЬЗОВАТЕЛЯ ===" -ForegroundColor Cyan
try {
    $followResponse = Invoke-RestMethod -Uri "$API_URL/api/users/$targetUserId/follow" `
        -Method Post `
        -Headers @{ 
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "Успешно подписались на пользователя $targetUserEmail" -ForegroundColor Green
    $followResponse | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА ПОДПИСКИ:" -ForegroundColor Red
    Write-Host $_.Exception.Message
}

# Получить список подписок
Write-Host "`n=== ПОЛУЧЕНИЕ СПИСКА ПОДПИСОК ===" -ForegroundColor Cyan
try {
    $followingResponse = Invoke-RestMethod -Uri "$API_URL/api/users/$userId/following" `
        -Method Get `
        -Headers @{ 
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "Список подписок:" -ForegroundColor Green
    $followingResponse | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ ПОДПИСОК:" -ForegroundColor Red
    Write-Host $_.Exception.Message
}

# Получить список подписчиков
Write-Host "`n=== ПОЛУЧЕНИЕ СПИСКА ПОДПИСЧИКОВ ===" -ForegroundColor Cyan
try {
    $followersResponse = Invoke-RestMethod -Uri "$API_URL/api/users/$targetUserId/followers" `
        -Method Get `
        -Headers @{ 
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "Список подписчиков пользователя ${targetUserEmail}:" -ForegroundColor Green
    $followersResponse | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ ПОДПИСЧИКОВ:" -ForegroundColor Red
    Write-Host $_.Exception.Message
}