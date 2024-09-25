go build .
if exist "../../Desktop/developer-tools.exe" del "../../Desktop/developer-tools.exe"
move ./developer-tools.exe "../../Desktop/developer-tools.exe"