# delete service if it already exists
if (Get-Service ftpbeat -ErrorAction SilentlyContinue) {
  $service = Get-WmiObject -Class Win32_Service -Filter "name='ftpbeat'"
  $service.StopService()
  Start-Sleep -s 1
  $service.delete()
}

$workdir = Split-Path $MyInvocation.MyCommand.Path

# create new service
New-Service -name ftpbeat `
  -displayName ftpbeat `
  -binaryPathName "`"$workdir\\ftpbeat.exe`" -c `"$workdir\\ftpbeat.yml`" -path.home `"$workdir`" -path.data `"C:\\ProgramData\\ftpbeat`""
