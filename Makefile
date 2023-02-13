mocks: mocks-clean mocks-generate

mocks-generate:
	mockery --all --output=mocks --case=underscore --keeptree

mocks-clean:
	rm -rf mocks/*