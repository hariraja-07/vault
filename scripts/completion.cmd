@echo off
rem vault CMD Completion
rem Install: Add this line to your autoexec.nt or startup
rem prompt $P$G

doskey vault=python %USERPROFILE%\vault.py $*

:completion
if "%1"=="" goto end
if "%1"=="set" goto set_key
if "%1"=="get" goto get_key
if "%1"=="remove" goto remove_key
if "%1"=="list" goto list_opts
if "%1"=="help" goto help_opts
goto end

:set_key
echo set get remove list help completion --force -F --full -f
goto end

:get_key
echo set get remove list help completion --force -F --full -f
goto end

:remove_key
echo set get remove list help completion --force -F --full -f
goto end

:list_opts
echo --full -f set get remove list help completion --force -F
goto end

:help_opts
echo set get remove list help completion --force -F --full -f
goto end

:end
