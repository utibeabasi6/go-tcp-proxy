# Go tcp proxy

A simple tcp proxy written in golang

# How it works

The proxy listens on each of the ports defined in each of the apps in the config.json file. It proxies requests to any one of the targets defined and loadbalances between them.

# Testing