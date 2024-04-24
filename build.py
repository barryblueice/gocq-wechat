import subprocess

target_file = "main.go recieve.go send.go WebsocketReverse.go"

for arch in ["windows/amd64", "windows/386", "windows/arm64", "windows/arm"]:
    os_name, arch_name = arch.split("/")

    ps_command = [
        'powershell',
        '-Command',
        f'$env:CGO_ENABLED="0"; $env:GOOS="{os_name}"; $env:GOARCH="{arch_name}"; go build -o .\\release\\gocq-wechat-{os_name}-{arch_name}.exe {target_file}'
    ]

    # 使用 subprocess 运行命令
    subprocess.run(ps_command, shell=True)
