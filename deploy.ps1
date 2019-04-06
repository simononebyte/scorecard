
$source = @(
    "scorecard.exe",
    "scorecard.json"
)

$dest = "$($ENV:USERPROFILE)\Autotask Workplace\Management\Traction\Metrics\"

foreach ($f in $source) {
    Write-Host "Deploying $f - " -NoNewline
    Copy-Item -Path ".\$f" -Destination $dest -Force
    if ($?) {
        Write-Host "OK" -ForegroundColor Green
    }

}