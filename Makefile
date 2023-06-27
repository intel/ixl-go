export GOAMD64=v3

.PHONY: lint docs test  fuzz

lint:
	golangci-lint run  ./...
test:
	go test -count=1 -timeout 30s  -v ./... 
docs:
	gomarkdoc ./filter -o ./filter/doc.md
	gomarkdoc ./crc -o ./crc/doc.md
	gomarkdoc ./compress -o ./compress/doc.md
	gomarkdoc ./datamove -o ./datamove/doc.md
	gomarkdoc ./util/mem -o ./util/mem/doc.md

fuzz: compress_fuzz datamove_fuzz crc_fuzz filter_fuzz

FUZZ_WORKERS=8
FUZZ_TIME=10s

COMPRESS_FUZZ=FuzzDeflate FuzzWriter_Write FuzzInflate

compress_fuzz:
	$(foreach fuzz,$(COMPRESS_FUZZ),go test -parallel $(FUZZ_WORKERS) -json -run=^$$ -fuzz=$(fuzz) -fuzztime=$(FUZZ_TIME) ./compress ;)


DATAMOVE_FUZZ=FuzzCopy

datamove_fuzz:
	$(foreach fuzz,$(DATAMOVE_FUZZ),go test -parallel $(FUZZ_WORKERS) -json -run=^$$ -fuzz=$(fuzz) -fuzztime=$(FUZZ_TIME) ./datamove ;)


CRC_FUZZ=FuzzCRC64 FuzzCRC32

crc_fuzz:
	$(foreach fuzz,$(CRC_FUZZ),go test -parallel $(FUZZ_WORKERS) -json -run=^$$ -fuzz=$(fuzz) -fuzztime=$(FUZZ_TIME) ./crc ;)

FILTER_FUZZ=FuzzScan
filter_fuzz:
	$(foreach fuzz,$(FILTER_FUZZ),go test -parallel $(FUZZ_WORKERS) -json -run=^$$ -fuzz=$(fuzz) -fuzztime=$(FUZZ_TIME) ./filter ;)
