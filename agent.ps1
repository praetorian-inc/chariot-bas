[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$api = 'https://d0qcl2e18h.execute-api.us-east-2.amazonaws.com/chariot'
$asset = '<asset>'
$account = '<account>'
$dat = ""

function Execute {
    Param([String]$File)

    $result = New-Object PSObject -Property @{
        Code = 0
        Stdout = ""
    }

    try {
        $stdoutTempFile = New-TemporaryFile
        $stderrTempFile = New-TemporaryFile

        $proc = Start-Process -FilePath $File -NoNewWindow -PassThru -RedirectStandardOutput $stdoutTempFile -RedirectStandardError $stderrTempFile
        $proc | Wait-Process -Timeout 45 -ErrorAction SilentlyContinue -ErrorVariable timeoutVar

        if ($timeoutVar) {
            $proc | Stop-Process
            $result.Code = 52
            return $result
        }

        $result.Stdout = Get-Content -Path $stdoutTempFile -Raw
        $result.Code = if (Test-Path $File) {$proc.ExitCode} Else {52}
    } catch {
        $result.Code = 49
    }

    return $result
}

$stdout = ""
while ($true) {
    try {
        $headers = @{
            "dos" = "windows-$Env:PROCESSOR_ARCHITECTURE" 
            "dat" = $dat
            "account" = $account
            "asset" = $asset
        }

        $response = Invoke-WebRequest -Uri $api -Method POST -Headers $headers -Body $stdout -UseBasicParsing -MaximumRedirection 0 -ErrorAction SilentlyContinue

        $redirectUrl = $response.Headers['Location']

        if ($redirectUrl -match "([0-9a-fA-F]{8}[0-9a-fA-F]{4}[0-9a-fA-F]{4}[0-9a-fA-F]{4}[0-9a-fA-F]{12})") {
            $uuid = $matches[1]
        } else {
            $uuid = ""
        }

        if ($uuid) {
            $outFile = Join-Path -Path $env:TEMP -ChildPath $uuid
            Invoke-WebRequest -Uri $redirectUrl -OutFile $outFile -UseBasicParsing

            $result = Execute $outFile
            $dat = "${uuid}:$($result.Code)"

            $stdout = $result.Stdout
        } else {
            if($outFile -ne $null) {
                Remove-Item $outFile -ErrorAction SilentlyContinue
            }
            $dat = ""
            Start-Sleep -Seconds 60
        }
    } catch {
        $dat = ""
        Start-Sleep -Seconds 60
    }
}
