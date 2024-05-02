go build -o out/
# run web server in backgroud
# Get the current date and time in the format you specified
current_date=$(date +"%Y_%m_%d_%H_%M_%S")
mkdir -p ./out/logs
# Run your web server in the background and redirect output to a file named with the current date and time
nohup ./out/web_server > ./out/logs/${current_date}_out.txt 2>&1 &

