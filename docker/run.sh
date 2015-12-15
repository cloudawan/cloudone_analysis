#!/bin/bash

# Check whether Elastic Search is up
# Parameter
retry_amount=120
elastic_search_host_list_with_quotation=${ELASTICSEARCH_CLUSTER_HOST//,}
elastic_search_host_list=${elastic_search_host_list_with_quotation//\"}
elastic_search_port=$ELASTICSEARCH_CLUSTER_PORT

is_elastic_search_up() {
  for elastic_search_host in $elastic_search_host_list
  do
    elastic_search_url="http://$elastic_search_host:$elastic_search_port"
    elastic_search_response=$(curl -m 1 "$elastic_search_url")

    if [[ $elastic_search_response == *"\"elasticsearch\""* ]]; then
      return 1
    fi
  done
  return 0
}

# Elastic Search
for ((i=0;i<$retry_amount;i++))
do
  echo "ping $i times to Elastic Search"
  is_elastic_search_up
  elastic_search_result=$?
  if [ $elastic_search_result == 1 ]; then
	break
  fi
  sleep 1
done

if [ $i == $retry_amount ]; then
  echo "Could not get ping response from Elastic Search"
  exit -1
fi



# Use environment
sed -i "s/{{KUBEAPI_CLUSTER_HOST_AND_PORT}}/$KUBEAPI_CLUSTER_HOST_AND_PORT/g" /etc/cloudone_analysis/configuration.json
sed -i "s/{{ELASTICSEARCH_CLUSTER_HOST}}/$ELASTICSEARCH_CLUSTER_HOST/g" /etc/cloudone_analysis/configuration.json
sed -i "s/{{ELASTICSEARCH_CLUSTER_PORT}}/$ELASTICSEARCH_CLUSTER_PORT/g" /etc/cloudone_analysis/configuration.json

cd /src/cloudone_analysis
./cloudone_analysis &

while :
do
	sleep 1
done

