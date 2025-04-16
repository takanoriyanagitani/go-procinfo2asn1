#!/bin/sh

printf \
	'%s\n' \
	'pid:42	rss:255	cpu:30.2	mem:50.4	state:S	time:00:00:01' \
	'pid:42	rss:255	cpu:31.5	mem:51.6	state:S	time:00:00:01' |
	./piltsv2asn1 |
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
