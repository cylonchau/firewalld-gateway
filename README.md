## Fiewall Gateway Uranus

In Greek mythology, Uranus king of gods. The firewall gateway is the Uranus of iptables for many hosts

[中文](./README-CN.md)

fiewall gateway is a firewall central controller as firewalld

## Features

- Full firewalld features
- Full D-BUS API convert to REST API.
- Based dbus remotely.
- HTTP restful API.
- Support HA (Based Kubernetes)
- Asynchronous batch interface (only add).
- Can control thousands of linux machine via firewall gateway remotely.
- Support change tempate of thousands of machine fastly.
- Support wrong operation backoff.
- Support delay command effect.
- Support iptables NAT ipset timer task.
- Support template switch (only enable db).
- Only HTTP Service (without store).

## TODO
- [X] Asynchronous batch process
- [X] optional API on (v3 only)
- [ ] rpm spec
- [ ] Delay task
- [ ] UI
- [ ] Authtication.
- [ ] Based Kubernetes HA.
- [ ] Prometheus Metics.
- [ ] WAF SDK.
- [ ] Deplyment on Kubernetes


## Deploy on binary

```bash
git clone ..
make
```

## Deplyment on kubernetes

```

```

## Thanks libs
- [kubernetes workqueue](https://github.com/kubernetes/kubernetes)
- [klog](https://github.com/kubernetes/kubernetes)
- [godbus](https://github.com/godbus/dbus)
- [gin](https://.com/gin-gonic/gin)
- [viper](https://github.com/spf13/viper)

## use

[HTTP API DOC](https://documenter.getpostman.com/view/12796679/UV5agGNr)

- v1 runtime resource.
- v2 permanent resource.
- v3 Asynchronous batck opreation.

## FAQ

### Why not use ssh or ansible tools.

Because D-Bus support remotely and firewalld implemented full D-Bus API, so we can batch manage iptables rules via firealld.

### How diffrence your project and other

firewall gateway implemented full dbus API convert to HTTP API, so can control thousands of machine via gateway. And ohter project update iptables via agent scripts. or only run on one machines.


### Is enable D-Bus remotely safe?

We can open D-Bus port only accpet gateway's IP, so is safed

default if you machine hacked, enable of disable D-Bus remote, it doesn't make any sense. Because hacker can run any command on your machine.

If you machine Is safe, so we can through open D-Bus port only accpet gateway's IP, so can management iptables rules via gateway and UI
