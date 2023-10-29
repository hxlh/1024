###
 # @Date: 2023-10-29 09:12:15
 # @LastEditors: hxlh
 # @LastEditTime: 2023-10-29 09:12:21
 # @FilePath: /1024/build-server.sh
### 
rm -R build
mkdir build
mkdir build/static
cd server/src
go build -o server
mv server ../../build/server
cp server.yml ../../build/server.yml