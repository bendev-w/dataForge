# Install DataForge on Windows:
#   irm https://raw.githubusercontent.com/bendev-w/dataForge/main/install.ps1 | iex
# Optional env: $env:VERSION = "v1.0.0"; $env:INSTALL_DIR = "C:\tools\dataforge"
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$repo = "bendev-w/dataForge"
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }

$version = $env:VERSION
if (-not $version) {
    $version = (Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest").tag_name
}

$archive = "dataforge_$($version.TrimStart('v'))_windows_$arch.zip"
$base = "https://github.com/$repo/releases/download/$version"
$tmp = Join-Path ([IO.Path]::GetTempPath()) ([Guid]::NewGuid())
New-Item -ItemType Directory -Path $tmp | Out-Null

try {
    Write-Host "Downloading $archive ..."
    Invoke-WebRequest "$base/$archive" -OutFile "$tmp\$archive"
    Invoke-WebRequest "$base/checksums.txt" -OutFile "$tmp\checksums.txt"

    $line = Get-Content "$tmp\checksums.txt" | Where-Object { $_ -match " $([regex]::Escape($archive))$" }
    $expected = ($line -split "\s+")[0]
    $actual = (Get-FileHash "$tmp\$archive" -Algorithm SHA256).Hash
    if (-not $expected -or $expected -ne $actual.ToLower()) { throw "Checksum verification failed" }

    Expand-Archive "$tmp\$archive" -DestinationPath $tmp -Force

    $dest = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\dataforge" }
    New-Item -ItemType Directory -Path $dest -Force | Out-Null
    Copy-Item "$tmp\dataforge.exe" "$dest\dataforge.exe" -Force

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (($userPath -split ";") -notcontains $dest) {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$dest", "User")
        Write-Host "Added $dest to your PATH (restart your terminal)."
    }
    Write-Host "Installed dataforge $version to $dest"
} finally {
    Remove-Item $tmp -Recurse -Force -ErrorAction SilentlyContinue
}
