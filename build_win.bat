@echo off

:make_built_folder
echo Make built folder
mkdir built
set time_str=%date:/=%_%time%
set time_str=%time_str: =%
set time_str=%time_str::=%
set time_str=%time_str:.=%
set time_str=%time_str:周=%
set time_str=%time_str:星期=%
set time_str=%time_str:一=%
set time_str=%time_str:二=%
set time_str=%time_str:三=%
set time_str=%time_str:四=%
set time_str=%time_str:五=%
set time_str=%time_str:六=%
set time_str=%time_str:七=%
for /f %%i in (version) do set version_str=%%i
set built_folder=built\%version_str%_win_%time_str%
mkdir %built_folder%

:go_build
echo Run go build
go build -o %built_folder%\blog_shell.exe

:copy_files
echo Copy config.json
copy config.json %built_folder%\

echo Build done. Check %built_folder%.
