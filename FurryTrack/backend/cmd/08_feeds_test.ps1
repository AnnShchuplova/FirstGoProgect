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
    $userId = $loginResponse.data.user.ID
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
        Write-Host "У пользователя нет питомцев, пост будет создан без pet_id" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "ОШИБКА ПОЛУЧЕНИЯ ПИТОМЦЕВ:" -ForegroundColor Red
    Write-Host $_.Exception.Message
    .\test_feed.ps1 = $null
}
Write-Host $userId

# Создание тестового поста
Write-Host "`n=== СОЗДАНИЕ ТЕСТОВОГО ПОСТА ===" -ForegroundColor Cyan
$postData = @{
    AutorID = $userId
    content = "Тестовый пост для проверки ленты, маркет"
    post_type = "market"
    
}


$postData.Add("Pet_id", $petId)

$postData | Format-List | Out-Host

$postBody = $postData | ConvertTo-Json

try {
    $postResponse = Invoke-RestMethod -Uri "$API_URL/api/posts" `
        -Method Post `
        -Body $postBody `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    $postId = $postResponse.ID
    Write-Host $postResponse
    Write-Host "Пост создан (ID: $postId)" -ForegroundColor Green
    Write-Host "Детали поста:" -ForegroundColor Cyan
    $postResponse.data | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА СОЗДАНИЯ ПОСТА:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
    exit
}
Write-Host $postId

# Тестирование лент
$feedTests = @(
    @{Name = "Основная лента"; Endpoint = "/api/feed/main"},
    @{Name = "Лента продаж"; Endpoint = "/api/feed/market"},
    @{Name = "Лента подписок"; Endpoint = "/api/feed/following"}
)

foreach ($test in $feedTests) {
    Write-Host "`n=== ТЕСТИРУЕМ $($test.Name) ===" -ForegroundColor Cyan
    try {
        $response = Invoke-RestMethod -Uri "$API_URL$($test.Endpoint)" `
            -Method Get `
            -Headers @{
                "Authorization" = "Bearer $token"
                "Content-Type" = "application/json"
            } `
            -ErrorAction Stop

        $count = $response.data.Count
        Write-Host "УСПЕХ: Получено $count постов" -ForegroundColor Green
        if ($count -gt 0) {
            Write-Host "Последний пост в ленте:" -ForegroundColor Cyan
            $response.data[0] | Format-List | Out-Host
        }
    }
    catch {
        Write-Host "ОШИБКА:" -ForegroundColor Red
        Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $errorBody = $reader.ReadToEnd()
            Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
        }
    }
}

# Тестирование лайков
Write-Host "`n=== ТЕСТИРОВАНИЕ ЛАЙКОВ ===" -ForegroundColor Cyan
try {
    $likeResponse = Invoke-RestMethod -Uri "$API_URL/api/posts/$postId/like" `
        -Method Post `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "Лайк успешно поставлен" -ForegroundColor Green
    $likeResponse  | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА ЛАЙКА:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
}

# Тестирование комментариев
Write-Host "`n=== ТЕСТИРОВАНИЕ КОММЕНТАРИЕВ ===" -ForegroundColor Cyan
$commentBody = @{
    content = "Тестовый комментарий"
} | ConvertTo-Json

try {
    $commentResponse = Invoke-RestMethod -Uri "$API_URL/api/posts/$postId/comments" `
        -Method Post `
        -Body $commentBody `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -ErrorAction Stop

    Write-Host "Комментарий успешно добавлен" -ForegroundColor Green
    Write-Host "Детали комментария:" -ForegroundColor Cyan
    $commentResponse.data | Format-List | Out-Host
}
catch {
    Write-Host "ОШИБКА КОММЕНТАРИЯ:" -ForegroundColor Red
    Write-Host "Статус код: $($_.Exception.Response.StatusCode.value__)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "Тело ошибки: $errorBody" -ForegroundColor Yellow
    }
}

Write-Host "`n=== ТЕСТИРОВАНИЕ ЗАВЕРШЕНО ===" -ForegroundColor Green