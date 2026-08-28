# One-click install for Quantis on Windows.
# Usage (elevated recommended for Program Files):
#   powershell -ExecutionPolicy Bypass -File tools\install\install-windows.ps1
param(
    [string]$InstallDir = "$env:ProgramFiles\Quantis",
    [switch]$NoAutostart,
    [switch]$Quiet
)
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path  # repo root
$win  = Join-Path $root 'hosted\bin\windows-amd64'

Write-Host "== Quantis one-click installer (Windows) =="

# ensure the Windows binary exists
if (-not (Test-Path (Join-Path $win 'quantisd.exe'))) {
    Write-Host "building binaries..."
    & (Join-Path $root 'tools\build-hosted.ps1') | Out-Null
}

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

Write-Host "copying quantisd..."
Copy-Item (Join-Path $win 'quantisd.exe') (Join-Path $InstallDir 'quantisd.exe') -Force

Write-Host "building reconciler..."
Push-Location (Join-Path $root 'hosted\reconciler')
try {
    $env:GOOS = 'windows'; $env:GOARCH = 'amd64'
    & go build -o (Join-Path $InstallDir 'reconciler.exe') .
} finally {
    $env:GOOS = $null; $env:GOARCH = $null
    Pop-Location
}
if ($LASTEXITCODE -ne 0) { throw 'reconciler build failed' }

# autostart via scheduled task (runs at logon, no admin service needed)
if (-not $NoAutostart) {
    Write-Host "registering autostart (scheduled task 'Quantis')..."
    schtasks /Create /F /TN "Quantis" `
        /TR "`"$InstallDir\quantisd.exe`"" /SC ONLOGON /RL LIMITED | Out-Null
    schtasks /Create /F /TN "Quantis-Reconciler" `
        /TR "`"$InstallDir\reconciler.exe`"" /SC ONLOGON /RL LIMITED | Out-Null
}

# start services now (hidden, background)
Write-Host "starting services..."
Start-Process -FilePath (Join-Path $InstallDir 'quantisd.exe') -WindowStyle Hidden
Start-Process -FilePath (Join-Path $InstallDir 'reconciler.exe') -WindowStyle Hidden

if (-not $Quiet) {
    Start-Sleep 1
    Write-Host ""
    Write-Host "Quantis installed to: $InstallDir" -ForegroundColor Green
    Write-Host "  resources: http://localhost:9090/ns"
    Write-Host "  reconciler: http://localhost:9091/decl"
    if (-not $NoAutostart) { Write-Host "  autostart: ON (at next logon)" }
}
