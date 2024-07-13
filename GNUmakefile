default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 $(if $(PGVECTORS_CLOUD_API_KEY),PGVECTORS_CLOUD_API_KEY=$(PGVECTORS_CLOUD_API_KEY) )$(if $(PGVECTORS_CLOUD_API_URL),PGVECTORS_CLOUD_API_URL=$(PGVECTORS_CLOUD_API_URL) )go test ./... -v $(TESTARGS) -timeout 120m
