# Uninstall Quantis from Windows (reverse of install-windows.ps1).
# Usage: powershell -ExecutionPolicy Bypass -File tools\install\uninstall-windows.ps1
param([string]$InstallDir = "$env:ProgramFiles\Quantis")
$ErrorActionPreference = 'SilentlyContinue'

Write-Host "== Quantis uninstaller (Windows) =="
Stop-Process -Name quantisd -ErrorAction SilentlyContinue
Stop-Process -Name reconciler -ErrorAction SilentlyContinue
schtasks /Delete /F /TN "Quantis"
schtasks /Delete /F /TN "Quantis-Reconciler"
Remove-Item -Recurse -Force $InstallDir
Write-Host "Removed $InstallDir"
Write-Host "Uninstall complete."
