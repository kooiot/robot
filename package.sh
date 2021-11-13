# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

iot_tun_version=`./bin/robot-server --version`
echo "build version: $iot_tun_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        iot_tun_dir_name="iot_tun_${iot_tun_version}_${os}_${arch}"
        iot_tun_path="./packages/iot_tun_${iot_tun_version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./robot-client_${os}_${arch}.exe" ]; then
                continue
            fi
            if [ ! -f "./robot-server_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${iot_tun_path}
            mv ./robot-client_${os}_${arch}.exe ${iot_tun_path}/robot-client.exe
            mv ./robot-server_${os}_${arch}.exe ${iot_tun_path}/robot-server.exe
        else
            if [ ! -f "./robot-client_${os}_${arch}" ]; then
                continue
            fi
            if [ ! -f "./robot-server_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${iot_tun_path}
            mv ./robot-client_${os}_${arch} ${iot_tun_path}/robot-client
            mv ./robot-server_${os}_${arch} ${iot_tun_path}/robot-server
        fi
        cp ../LICENSE ${iot_tun_path}
        cp -rf ../conf/* ${iot_tun_path}

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${iot_tun_dir_name}.zip ${iot_tun_dir_name}
        else
            tar -zcf ${iot_tun_dir_name}.tar.gz ${iot_tun_dir_name}
        fi
        cd ..
        rm -rf ${iot_tun_path}
    done
done

cd -
