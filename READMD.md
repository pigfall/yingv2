# Intoduction
A toy VPN server side.

[Android client link](https://github.com/pigfall/ymobile)

## Build from source
### Prerequries
* Go >= 1.17
```
git clone https://github.com/pigfall/yingv2.git
cd yingv2/cmd/server
go build .
```

## Usage
### Prerequired
* Linux system
```
# output config template
server -demoConfig > config.json
# run server
sudo ./server -confPath config.json
```
