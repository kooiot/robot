# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

robot_version=`./bin/robot-server --version`
echo "build version: $robot_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        robot_dir_name="robot_${robot_version}_${os}_${arch}"
        robot_path="./packages/robot_${robot_version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./robot-client_${os}_${arch}.exe" ]; then
                continue
            fi
            if [ ! -f "./robot-server_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${robot_path}
            mv ./robot-client_${os}_${arch}.exe ${robot_path}/robot-client.exe
            mv ./robot-server_${os}_${arch}.exe ${robot_path}/robot-server.exe
        else
            if [ ! -f "./robot-client_${os}_${arch}" ]; then
                continue
            fi
            if [ ! -f "./robot-server_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${robot_path}
            mv ./robot-client_${os}_${arch} ${robot_path}/robot-client
            mv ./robot-server_${os}_${arch} ${robot_path}/robot-server
        fi
        cp ../LICENSE ${robot_path}
        cp -rf ../conf/* ${robot_path}

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${robot_dir_name}.zip ${robot_dir_name}
        else
            tar -zcf ${robot_dir_name}.tar.gz ${robot_dir_name}
        fi
        cd ..
        rm -rf ${robot_path}
    done
done

cd -
