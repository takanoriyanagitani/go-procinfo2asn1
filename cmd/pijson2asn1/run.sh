#!/bin/sh

jq -c -n '[
	{pid: 42, rss: 255, cpu: 99, mem: 30, state: "", time: ""},
	{pid: 43, rss: 155, cpu: 80, mem: 20, state: "", time: ""}
]' |
	jq -c '.[]' |
	./pijson2asn1 |
	xxd -ps |
	tr -d '\n' |
	python3 \
		-m asn1tools \
		convert \
		-i der \
		-o xer \
		./procinfo.asn \
		SpiArray \
		- |
	bat --language=xml
