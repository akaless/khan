# ═══════════════════════════════════════════════════════════
# خان (Khan) — Makefile
# چت سازمانی سبک فارسی
# ═══════════════════════════════════════════════════════════

# ─── نسخه ───
VERSION ?= 1.0.1
LDFLAGS := -s -w -X main.version=$(VERSION)

# ─── مسیرها ───
BIN_DIR := dist
SRC := ./cmd/server

# ─── پلتفرم‌ها ───
PLATFORMS := linux/amd64 windows/amd64 darwin/amd64 darwin/arm64

.PHONY: all build build-all test vet fmt lint clean run windows linux mac docker release check

## 🏗 همه بیلدها
all: build-all

## 🔨 بیلد برای پلتفرم فعلی
build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/khan-$(shell go env GOOS) $(SRC)
	@echo "✅ built: $(BIN_DIR)/khan-$(shell go env GOOS)"

## 🌍 بیلد برای همه پلتفرم‌ها
build-all:
	@mkdir -p $(BIN_DIR)
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		name="khan"; \
		if [ "$$os" = "windows" ]; then name="khan.exe"; fi; \
		if [ "$$os" = "darwin" ] && [ "$$arch" = "amd64" ]; then name="khan-macos-intel"; fi; \
		if [ "$$os" = "darwin" ] && [ "$$arch" = "arm64" ]; then name="khan-macos-arm64"; fi; \
		if [ "$$os" = "linux" ]; then name="khan-linux"; fi; \
		echo "→ building $$os/$$arch → $(BIN_DIR)/$$name"; \
		cd code && GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ../$(BIN_DIR)/$$name ./cmd/server || exit 1; \
		cd ..; \
	done
	@echo "✅ all builds complete"
	@ls -la $(BIN_DIR)/

## 🖥 بیلد ویندوز
windows:
	@mkdir -p $(BIN_DIR)
	cd code && GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ../$(BIN_DIR)/khan.exe ./cmd/server
	@echo "✅ windows: $(BIN_DIR)/khan.exe"

## 🐧 بیلد لینوکس
linux:
	@mkdir -p $(BIN_DIR)
	cd code && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ../$(BIN_DIR)/khan-linux ./cmd/server
	@echo "✅ linux: $(BIN_DIR)/khan-linux"

## 🍎 بیلد مک (هر دو معماری)
mac:
	@mkdir -p $(BIN_DIR)
	cd code && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ../$(BIN_DIR)/khan-macos-intel ./cmd/server
	cd code && GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ../$(BIN_DIR)/khan-macos-arm64 ./cmd/server
	@echo "✅ mac: intel + arm64"

## 🐳 ساخت داکر
docker:
	docker build -t khan:$(VERSION) .

## 🧪 تست
test:
	cd code && go test ./... -v

## 🔍 تحلیل استاتیک
vet:
	cd code && go vet ./...

## ✨ قالب‌بندی
fmt:
	cd code && gofmt -l .

## 🔬 لینتر
lint: vet fmt
	@echo "✅ lint passed"

## 🧹 پاک‌سازی
clean:
	rm -rf $(BIN_DIR)
	@echo "✅ cleaned"

## 🚀 اجرای محلی
run:
	cd code && go run ./cmd/server

## 🔐 بررسی امنیتی
check:
	bash scripts/khan-security-check.sh

## 📦 آماده انتشار
release: lint test build-all
	@echo "🎉 release $(VERSION) ready!"
	@ls -la $(BIN_DIR)/

## ❓ کمک
help:
	@echo "خان — اهداف Makefile:"
	@echo "  make build       بیلد برای سیستم فعلی"
	@echo "  make build-all   بیلد همه پلتفرم‌ها (win/linux/mac)"
	@echo "  make windows     فقط ویندوز"
	@echo "  make linux       فقط لینوکس"
	@echo "  make mac         فقط مک"
	@echo "  make docker      ساخت داکر"
	@echo "  make test        تست‌ها"
	@echo "  make vet         تحلیل استاتیک"
	@echo "  make lint        لینت + قالب‌بندی"
	@echo "  make run         اجرای محلی"
	@echo "  make check       بررسی امنیتی"
	@echo "  make release     انتشار کامل"
	@echo "  make clean       پاک‌سازی"
