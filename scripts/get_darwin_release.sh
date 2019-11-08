#!/bin/bash

wget https://github.com/dlampsi/cataloger/releases/download/$1/cataloger_$1_darwin_amd64.zip \
&& unzip cataloger_$1_darwin_amd64.zip \
&& cd cataloger_$1_darwin_amd64 \
&& mv cataloger_darwin_amd64 /usr/local/bin/cataloger \
&& chmod +x /usr/local/bin/cataloger \
&& cd .. && rm -rf cataloger_$1_darwin_amd64*