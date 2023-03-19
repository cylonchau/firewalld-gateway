## fiewall gateway

[中文](./README-CN.md)

fiewall gateway is a firewall central controller as firewalld

## Features

- Full dbus api.
- Full firewalld features.
- Based dbus remotely.
- Http restful api.
- Support Kubernetes HA
- Async batch task (only add).
- Can control thousands of linux machine via firewall gateway remotely.
- Support change tempate of thousands of machine fastly.
- Support wrong operation.
- Support delay command effect.
- Support IPtables NAT
- Support template (only enable db).
- Only HTTP API, without DB & Store.

## TODO
- UI 
- Authtication.
- Based Kubernetes HA.
- Prometheus Metics.
- waf SDK.


## deploy

```bash
git clone ..
make
```

## use

[HTTP API DOC](https://documenter.getpostman.com/view/12796679/UV5agGNr)

## FAQ

### Why not use ssh or ansible tools.

Because D-Bus support remotely and firewalld implemented full D-Bus API, so we can batch manage iptables rules via firealld.

### How diffrence your project and other

firewall gateway implemented full dbus API convert to HTTP API, so can control thousands of machine via gateway. And ohter project update iptables via agent scripts. or only run on one machines.


### Is enable D-Bus remotely safe?

We can open D-Bus port only accpet gateway's IP, so is safed

default if you machine hacked, enable of disable D-Bus remote, it doesn't make any sense. Because hacker can run any commond on your machine.

If you machine Is safe, so we can through open D-Bus port only accpet gateway's IP, so can management iptables rules via gateway and UI
