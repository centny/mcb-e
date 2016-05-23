@echo off
set PATH=%PATH%;%MSYS_HOME%\usr\bin
call VsMsBuildCmd.bat
cd ..\cswf.ffcm
call pkg.bat
cd ..\cswf.doc
call pkg.bat
cd ..\mcb-e
bash -c "./build.sh  win.x64 ../cswf.doc/build/cswf.doc/"