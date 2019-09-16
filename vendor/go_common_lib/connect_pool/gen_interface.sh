#! /usr/bin/bash
# -- add one line cmd when you have a new recommend service
thrift --gen go:package_prefix=go_common_lib/connect_pool/gen-go/ thrift_src/rec_tag_search.thrift
thrift --gen go:package_prefix=go_common_lib/connect_pool/gen-go/ thrift_src/response.thrift
