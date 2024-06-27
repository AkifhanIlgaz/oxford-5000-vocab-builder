#!/bin/bash
set -e

mongoimport --db oxford-5000 --collection words --type json --file /data/oxford-5000.json --jsonArray