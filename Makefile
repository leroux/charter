# The charter's build. `make check` is everything CI runs on Go and the viewer.
.PHONY: check check-go claims viewer viewer-dev FORCE

NODE_MODULES := node_modules/.modules.yaml

$(NODE_MODULES): pnpm-lock.yaml package.json claims-viewer/package.json
	pnpm install --frozen-lockfile
	touch $@

check: check-go example viewer

# The gate parses TypeScript citations with node_modules/typescript, so the
# JavaScript workspace is installed first.
check-go: $(NODE_MODULES)
	go mod tidy -diff && go vet ./... && go test -count=1 ./...

# Regenerate the charter's own claims/specs.go, claims/index.ts and digest
# after editing an area document. Commit what it writes.
claims: $(NODE_MODULES)
	go run ./cmd/claimsgen

claims-viewer/public/charter.json: FORCE
	mkdir -p claims-viewer/public
	go run ./cmd/claimsmap -charter self -format json -o $@
	printf '%s\n' '[{"data": "charter.json", "label": "charter"}]' > claims-viewer/public/charters.json

viewer: claims-viewer/public/charter.json $(NODE_MODULES)
	cd claims-viewer && pnpm run typecheck && pnpm test && pnpm run build

viewer-dev: claims-viewer/public/charter.json $(NODE_MODULES)
	cd claims-viewer && pnpm run dev

# The example adopter, a module of its own: its tests and its gate, and its
# own claimsgen's lint. go test ./... runs the same through adopter_test.go.
.PHONY: example
example:
	cd example && go vet ./... && go test -count=1 ./... && go run ./cmd/claimsgen lint
