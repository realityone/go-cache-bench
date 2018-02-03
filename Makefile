path=$(shell pwd)
export BENCHTIME=3s
export KEY_METHOD=random

define bench_cmd
	$(eval RESULT=$(path)/result/$@/)
	mkdir -p $(RESULT)
	@echo "Bench with $@" > $(RESULT)/output.txt
	@echo "BENCHTIME=$(BENCHTIME)" >> $(RESULT)/output.txt
	@echo "KEY_METHOD=$(KEY_METHOD)" >> $(RESULT)/output.txt
	go test -o $(RESULT)/$@.test -benchtime=$(BENCHTIME) -bench="$@" -outputdir=$(RESULT) -benchmem -cpuprofile=cpu.prof -memprofile=memory.prof -v ./cache_test/ >> $(RESULT)/output.txt
	go-torch $(RESULT)/$@.test $(RESULT)/cpu.prof -f $(RESULT)/cpu.svg
	go-torch $(RESULT)/$@.test $(RESULT)/memory.prof -f $(RESULT)/memory.svg
endef

gcache:
	$(call bench_cmd)

groupcache:
	$(call bench_cmd)

cache2go:
	$(call bench_cmd)

all: clean gcache groupcache cache2go

clean:
	rm -rf result