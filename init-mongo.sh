#!/bin/bash
set -e

mongoimport --username root --password password123 --authenticationDatabase admin --db oxford-5000 --collection words --type json --file /data/oxford-5000.json --jsonArray