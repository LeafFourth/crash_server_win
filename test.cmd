echo off
echo hello
set i=0

:LOOP
    set /a i=%i% + 1
    echo %i%

IF %i% lss 10000 goto LOOP