#!/usr/bin/env bash

git add model
git commit -m "dumper: $(date '+%d-%m-%Y_%T')"
git push
