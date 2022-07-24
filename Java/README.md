# Desired state
A single application which can be used to serialize user specified payloads
using specified serializers (in our case AVRO or JSON) and publish that message 
dynamically to the topic specified by the user. Additionally, would like the topic
to be allowed to be usable for multiple type publications.


## Setup
Below values can be changed in the `application.yaml`
- Port `18089` free to run on
- Kafka cluster running on `127.0.0.1:9092`
- Schema Registry running `127.0.0.1:8081`
  - With some test schemas preloaded
- Postman (or equivalent tool) to POST data to app
- Available profiles `default`, `JSON`, `AVRO`

## API Requests
### Payload
- **topic**
  - _Data-Type_: string
  - _Description_: the topic for the message to be published to
  - _Required_: ALWAYS
- **avroSource**
  - _Data-Type_: string
  - _Description_: the subject name your AVRO schema is registered under
  - _Required_: ALWAYS
- **payload**
  - _Data-Type_: JSON
  - _Description_: the data contained in your message, this is any valid JSON token
  - _Required_: AVRO
- **headers**
  - _Data-Type_: JSON Map Key(string) Value(string)
  - _Description_: key-value data added as headers to message
  - _Required_: NO

#### Notes
- Content type header set based on API header is `content-type` change in `MessageUtils.java`

### Produce AVRO encoded message
    POST http://127.0.0.1:18089/avro
    {
        "topic": "<your topic here>",
        "avroSource": "<your schema subject>",
        "payload": {
            "<examplePayloadKey1>": "<examplePayloadValue1>",
            "<examplePayloadKey2>": "<examplePayloadValue2>"
        },
        "headers": {
            "<exampleHeaderKey1>": "<exampleHeaderValue1>",
            "<exampleHeaderKey2>": "<exampleHeaderValue2>"
        }
    }


### Produce JSON formatted message
    POST http://127.0.0.1:18089/json
    {
        "topic": "<your topic here>",
        "payload": {
            "<examplePayloadKey1>": "<examplePayloadValue1>",
            "<examplePayloadKey2>": "<examplePayloadValue2>"
        },
        "headers": {
            "<exampleHeaderKey1>": "<exampleHeaderValue1>",
            "<exampleHeaderKey2>": "<exampleHeaderValue2>"
        }
    }

## Optional Configurations
Using either `AVRO` or `JSON` as the spring profile locks the application into producing events of that type.
Doing this will configure the application to only expose the rest API corresponding to the datatype.

## Notes
- The `converter` and `org.apache.avro.io` packages placed at `src/main/java` are used to represent the below dependency
while preventing security vulnerabilities.

```
<!-- https://mvnrepository.com/artifact/tech.allegro.schema.json2avro/converter -->
<dependency>
<groupId>tech.allegro.schema.json2avro</groupId>
<artifactId>converter</artifactId>
<version>0.2.2</version>
</dependency>
```