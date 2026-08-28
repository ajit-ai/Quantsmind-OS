# Assembles release artifacts for all platforms into dist/<os>/
# Produces: Windows (folder + install script), Linux (+ .deb tree),
# macOS, FreeBSD/OpenBSD/NetBSD as .tar.gz archives.
# Usage: powershell -File tools\package.ps1
$ErrorActionPreference = 'Stop'
$root = Split-Path $PSScriptRoot -Parent

# 1) build all quantisd binaries
& (Join-Path $root 'tools\build-hosted.ps1') | ForEach-Object { Write-Host $_ }

# 2) working area + dist
$stage = Join-Path $env:TEMP 'quantis-pkg'
$dist  = Join-Path $root 'dist'
Remove-Item $stage -Recurse -Force -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $stage | Out-Null
New-Item -ItemType Directory -Force -Path $dist  | Out-Null

# each target: quantisd + reconciler + installer
$targets = @(
    @{ os='linux';   arch='amd64'; bindir='linux-amd64';    ext='' },
    @{ os='linux';   arch='arm64'; bindir='linux-arm64';    ext='' },
    @{ os='darwin';  arch='amd64'; bindir='darwin-amd64';   ext='' },
    @{ os='darwin';  arch='arm64'; bindir='darwin-arm64';   ext='' },
    @{ os='freebsd'; arch='amd64'; bindir='freebsd-amd64';  ext='' },
    @{ os='openbsd'; arch='amd64'; bindir='openbsd-amd64';  ext='' },
    @{ os='netbsd';  arch='amd64'; bindir='netbsd-amd64';   ext='' }
)

# build reconciler for each target (build from inside its own module dir)
function Build-Reconciler($os, $arch, $destDir) {
    $env:GOOS = $os; $env:GOARCH = $arch
    Push-Location (Join-Path $root 'hosted\reconciler')
    try { & go build -trimpath -o (Join-Path $destDir 'reconciler') . }
    finally {
        Pop-Location
        $env:GOOS = $null; $env:GOARCH = $null
    }
    if ($LASTEXITCODE -ne 0) { throw "reconciler $os/$arch failed" }
}

foreach ($t in $targets) {
    $d = Join-Path $stage $t.bindir
    New-Item -ItemType Directory -Force -Path $d | Out-Null
    Copy-Item (Join-Path $root ("hosted\bin\" + $t.bindir + '\quantisd' + $t.ext)) (Join-Path $d 'quantisd') -Force
    Build-Reconciler $t.os $t.arch $d
    Copy-Item (Join-Path $root 'tools\install\README-DIST.md') (Join-Path $d 'README.md') -Force
    New-Item -ItemType Directory -Force -Path (Join-Path $stage ($t.bindir + '\tools')) | Out-Null
    Copy-Item (Join-Path $root 'tools\demo.ps1') (Join-Path $stage ($t.bindir + '\tools\demo.ps1')) -Force
}

# Windows layout (folder + installer, no gz)
$winDir = Join-Path $stage 'windows-amd64'
New-Item -ItemType Directory -Force -Path $winDir | Out-Null
Copy-Item (Join-Path $root 'hosted\bin\windows-amd64\quantisd.exe') (Join-Path $winDir 'quantisd.exe') -Force
Build-Reconciler 'windows' 'amd64' $winDir   # produces reconciler (no .exe); rename
Rename-Item (Join-Path $winDir 'reconciler') (Join-Path $winDir 'reconciler.exe')
Copy-Item (Join-Path $root 'tools\install\install-windows.ps1') (Join-Path $winDir 'install-windows.ps1') -Force

# copy to dist
Copy-Item $stage\* $dist -Recurse -Force

# tar.gz each non-windows target (userland tar.exe on Win10+)
foreach ($t in $targets) {
    $src = Join-Path $dist $t.bindir
    $gz  = Join-Path $dist ("quantis-" + $t.os + "-" + $t.arch + ".tar.gz")
    Push-Location $dist
    try { & tar -czf $gz $t.bindir } finally { Pop-Location }
}

""
"== PACKAGED into: $dist =="
Get-ChildItem $dist -Recurse -File | ForEach-Object { $_.FullName.Replace($dist + '\', '') }
