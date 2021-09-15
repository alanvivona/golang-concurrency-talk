step-0:
	git checkout -f main

step-1:
	git checkout -f 1-conurrent-v2

step-2:
	git checkout -f 1-channels-v1

local:
	./run

site:
	cd test-site && ./run