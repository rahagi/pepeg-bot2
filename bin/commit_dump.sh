#!/usr/bin/env bash

git pull
git add model
git commit -m "dumper: $(date '+%d-%m-%Y_%T')"
git push
