for ( $i = 0; $i -lt $args.count; $i++ ) {
    Switch ($args[$i]) {
        "fmt" {
            go build -o ./build/fmt.exe ./fmt/main.go
            write-host "building $($args[$i]) finished with code $LASTEXITCODE"
        }
        "install" {
            Copy-Item -Path ".\build\*.exe" -Destination "C:\Tools" -Force -Verbose
        }
    } 
}