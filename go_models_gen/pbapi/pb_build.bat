@echo off

D:\Qt\Tools\mingw810_64\bin\protoc.exe --go_out=. --plugin=protoc-gen-go=D:\go_work\gopath\bin\protoc-gen-go.exe *.proto
rem D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\protoc.exe -I ./ --go_out=plugins=grpc:. --plugin=protoc-gen-go=D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\protoc-gen-go.exe ui_api.proto
rem D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\protoc.exe -I ./ --go_out=plugins=grpc:. --plugin=protoc-gen-go=D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\protoc-gen-go.exe svr_api.proto


rem 生成`msg.pb` 消息类
D:\Qt\Tools\mingw810_64\bin\protoc.exe -I ./ --cpp_out=..\..\ *.proto
rem D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\protoc.exe --grpc_out=D:\qt_work\TestRpcU --plugin=protoc-gen-grpc=D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\grpc_cpp_plugin.exe ui_api.proto
rem D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\protoc.exe -I ./ --cpp_out=D:\qt_work\TestRpcU\ ui_api.proto
rem D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\protoc.exe --grpc_out=D:\qt_work\TestRpcU --plugin=protoc-gen-grpc=D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\grpc_cpp_plugin.exe svr_api.proto
rem D:\qt_work\Qt5.12.10\Tools\mingw730_64\bin\protoc.exe -I ./ --cpp_out=D:\qt_work\TestRpcU\ svr_api.proto
pause
