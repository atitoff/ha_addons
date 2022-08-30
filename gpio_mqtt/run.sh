#!/usr/bin/env bash

CONFIG_PATH=/data/options.json

if [[ -r "$CONFIG_PATH" ]]
then
  MqttPort="$(jq --raw-output '.MqttPort // empty' "$CONFIG_PATH")"
  MqttPort="$(jq --raw-output '.MqttHost // empty' "$CONFIG_PATH")"
fi

export MqttPort="${MqttPort:-1883}"
export MqttHost="${MqttHost:-localhost}"

exec /gpio_mqtt
