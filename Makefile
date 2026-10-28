.PHONY: build build-frontend test vet race clean dev

build-frontend:
	cd frontend && pnpm install --no-frozen-lockfile && pnpm build
	@mkdir -p httpapi/web web
	@cp -r frontend/dist/* httpapi/web/
	@cp -r frontend/dist/* web/
	@rm -f httpapi/web/app.js httpapi/web/style.css web/app.js web/style.css

build: build-frontend
	go build ./...

test:
	go test ./... -count=1

vet:
	go vet ./...

race:
	go test -race ./...

clean:
	rm -f goremote goremote_ui goremote_new2
	rm -rf frontend/dist frontend/node_modules/.cache
