go build -o out/
# run web server in background
nohup ./out/web_server > ./out/out.text  2>&1  & 

