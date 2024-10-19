@echo off
echo Condensing .go and .mod files into condensed.txt...
(for /r %%i in (*.go *.mod) do @type "%%i") > condensed.txt
echo Done! condensed.txt has been created.
pause