cd ~/isubata/webapp/go
make
sudo systemctl restart isubata.golang.service
cd ~/isubata/bench
./bin/bench -remotes=127.0.0.1 -output result.json
jq < result.json
cd ~/isubata

