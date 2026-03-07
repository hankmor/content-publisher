@echo off

REM 启动脚本，用于在 Windows 系统上启动微信发布工具

echo 启动微信发布工具...

REM 检查可执行文件是否存在
if exist .\wechat-publisher.exe set EXECUTABLE=wechat-publisher.exe
if exist .\wechat-publisher-windows-amd64.exe set EXECUTABLE=wechat-publisher-windows-amd64.exe
if exist .\output\wechat-publisher-windows-amd64.exe set EXECUTABLE=output\wechat-publisher-windows-amd64.exe

if not defined EXECUTABLE (
    echo 错误：未找到可执行文件！
    echo 请先运行 build.sh 脚本构建应用。
    pause
    exit /b 1
)

REM 检查配置文件是否存在
if not exist .\configs\config.yaml (
    echo 警告：配置文件不存在！
    echo 请根据 configs\config.yaml.example 创建 configs\config.yaml 文件并填写相关配置。
    if exist .\configs\config.yaml.example (
        echo 正在复制配置文件模板...
        copy .\configs\config.yaml.example .\configs\config.yaml
        echo 请编辑 configs\config.yaml 文件，填写微信公众号的 AppID、AppSecret 和 Token。
    )
)

REM 启动应用
echo 正在启动应用...
%EXECUTABLE%
pause
