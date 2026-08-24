# Boot Quantis / 9front in QEMU on Windows.
# Usage:
#   .\tools\run-qemu.ps1 -Install            # first run: creates disk, boots ISO
#   .\tools\run-qemu.ps1                     # daily: boots installed disk
# Requires: qemu-system-x86_64 on PATH or in "C:\Program Files\qemu"
param(
    [switch]$Install,
    [string]$Disk = "$PSScriptRoot\..\vm\quantis.qcow2",
    [string]$Iso  = "$PSScriptRoot\..\vm\9front.iso",
    [int]$MemoryMB = 2048,
    [int]$Cores = 4
)
$ErrorActionPreference = 'Stop'

$qemu = Get-Command qemu-system-x86_64 -ErrorAction SilentlyContinue
if (-not $qemu) {
    $candidate = 'C:\Program Files\qemu\qemu-system-x86_64.exe'
    if (Test-Path $candidate) { $qemu = @{ Source = $candidate } }
    else { throw 'QEMU not found. Install with: winget install SoftwareFreedomConservancy.QEMU' }
}

New-Item -ItemType Directory -Force -Path (Split-Path $Disk) | Out-Null

# Acceleration: WHPX (Windows Hypervisor Platform) if enabled, else TCG software emulation.
$accel = 'tcg,thread=multi'
$whpx = Get-WindowsOptionalFeature -Online -FeatureName HypervisorPlatform -ErrorAction SilentlyContinue
if ($whpx -and $whpx.State -eq 'Enabled') { $accel = 'whpx' }

if ($Install) {
    if (-not (Test-Path $Iso)) { throw "ISO missing at $Iso. Run tools\get-iso.ps1 first." }
    if (-not (Test-Path $Disk)) {
        & $qemu.Source img create -f qcow2 $Disk 20G
    }
    & $qemu.Source `
        -machine q35,accel=$accel -cpu max -smp $Cores -m $MemoryMB `
        -drive file=$Disk,if=nvme,format=qcow2 `
        -cdrom $Iso -boot d `
        -netdev user,id=n0 -device rtl8139,netdev=n0 `
        -device VGA
} else {
    if (-not (Test-Path $Disk)) { throw "No disk at $Disk. Run with -Install first." }
    & $qemu.Source `
        -machine q35,accel=$accel -cpu max -smp $Cores -m $MemoryMB `
        -drive file=$Disk,if=nvme,format=qcow2 `
        -netdev user,id=n0 -device rtl8139,netdev=n0 `
        -device VGA
}
