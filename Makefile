run:
	go run exporter -d=./data -q="select n.id id1, n.body, n2.id from notifications n inner join notifications n2 on n.id = n2.id"
build:
	go build exporter
test:
	go test exporter -v
