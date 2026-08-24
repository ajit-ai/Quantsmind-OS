# Downloads the current 9front amd64 ISO into vm\9front.iso
# Source: iso.only9fans.com (official nightly ISO host linked from 9front.org/download)
$ErrorActionPreference = 'Stop'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$vmDir  = Join-Path $PSScriptRoot '..\vm'
New-Item -ItemType Directory -Force -Path $vmDir | Out-Null
$host_  = 'http://iso.only9fans.com'

function Get-WithRetry($url) {
    foreach ($i in 1..6) {
        try { return Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 60 }
        catch { Start-Sleep 3 }
    }
    throw "unreachable: $url"
}

$page = Get-WithRetry $host_
$link = ($page.Links | ForEach-Object { $_.href } |
    Where-Object { $_ -match '\.amd64\.iso\.gz$' } |
    Select-Object -Last 1)
if (-not $link) { throw 'no amd64 iso.gz found at iso host' }
if ($link -notmatch '^https?://') { $link = $host_ + $link }
Write-Host "Downloading: $link"

$tmp    = Join-Path $vmDir '9front.iso.gz'
$isoOut = Join-Path $vmDir '9front.iso'

# Resumable download: keeps appending across interruptions via HTTP Range.
$totalLen = -1
while ($true) {
    $have = if (Test-Path $tmp) { (Get-Item $tmp).Length } else { 0 }
    $req = [Net.HttpWebRequest]::Create($link)
    $req.AllowAutoRedirect = $true
    if ($have -gt 0) {
        $req.AddRange($have)
        Write-Host ("resuming at {0:N0} MB" -f ($have / 1MB))
    }
    try {
        $resp = $req.GetResponse()
        $totalLen = $resp.ContentLength + $have
        $stream = $resp.GetResponseStream()
        $out = [IO.File]::Open($tmp, [IO.FileMode]::Append, [IO.FileAccess]::Write)
        $buf = New-Object byte[] (4 * 1024 * 1024)
        try {
            while (($n = $stream.Read($buf, 0, $buf.Length)) -gt 0) { $out.Write($buf, 0, $n) }
        } finally { $out.Dispose(); $stream.Dispose(); $resp.Close() }
    } catch [Net.WebException] {
        Write-Host "interrupted: $_ ; retrying in 5s"
        Start-Sleep 5
        continue
    }
    break
}
"Downloaded {0:N0} MB" -f ((Get-Item $tmp).Length / 1MB)

Write-Host "Decompressing -> $isoOut"
$input  = [IO.File]::OpenRead($tmp)
$output = [IO.File]::Create($isoOut)
try {
    (New-Object IO.Compression.GzipStream($input,[IO.Compression.CompressionMode]::Decompress)).CopyTo($output)
} finally { $output.Dispose(); $input.Dispose() }
Remove-Item $tmp
"Done: $isoOut ({0:N0} MB)" -f ((Get-Item $isoOut).Length / 1MB)
