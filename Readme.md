## Exporter

Cli tool to export many rows of a complex SQL query to a csv file.

```shell
$ exporter -d=./data -q="select n.id id1, n.body, n2.id from notifications n inner join notifications n2 on n.id = n2.id"
```
