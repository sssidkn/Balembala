param(
    [string]$BaseUrl = "http://localhost:80",
    [string]$Email = "",
    [string]$Password = "123456",
    [string]$Username = "vegeta",
    [int]$Contacts = 20,
    [string]$TargetsPath = "",
    [string]$StatePath = ""
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($TargetsPath)) {
    $TargetsPath = Join-Path $repoRoot "targets.txt"
}
if ([string]::IsNullOrWhiteSpace($StatePath)) {
    $StatePath = Join-Path $PSScriptRoot "state.json"
}

function Invoke-Json {
    param(
        [string]$Method,
        [string]$Url,
        [hashtable]$Headers = @{},
        $Body = $null
    )

    $params = @{
        Method = $Method
        Uri = $Url
        ContentType = "application/json"
    }

    if ($Headers.Count -gt 0) {
        $params.Headers = $Headers
    }
    if ($null -ne $Body) {
        if ($Body -is [string]) {
            $params.Body = $Body
        } else {
            $params.Body = $Body | ConvertTo-Json -Depth 20 -Compress
        }
    }

    try {
        return Invoke-RestMethod @params
    } catch {
        $details = $_.Exception.Message
        if ($_.ErrorDetails -and $_.ErrorDetails.Message) {
            $details = $_.ErrorDetails.Message
        }
        throw "$Method $Url failed: $details"
    }
}

function Get-RequiredProperty {
    param(
        $Object,
        [string]$Property,
        [string]$Source
    )

    if ($null -eq $Object -or $null -eq $Object.$Property) {
        throw "Cannot read '$Property' from $Source response."
    }

    return $Object.$Property
}

$stamp = Get-Date -Format "yyyyMMddHHmmss"
if ([string]::IsNullOrWhiteSpace($Email)) {
    $Email = "vegeta-$stamp@example.com"
}

$registerBody = @{
    email = $Email
    password = $Password
    username = $Username
}

try {
    $auth = Invoke-Json -Method "POST" -Url "$BaseUrl/api/register" -Body $registerBody
} catch {
    $loginBody = @{
        email = $Email
        password = $Password
    }
    $auth = Invoke-Json -Method "POST" -Url "$BaseUrl/api/login" -Body $loginBody
}

$token = Get-RequiredProperty -Object $auth -Property "token" -Source "auth"
$headers = @{ Authorization = "Bearer $token" }

$template = Invoke-Json `
    -Method "POST" `
    -Url "$BaseUrl/api/templates" `
    -Headers $headers `
    -Body @{
        title = "vegeta-template-$stamp"
        message = "Load test message $stamp"
    }

$templateObject = Get-RequiredProperty -Object $template -Property "template" -Source "create template"
$templateId = [int64](Get-RequiredProperty -Object $templateObject -Property "id" -Source "create template")

$contactIds = @()
for ($i = 1; $i -le $Contacts; $i++) {
    $contact = Invoke-Json `
        -Method "POST" `
        -Url "$BaseUrl/api/contacts" `
        -Headers $headers `
        -Body @{
            name = "Vegeta Contact $i"
            email = "vegeta-$stamp-$i@example.com"
        }

    $contactObject = Get-RequiredProperty -Object $contact -Property "contact" -Source "create contact"
    $contactIds += [int64](Get-RequiredProperty -Object $contactObject -Property "id" -Source "create contact")
}

Invoke-Json `
    -Method "POST" `
    -Url "$BaseUrl/api/template/$templateId" `
    -Headers $headers `
    -Body @{ contactsId = $contactIds } | Out-Null

$target = @"
POST $BaseUrl/api/send/$templateId
Authorization: Bearer $token
Content-Type: application/json

"@

Set-Content -Path $TargetsPath -Value $target -Encoding ascii

$state = @{
    baseUrl = $BaseUrl
    email = $Email
    templateId = $templateId
    contacts = $contactIds.Count
    targetsPath = (Resolve-Path $TargetsPath).Path
    createdAt = (Get-Date).ToString("o")
}

$state | ConvertTo-Json -Depth 20 | Set-Content -Path $StatePath -Encoding ascii

Write-Host "Prepared vegeta target:"
Write-Host "  user: $Email"
Write-Host "  templateId: $templateId"
Write-Host "  contacts: $($contactIds.Count)"
Write-Host "  targets: $((Resolve-Path $TargetsPath).Path)"
