# Create an issue if the installtion script does not works

current_version="v1.1.0"
root_path="$HOME/templatify"
install_path="$root_path/templatify"
github_download_url="https://github.com/scientific-dev/templatify/releases/download"

if [ ${#1} != 0 ]; then
    current_version=$1
fi

echo "Installing templatify version: $current_version"
echo "Installing at $install_path"

platform="$(uname -s | tr '[:upper:]' '[:lower:]')"
case $platform in 
    darwin*)    os="macOS";;
    linux*)     os="linux";; 
    msys_nt*)   os="win";;
    mingw*)     os="win";;
    cygwin_nt*) os="win";;
    *) echo "Unknown os for installtion. Please create an issue on github to make templatify available for $platform!" && exit 1 ;;
esac

case "`uname -m`" in
    x86)    arch="amd64";;
    ia64)   arch="amd64";;
    i?86)   arch="amd64";;
    amd64)  arch="amd64";;
    x86_64) arch="amd64";;
    *) echo "You architecture is not supported. Please create an issue on github to make templatify available for this arch! If you think this is a mistake you can directly download it from github releases!" && exit 1;;
esac

target="$github_download_url/$current_version/templatify-$os-$arch"

if [ $os == "win" ]; then 
    target="$target.exe"
    install_path="$install_path.exe"
fi

if [ ! -d $root_path ]; then 
    mkdir $root_path
fi

echo "Installing from $target for $os $arch"
echo ""
echo "Installing..."
curl -# --location -o $install_path -O $target

if [ $? -ne 0 ]; then
    echo 'An error has occurred while installing templatify!'
    exit 0
fi

echo ""
echo "Successfully installed templatify into $install_path. Now you can manually set it to the path and create templates ;)"
