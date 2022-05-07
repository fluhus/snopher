:: Builds an example's go DLL and runs its python script.
:: Use from within the src directory: run.bat [DIR_NAME]
go build -o %1.dll -buildmode=c-shared ./%1 && ^
python %1\%1.py
