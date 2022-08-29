#!/usr/bin/env bash

CONFIG_PATH=/data/options.json

if [[ -r "$CONFIG_PATH" ]]
then
  MQTT_PORT="$(jq --raw-output '.MqttPort // empty' "$CONFIG_PATH")"
fi

export MQTT_PORT="${MQTT_PORT:-pi}"

exec /gpio_mqtt

