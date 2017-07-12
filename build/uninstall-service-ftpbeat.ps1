# delete service if it exists
if (Get-Service ftpbeat -ErrorAction SilentlyContinue) {
  $service = Get-WmiObject -Class Win32_Service -Filter "name='ftpbeat'"
  $service.delete()
}
