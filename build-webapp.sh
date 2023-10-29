###
 # @Date: 2023-10-29 08:15:06
 # @LastEditors: hxlh
 # @LastEditTime: 2023-10-29 09:58:52
 # @FilePath: /1024/build-webapp.sh
### 
rm build/static/*
cd webapp
npm run build
mv dist/* ../build/static/