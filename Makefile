# The charter's build. `make check` before pushing runs all of it.
.PHONY: check check-go claims examples FORCE

NODE_MODULES := node_modules/.modules.yaml
CHARTER := build/charter

$(NODE_MODULES): pnpm-lock.yaml package.json
	pnpm install --frozen-lockfile
	touch $@

$(CHARTER): FORCE
	go build -o $@ ./cmd/charter

check: check-go examples

# The tooling's tests, the charter's own gate among them (cmd/charter's
# adopter tests run the examples' gates too), then the gate as an adopter
# runs it. The TypeScript reader's tests parse with the root's
# node_modules/typescript, so the workspace is installed first.
check-go: $(NODE_MODULES) $(CHARTER)
	go mod tidy -diff && go vet ./... && go test -count=1 ./...
	$(CHARTER) check

# Regenerate the charter's own claims/specs.go after editing an area
# document, and each example's. Commit what it writes.
claims: $(CHARTER)
	$(CHARTER) gen
	for d in examples/*/; do (cd $$d && ../../$(CHARTER) gen) || exit 1; done

# The examples, charters of their own: each one's gate and report.
examples: $(CHARTER)
	for d in examples/*/; do (cd $$d && ../../$(CHARTER) check && ../../$(CHARTER) report -md >/dev/null) || exit 1; done
