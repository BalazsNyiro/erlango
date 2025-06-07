#!/usr/bin/env bash

DIR_START=$(pwd)
pwd

rm .coverage

python3 -m coverage run -m unittest *test.py
python3 -m coverage report
python3 -m coverage html

cd $DIR_START


