param(
  [string]$BinDir = (Join-Path $HOME ".local\bin")
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$ApiBase = "https://api.github.com/repos/soapwong703/personal-gitignore/releases/latest"
$ReleaseBase = "https://github.com/soapwong703/personal-gitignore/releases/latest/download"

$release = Invoke-RestMethod -Headers @{ Accept = "application/vnd.github+json" } -Uri $ApiBase
$version = $release.tag_name
if (-not $version) {
  throw "Unable to determine the latest release version."
}

$arch = $env:PROCESSOR_ARCHITECTURE
if ($env:PROCESSOR_ARCHITEW6432) {
  $arch = $env:PROCESSOR_ARCHITEW6432
}

switch ($arch.ToUpperInvariant()) {
  "AMD64" { $arch = "amd64" }
  "ARM64" { $arch = "arm64" }
  default {
    throw "Unsupported Windows architecture: $arch"
  }
}

$asset = "pgi_windows_${arch}.zip"
$url = "$ReleaseBase/$asset"

$tmpRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("pgi-" + [guid]::NewGuid().ToString("N"))
$archive = Join-Path $tmpRoot $asset
$extractDir = Join-Path $tmpRoot "extract"

New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
New-Item -ItemType Directory -Force -Path $tmpRoot | Out-Null
New-Item -ItemType Directory -Force -Path $extractDir | Out-Null

$client = New-Object System.Net.WebClient
try {
  Write-Output "Downloading pgi $version for windows/$arch..."
  $client.DownloadFile($url, $archive)
  Expand-Archive -Path $archive -DestinationPath $extractDir -Force

  $packageDir = Join-Path $extractDir "pgi_windows_${arch}"
  $source = Join-Path $packageDir "pgi.exe"
  $destination = Join-Path $BinDir "pgi.exe"
  Copy-Item -Force -Path $source -Destination $destination

  Write-Output "Installed pgi $version to $destination"
}
finally {
  $client.Dispose()
  Remove-Item -Recurse -Force $tmpRoot -ErrorAction SilentlyContinue
}