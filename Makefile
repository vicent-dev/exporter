run:
	go run exporter -f=notifications -d=. -q="select * from notifications limit 2000000"
build:
	go build exporter
test:
	go test exporter -v
