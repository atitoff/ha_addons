#!/usr/bin/env bashio

# #!/usr/bin/env bash
CONFIG_PATH=/data/options.json


MQTT_HOST=$(bashio::services mqtt "host")
MQTT_USER=$(bashio::services mqtt "username")
MQTT_PASSWORD=$(bashio::services mqtt "password")

if [[ -r "$CONFIG_PATH" ]]
then
  MqttPort="$(jq --raw-output '.MqttPort // empty' "$CONFIG_PATH")"
  MqttPort="$(jq --raw-output '.MqttHost // empty' "$CONFIG_PATH")"
fi

export MqttPort="${MqttPort:-1883}"
export MqttHost="${MqttHost:-localhost}"

exec /gpio_mqtt --mqtthost=$MQTT_HOST --mqttuser=$MQTT_USER --mqttpass=$MQTT_PASSWORD --mqttport=1883
