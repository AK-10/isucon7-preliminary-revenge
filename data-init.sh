cd ~/isubata/bench
./bin/gen-initial-dataset

zcat ~/isubata/bench/isucon7q-initial-dataset.sql.gz | sudo mysql isubata
