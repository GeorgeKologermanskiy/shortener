#!/usr/bin/bash

# copy nginx script to available
cp nginx.conf /etc/nginx/sites-available/nginx.conf

# create symlink
ls -s /etc/nginx/sites-available/nginx.conf /etc/nginx/sites-enabled/nginx.conf

systemctl restart nginx.service

systemctl status nginx.service