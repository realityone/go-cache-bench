path=$(shell pwd)
export BENCHTIME=3s
export KEY_METHOD=random

define bench_cmd
    mkdir -p result/$@
	go test -o result/$@/$@.test -benchtime=$(BENCHTIME) -bench="$@" -outputdir=$(path)/result/$@ -benchmem -cpuprofile=cpu.prof -memprofile=memory.prof -v ./cache_test/
	go-torch $(path)/result/$@/$@.test $(path)/result/$@/cpu.prof -f $(path)/result/$@/torch.svg
endef

gcache:
	$(call bench_cmd)

groupcache:
	$(call bench_cmd)

cache2go:
	$(call bench_cmd)

all: gcache groupcache cache2go

clean:
	rm -rf result