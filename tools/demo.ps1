# LIVE DEMO - Quantis hosted resource server on Windows
# Shows the ONE core idea: a resource is simultaneously
#   1) a file you can read        (like /proc)
#   2) an API endpoint you POST to       (like REST)
#   3) an event stream you can watch     (like Kafka/SSE)
#
# Usage from repo root:
#   powershell -File tools\demo.ps1
$ErrorActionPreference = 'Stop'

$bin  = Resolve-Path 'hosted\bin\windows-amd64\quantisd.exe'
$port = 9090
$base = "http://localhost:$port"

Write-Host ""
Write-Host "=========== QUANTIS USABILITY DEMO (Windows) ===========" -ForegroundColor Cyan
Write-Host "Starting quantisd on port $port ..." -ForegroundColor Gray
$job = Start-Job -ScriptBlock { param($b) & $b } -ArgumentList $bin
Start-Sleep 2

try {
    Write-Host ""
    Write-Host "[1] NAMESPACE - what does this system expose?" -ForegroundColor Yellow
    Write-Host "GET /ns  (introspection: everything is self-describing)" -ForegroundColor Gray
    $ns = cmd /c "curl.exe -s $base/ns"
    Write-Host $ns

    Write-Host ""
    Write-Host "[2] READ - every resource is a FILE you can cat." -ForegroundColor Yellow
    Write-Host "GET /r/sys/mem" -ForegroundColor Gray
    cmd /c "curl.exe -s $base/r/sys/mem"
    Write-Host "GET /r/srv/status" -ForegroundColor Gray
    cmd /c "curl.exe -s $base/r/srv/status"

    Write-Host ""
    Write-Host "[3] WRITE - every resource is an API you POST to." -ForegroundColor Yellow
    Write-Host "POST /r/ai/models/qm-demo/infer" -ForegroundColor Gray
    cmd /c "curl.exe -s -X POST -d ""hello quantis stream tokens"" $base/r/ai/models/qm-demo/infer"

    Write-Host ""
    Write-Host "[4] WATCH - every resource is an EVENT STREAM (SSE)." -ForegroundColor Yellow
    Write-Host "Watching /dev/net0 for 3 seconds - live frames, no polling:" -ForegroundColor Gray
    cmd /c "curl.exe -s -N --max-time 3 $base/r/dev/net0/watch"
    Write-Host ""
    Write-Host "(end of stream capture)" -ForegroundColor DarkGray

    Write-Host ""
    Write-Host "[5] MODEL AS A FILE - MedScan triage via the resource." -ForegroundColor Yellow
    Write-Host "POST /r/ai/models/medscan/infer (watch the triage stream live):" -ForegroundColor Gray
    $report = "Chest CT: 9mm spiculated nodule in the right upper lobe with a small pleural effusion. No acute consolidation."
    $w = Start-Job -ScriptBlock { param($u) cmd /c "curl.exe -s -N --max-time 4 $u" } -ArgumentList "$base/r/ai/models/medscan/infer/watch"
    Start-Sleep 1
    cmd /c "curl.exe -s -X POST -d ""$report"" $base/r/ai/models/medscan/infer"
    Wait-Job $w | Out-Null
    Receive-Job $w
    Remove-Job $w -Force -ErrorAction SilentlyContinue
    Write-Host ""
    Write-Host "GET /r/ai/models/medscan/infer (cached structured result):" -ForegroundColor Gray
    cmd /c "curl.exe -s $base/r/ai/models/medscan/infer"
} finally {
    Stop-Job $job -ErrorAction SilentlyContinue
    Remove-Job $job -ErrorAction SilentlyContinue
}

Write-Host ""
Write-Host "=========== WHAT THIS PROVES ===========" -ForegroundColor Green
Write-Host "1 FILE  : the same object that a shell would cat"
Write-Host "2 API   : the same object that a program POSTs to"
Write-Host "3 EVENTS: the same object that a dashboard subscribes to"
Write-Host ""
Write-Host "That is the ENTIRE Quantis thesis in 4 commands." -ForegroundColor Green
