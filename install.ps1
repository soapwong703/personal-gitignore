$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$ReleaseBase = "https://github.com/soapwong703/personal-gitignore/releases/latest/download"
$BinDir = Join-Path $HOME ".local\bin"

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

function Get-LatestVersion {
  param(
    [string]$DownloadUrl
  )

  $request = [System.Net.HttpWebRequest]::Create($DownloadUrl)
  $request.Method = "HEAD"
  $request.AllowAutoRedirect = $false

  try {
    $response = $request.GetResponse()
    $location = $response.Headers["Location"]
    $response.Dispose()
  }
  catch [System.Net.WebException] {
    $response = $_.Exception.Response
    if ($response) {
      $location = $response.Headers["Location"]
      $response.Dispose()
    }
    else {
      throw
    }
  }

  if (-not $location) {
    throw "Unable to determine the latest release version."
  }

  if ($location -match "/releases/download/([^/]+)/") {
    return $Matches[1]
  }

  throw "Unable to determine the latest release version."
}

$version = Get-LatestVersion -DownloadUrl $url

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

  $pathEntries = $env:Path -split ';'
  if ($pathEntries -notcontains $BinDir) {
    Write-Warning @"
${BinDir} is not on PATH.
Add it to PATH with:

  `$env:Path = `"$BinDir;$env:Path`"
"@
  }
}
finally {
  $client.Dispose()
  Remove-Item -Recurse -Force $tmpRoot -ErrorAction SilentlyContinue
}