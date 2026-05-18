param(
    [string]$BaseUrl = "http://localhost:80",
    [string]$Rate = "10/s",
    [string]$Duration = "30s",
    [int]$Contacts = 20,
    [string]$TargetsPath = "",
    [string]$OutputBin = "",
    [string]$OutputHtml = "",
    [switch]$SkipPrepare
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($TargetsPath)) {
    $TargetsPath = Join-Path $repoRoot "targets.txt"
}
if ([string]::IsNullOrWhiteSpace($OutputBin)) {
    $OutputBin = Join-Path $repoRoot "normal.bin"
}
if ([string]::IsNullOrWhiteSpace($OutputHtml)) {
    $OutputHtml = Join-Path $repoRoot "normal.html"
}

$vegetaCommand = Get-Command vegeta -ErrorAction SilentlyContinue
if ($vegetaCommand) {
    $vegetaPath = $vegetaCommand.Source
} else {
    $localVegeta = Join-Path $repoRoot "vegeta.exe"
    if (Test-Path $localVegeta) {
        $vegetaPath = $localVegeta
    } else {
        throw "vegeta was not found in PATH or repo root. Install it first: go install github.com/tsenart/vegeta@latest"
    }
}

if (-not $SkipPrepare) {
    & (Join-Path $PSScriptRoot "prepare.ps1") `
        -BaseUrl $BaseUrl `
        -Contacts $Contacts `
        -TargetsPath $TargetsPath
}

& $vegetaPath attack `
    -targets $TargetsPath `
    -rate $Rate `
    -duration $Duration `
    -output $OutputBin

& $vegetaPath report $OutputBin
& $vegetaPath plot -title "Balembala send load test" $OutputBin | Out-File -FilePath $OutputHtml -Encoding utf8

Write-Host "Saved binary results: $((Resolve-Path $OutputBin).Path)"
Write-Host "Saved html plot: $((Resolve-Path $OutputHtml).Path)"
