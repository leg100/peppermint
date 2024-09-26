#!/usr/bin/env bash

# in foreground, continously run app
while true; do
    _build/moto postgres:///lr land_registry_price_paid_uk
done
