#!/bin/bash

psql -U boot -d olx_clone -f /home/mr_ali/golang/project/olx_clone/migrations/update_product_status.sql
