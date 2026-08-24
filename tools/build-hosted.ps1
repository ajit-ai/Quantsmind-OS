# Builds hosted/quantisd for every supported platform -> hosted/bin/<os>-<arch>/
# Usage: powershell -NoProfile -ExecutionPolicy Bypass -File tools\build-hosted.ps1
$ErrorActionPreference = 'Stop'
$here = Split-Path $PSScriptRoot -Parent
$src  = Join-Path $here 'hosted'
$out  = Join-Path $src 'bin'

$targets = @(
    @{ os = 'windows'; arch = 'amd64'; ext = '.exe' },
    @{ os = 'linux';   arch = 'amd64'; ext = '' },
    @{ os = 'linux';   arch = 'arm64'; ext = '' },
    @{ os = 'darwin';  arch = 'amd64'; ext = '' },
    @{ os = 'darwin';  arch = 'arm64'; ext = '' },
    @{ os = 'freebsd'; arch = 'amd64'; ext = '' },
    @{ os = 'openbsd'; arch = 'amd64'; ext = '' },
    @{ os = 'netbsd';  arch = 'amd64'; ext = '' }
)

Push-Location $src
try {
    $savedGOOS = $env:GOOS; $savedGOARCH = $env:GOARCH
    foreach ($t in $targets) {
        $env:GOOS = $t.os; $env:GOARCH = $t.arch
        $dir = Join-Path $out ("{0}-{1}" -f $t.os, $t.arch)
        New-Item -ItemType Directory -Force -Path $dir | Out-Null
        $bin = Join-Path $dir ("quantisd" + $t.ext)
        go build -trimpath -ldflags "-s -w" -o $bin .
        if ($LASTEXITCODE -ne 0) { throw "build failed: $($t.os)/$($t.arch)" }
        Write-Host ("OK  {0,-14} -> {1}" -f "$($t.os)/$($t.arch)", $bin)
    }
    $env:GOOS = $savedGOOS; $env:GOARCH = $savedGOARCH
} finally { Pop-Location }
"All platform binaries built."
