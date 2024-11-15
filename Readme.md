## Exporter

Cli tool to export a lot of unordered rows of a complex SQL query to a csv file.

```shell
$ exporter -f=exported_data -d=./data -q="select * from notifications"
```

### Performance

~ 2.2M rows simple table takes 10 seconds in my laptop. Exporter creates X goroutines to process the data depending on the cores of the machine.


```shell
λ time ./exporter -q="select client_id client, client_group_id grup, body, created_at from notifications"
2024/11/14 17:55:24.305837 - [INFO] - Snapshot created.
2024/11/14 17:55:24.837202 - [INFO] - Exporting records from 0 to 110356.
2024/11/14 17:55:24.916322 - [INFO] - Exporting records from 2207120 to 2317476.
2024/11/14 17:55:24.916477 - [INFO] - Exporting records from 1324272 to 1434628.
2024/11/14 17:55:24.916478 - [INFO] - Exporting records from 1213916 to 1324272.
2024/11/14 17:55:24.916490 - [INFO] - Exporting records from 331068 to 441424.
2024/11/14 17:55:24.916476 - [INFO] - Exporting records from 110356 to 220712.
2024/11/14 17:55:24.916508 - [INFO] - Exporting records from 1434628 to 1544984.
2024/11/14 17:55:24.916486 - [INFO] - Exporting records from 220712 to 331068.
2024/11/14 17:55:24.916540 - [INFO] - Exporting records from 1103560 to 1213916.
2024/11/14 17:55:24.916613 - [INFO] - Exporting records from 993204 to 1103560.
2024/11/14 17:55:24.916622 - [INFO] - Exporting records from 772492 to 882848.
2024/11/14 17:55:24.916636 - [INFO] - Exporting records from 441424 to 551780.
2024/11/14 17:55:24.916646 - [INFO] - Exporting records from 551780 to 662136.
2024/11/14 17:55:24.916661 - [INFO] - Exporting records from 662136 to 772492.
2024/11/14 17:55:24.916670 - [INFO] - Exporting records from 1986408 to 2096764.
2024/11/14 17:55:24.916683 - [INFO] - Exporting records from 1765696 to 1876052.
2024/11/14 17:55:24.916685 - [INFO] - Exporting records from 882848 to 993204.
2024/11/14 17:55:24.916691 - [INFO] - Exporting records from 1876052 to 1986408.
2024/11/14 17:55:24.916702 - [INFO] - Exporting records from 1655340 to 1765696.
2024/11/14 17:55:24.916712 - [INFO] - Exporting records from 2096764 to 2207120.
2024/11/14 17:55:24.916720 - [INFO] - Exporting records from 1544984 to 1655340.
2024/11/14 17:55:26.288214 - [INFO] - Count records exported: 2207133
2024/11/14 17:55:26.288255 - [INFO] - Data was exported: ./data_1731603324.csv
2024/11/14 17:55:26.297787 - [INFO] - Snapshot removed.

real	0m10,064s
user	0m1,832s
sys	0m0,433s

exporter on  main 
λ wc -l data_1731603324.csv 
2207134 data_1731603324.csv

```
