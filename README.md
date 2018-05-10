[![Build Status](https://travis-ci.org/OrigamiLogic/gollum.svg?branch=master)](https://travis-ci.org/jensendw/gollum)

# Gollum
![My Precious!](https://github.com/jensendw/gollum/raw/master/gollum.gif)

Gollum collects consul key values and vault secrets from their respective paths and turns them into a set of exports that can be easily sourced.  

This is meant to be used as a temporary stop gap to doing this the right way which is adding support for both of these things inside your application.

## Download

You can find releases on the [releases page](https://github.com/OrigamiLogic/gollum/releases)

If I wanted to download this on a Mac I would do the following:
```shell
wget https://github.com/OrigamiLogic/gollum/releases/download/0.0.2/gollum-0.0.2.darwin.amd64
mv gollum-0.0.1.darwin.amd64 gollum
chmod 755 gollum
```
## Usage

```
  -consulurl string
    	URL for consul cluster
  -gatekeeperurl string
    	URL for vault gatekeeper
  -keypath string
    	Path to consul key values for configuration
  -vaultkeypath string
    	Vault path to key with secrets
  -vaulturl string
    	URL for vault
  -version
    	prints current app version
```

You would typically place this inside your docker container, but in this case we will just walk through a few scenarios

### Collecting Key Values from Consul

The following will connect to https://consul.domain.com and retrieve all the key values in the settings/myapp/production/ path and convert them into a source export file with environment variables, for example assume you have the following in consul
```
settings/myapp/production/APPLICATION_NAME=foo
settings/myapp/production/FAVORITE_DRINK=absinthe
settings/myapp/production/ANSWER_TO_THE_UNIVERSE=42
```

Running this command:
```
./gollum -consulurl=https://consul.domain.com -keypath=settings/myapp/production/
```
Produces this output:
```
export APPLICATION_NAME=foo
export FAVORITE_DRINK=absinthe
export ANSWER_TO_THE_UNIVERSE=42
```

### Collecting secrets from vault

Gollum assumes you're using [vault-gatekeeper-mesos](https://github.com/ChannelMeter/vault-gatekeeper-mesos) to handle the management of temporary tokens, AKA vault response wrapping.

Assuming your vault-gatekeeper-mesos instance is running at https://vault-gatekeeper.domain.com, your vault service is available at https://vault.domain.com:8200 and you have the following in vault at this path secret/myapp/production/secrets

```
secret/myapp/production/secrets
     MYSQL_PASSWORD=IloveTheSmellofSQLInTheMorning
     MONGODB_PASSWORD=CauseWhoWantsRelationships
     lowercase=shouting_is_fun
```

Running this command:
```
./gollum -gatekeeperurl=https://vault-gatekeeper.domain.com -vaulturl=https://vault.domain.com:8200 -vaultkeypath=secret/myapp/production/secrets
```
Produces this output:
```
export MYSQL_PASSWORD=IloveTheSmellofSQLInTheMorning
export MONGODB_PASSWORD=CauseWhoWantsRelationships
export lowercase=shouting_is_fun
```

### All the things
We can produce output for both consul and vault at the same time, given the assumptions of the examples above:
```shell
./gollum -consulurl=https://consul.domain.com -keypath=settings/myapp/production/ -gatekeeperurl=https://vault-gatekeeper.domain.com -vaulturl=https://vault.domain.com:8200 -vaultkeypath=secret/myapp/production/secrets
```

Would produce the following:
```
export APPLICATION_NAME=foo
export FAVORITE_DRINK=absinthe
export ANSWER_TO_THE_UNIVERSE=42
export MYSQL_PASSWORD=IloveTheSmellofSQLInTheMorning
export MONGODB_PASSWORD=CauseWhoWantsRelationships
export lowercase=shouting_is_fun
```

In the real world we would probably run gollum in an entrypoint.sh file within a docker container as follows:
```shell
./gollum -consulurl=https://consul.domain.com \
  -keypath=settings/myapp/production/ \
  -gatekeeperurl=https://vault-gatekeeper.domain.com \
  -vaulturl=https://vault.domain.com:8200
  -vaultkeypath=secret/myapp/production/secrets \
  > myvars

source myvars
```

This will make the values within vault and consul available as environment variables within your container.

## Development
1. go get https://github.com/OrigamiLogic/gollum
2. cd $GOPATH/src/github.com/OrigamiLogic/gollum
3. dep ensure
4. go test

## Contributing
1. Fork it ( https://github.com/OrigamiLogic/gollum )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
