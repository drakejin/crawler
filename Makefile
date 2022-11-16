
.PHONY: gen
gen:
	go run cmd/entgo_generate/main.go

.PHONY: mig
mig:
	go run cmd/entgo_migration/main.go

.PHONY: genmig
genmig: gen mig