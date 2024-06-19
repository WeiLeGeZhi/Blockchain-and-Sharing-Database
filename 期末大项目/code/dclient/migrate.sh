rm -rf ./src/assets/contracts/*
cp ../dserver/build/contracts/* ./src/assets/contracts/

Remove-Item .\src\assets\contracts\* -Force
Copy-Item ..\dserver\build\contracts\* .\src\assets\contracts\
