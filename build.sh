# Create an issue for any kind of build issues.

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/386" "linux/amd64" "linux/arm64")
suffix=("win64" "win32" "macOS64" "linux32" "linux64" "linux-arm64")
index=0

for platform in "${platforms[@]}"
do 

    output_name="templatify-${suffix[$index]}"
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    if [ $GOOS == "windows" ]; then 
        output_name="$output_name.exe"
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o build/$output_name *.go

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 0
    fi

    echo "Created $output_name"
    index=$(( index+1 ))

done
