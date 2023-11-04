###
 # @Date: 2023-10-29 08:15:06
 # @LastEditors: hxlh
 # @LastEditTime: 2023-10-29 09:13:11
 # @FilePath: /1024/build-all.sh
### 
rm -R build
mkdir build
mkdir build/static
cd server/src
go build -o server
mv server ../../build/server
cp server.yml ../../build/server.yml
cd ../../
cd webapp
npm run build
mv dist/* ../build/static/