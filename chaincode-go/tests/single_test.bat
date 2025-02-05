
path=c:\Python311\;%PATH%

@REM create keys via python
cd ../dilithiumpy
python keys.py
@REM SIGN Message
python signmsg.py

cd ../tests

@REM Expected to success
go run v1.go

@REM Expected to fail
go run v2.go
