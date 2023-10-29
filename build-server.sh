###
 # @Date: 2023-10-29 09:12:15
 # @LastEditors: hxlh
 # @LastEditTime: 2023-10-29 09:59:25
 # @FilePath: /1024/build-server.sh
### 
mkdir build
cd server/src
go build -o server
mv server ../../build/server
cp server.yml ../../build/server.yml