#!/bin/bash

curl http://127.0.0.1:9090/admin/v1/plugins -X POST -i -H "Content-Type:application/json" -H "Accept:application/json" -w "\n" -d '{"type": "HTTPInput", "config": {"plugin_name": "test-httpinput", "url": "/test", "method": "POST", "headers_enum": {"name": ["bar", "bar1"]}, "request_body_io_key": "HTTP_REQUEST_BODY_IO"}}'
curl http://127.0.0.1:9090/admin/v1/plugins -X POST -i -H "Content-Type:application/json" -H "Accept:application/json" -w "\n" -d '{"type": "IOReader", "config": {"plugin_name": "test-ioreader", "input_key":"HTTP_REQUEST_BODY_IO", "output_key": "DATA"}}'
curl http://127.0.0.1:9090/admin/v1/plugins -X POST -i -H "Content-Type:application/json" -H "Accept:application/json" -w "\n" -d '{"type": "JSONValidator", "config": {"plugin_name": "test-jsonvalidator", "schema": "{\"title\": \"Record\",\"type\": \"object\",\"properties\": {\"name\": {\"type\": \"string\"}}, \"required\": [\"name\"]}", "data_key": "DATA"}}'
curl http://127.0.0.1:9090/admin/v1/plugins -X POST -i -H "Content-Type:application/json" -H "Accept:application/json" -w "\n" -d '{"type": "EaseMonitorJSONGidExtractor", "config": {"plugin_name": "test-jsongidextractor", "gid_key": "GID", "data_key": "DATA"}}'
curl http://127.0.0.1:9090/admin/v1/plugins -X POST -i -H "Content-Type:application/json" -H "Accept:application/json" -w "\n" -d '{"type": "KafkaOutput", "config": {"plugin_name": "test-kafkaoutput", "topic": "test", "brokers": ["192.168.98.130:9092"], "message_key_key": "GID", "data_key": "DATA"}}'
curl http://127.0.0.1:9090/admin/v1/pipelines -X POST -i -H "Content-Type:application/json" -H "Accept:application/json" -w "\n" -d '{"type": "LinearPipeline", "config": {"pipeline_name": "test-jsonpipeline", "plugin_names": ["test-httpinput", "test-ioreader", "test-jsonvalidator", "test-jsongidextractor", "test-kafkaoutput"], "parallelism": 10}}'
