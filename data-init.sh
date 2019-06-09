cd ~/isubata/bench
./bin/gen-initial-dataset

cd ~/isubata
sudo ./db/init.sh

zcat ~/isubata/bench/isucon7q-initial-dataset.sql.gz | sudo mysql isubata

