CC := go build

.PHONY: hue-cli
hue-cli: hue-cli.go
	$(CC) -o $@

clean:
	$(RM) hue-cli
