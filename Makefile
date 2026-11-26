# The charter's build. `make check` is everything CI runs on Go and the viewer.
.PHONY: check check-go claims examples viewer viewer-dev FORCE

NODE_MODULES := node_modules/.modules.yaml
CHARTER := build/charter

$(NODE_MODULES): pnpm-lock.yaml package.json claims-viewer/package.json
	pnpm install --frozen-lockfile
	touch $@

$(CHARTER): FORCE
	go build -o $@ ./cmd/charter

check: check-go examples viewer

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
	cd example && ../$(CHARTER) gen
	cd example-polyglot && ../$(CHARTER) gen

# The examples, charters of their own: each one's gate, lint and report.
examples: $(CHARTER)
	cd example && ../$(CHARTER) check && ../$(CHARTER) lint && ../$(CHARTER) report -md >/dev/null
	cd example-polyglot && ../$(CHARTER) check && ../$(CHARTER) lint && ../$(CHARTER) report -md >/dev/null

claims-viewer/public/charter.json: $(CHARTER) FORCE
	mkdir -p claims-viewer/public
	$(CHARTER) export -o $@
	printf '%s\n' '[{"data": "charter.json", "label": "charter"}]' > claims-viewer/public/charters.json

viewer: claims-viewer/public/charter.json $(NODE_MODULES)
	cd claims-viewer && pnpm run typecheck && pnpm test && pnpm run build

viewer-dev: claims-viewer/public/charter.json $(NODE_MODULES)
	cd claims-viewer && pnpm run dev
