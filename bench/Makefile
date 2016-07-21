results.bin: targets.txt
	@echo "Running attack..."
	@vegeta attack -duration 1m -targets=targets.txt > results.bin

plot.html: results.bin
	@echo "Running report..."
	@cat results.bin | vegeta report -reporter=plot > plot.html

clean:
	@rm -f plot.html results.bin

bench: clean plot.html
