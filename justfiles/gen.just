# generate README docs for commands
[group('gen')]
gen-docs:
    @go run scripts/cmds_docs.go
