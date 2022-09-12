# Localstack stuff

```sh
docker run --rm -it -p 4566:4566 -p 4571:4571 localstack/localstack
```

### Install or update the AWS CLI

Follow these steps from the command line to install the AWS CLI on Linux.

* Linux x86 (64-bit)

```sh
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

* Config AWS



```
awslocal rds create-db-cluster --db-cluster-identifier core-trans --engine postgresql --database-name core-trans 
{
    "DBCluster": {
        ...
        "Endpoint": "localhost:4510",
        "Port": 4510,  # may vary
        "DBClusterArn": "arn:aws:rds:us-east-1:000000000000:cluster:db1",
        ...
    }
}


aws secretsmanager create-secret --name dbpass --secret-string test
{
    "ARN": "arn:aws:secretsmanager:eu-central-1:1234567890:secret:dbpass-cfnAX",
    "Name": "dbpass",
    "VersionId": "fffa1f4a-2381-4a2b-a977-4869d59a16c0"
}

aws rds-data execute-statement --database test --resource-arn arn:aws:rds:us-east-1:000000000000:cluster:db1 --secret-arn arn:aws:secretsmanager:eu-central-1:1234567890:secret:dbpass-cfnAX --include-result-metadata --sql 'SELECT 123'
{
    "columnMetadata": [
        {
            "arrayBaseColumnType": 0,
            "isAutoIncrement": false,
            "isCaseSensitive": false,
            "isCurrency": false,
            "isSigned": true,
            "label": "?column?",
            "name": "?column?",
            "nullable": 0,
            "precision": 10,
            "scale": 0,
            "schemaName": "",
            "tableName": "",
            "type": 4,
            "typeName": "int4"
        }
    ],
    "numberOfRecordsUpdated": 0,
    "records": [
        [
            {
                "longValue": 123
            }
        ]
    ]
}
```