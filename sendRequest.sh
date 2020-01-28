#!/bin/bash

### Base distributed task Rest API url
baseRestApiUrl="https://32k2peqdgx7ub2us2xhkoumku5sysyikusqo7dxfn3whzdljfuqq@dev.azure.com/karauctionservices/openlane/_apis/build/definitions?api-version=5.0"
curl -X POST -H "Content-Type: application/json" -d @buildDefinition.json $baseRestApiUrl
