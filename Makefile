path=$(shell pwd)
export BENCHTIME=3s
export KEY_METHOD=random

gcache:
	mkdir -p result/$@
	go test -benchtime=$(BENCHTIME) -bench="$@" -outputdir=$(path)/result/$@ -benchmem -cpuprofile=cpu.prof -memprofile=memory.prof -v ./cache_test/
	mv cache_test.test $(path)/result/$@

groupcache:
	mkdir -p result/$@
	go test -benchtime=$(BENCHTIME) -bench="$@" -outputdir=$(path)/result/$@ -benchmem -cpuprofile=cpu.prof -memprofile=memory.prof -v ./cache_test/
	mv cache_test.test $(path)/result/$@

cache2go:
	mkdir -p result/$@
	go test -benchtime=$(BENCHTIME) -bench="$@" -outputdir=$(path)/result/$@ -benchmem -cpuprofile=cpu.prof -memprofile=memory.prof -v ./cache_test/
	mv cache_test.test $(path)/result/$@

all: gcache groupcache cache2go

clean:
	rm -rf result cache_test.test