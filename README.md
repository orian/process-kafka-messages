# process-kafka-messages

A simple Go utility to parse Kafka message definition's and use go' template to generate
code to encode and decode messages.

For easy tests one may run:

```
git clone --depth=1 git@github.com:apache/kafka.git
go run github.com/orian/process-kafka-messages -dir ../kafka/clients/src/main/resources/common/message/
```
