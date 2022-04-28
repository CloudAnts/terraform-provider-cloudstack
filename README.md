CloudStack Terraform Provider
=============================

A notice from CloudAnts
-----------------------

We, CloudAnts, are a managed hosting provider. We've forked the original repository to fix a number of bugs that have
been impeding our progress, and we've added a few missing features. At the moment, we have no intention of pushing back
changes to the original repository (through pull requests), but feel free to use this provider as you wish.

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 1.1.x
- [Go](https://golang.org/doc/install) >= 1.16+ (to build the provider plugin)

Using the Provider from Terraform registry
------------------------------------------
To install the CloudStack provider, copy and paste the below code into your Terraform configuration. Then, run terraform init.
```sh
terraform {
  required_providers {
    cloudstack = {
      source = "cloudants/cloudstack"
      version = "1.0.0"
    }
  }
}

provider "cloudstack" {
  # Configuration options
}
```
Check the [Terraform documentation](https://registry.terraform.io/providers/cloudants/cloudstack/latest/docs) for more
details on how to install and use the provider.

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version
1.16+ is *required*). You'll also need to correctly set up a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Clone repository to: `$GOPATH/src/github.com/apache/terraform-provider-cloudstack`

```sh
$ mkdir -p $GOPATH/src/github.com/cloudants; cd $GOPATH/src/github.com/cloudants
$ git clone git@github.com:cloudants/terraform-provider-cloudstack
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the
`$GOPATH/bin` directory. Enter the provider directory and build the provider:

```sh
$ cd $GOPATH/src/github.com/cloudants/terraform-provider-cloudstack
$ make build
$ ls $GOPATH/bin/terraform-provider-cloudstack
```

Once the build is ready, you have to copy the binary into Terraform locally (version appended).
On Linux this path is at ~/.terraform.d/plugins, and on Windows at %APPDATA%\terraform.d\plugins.

```sh
$ ls ~/.terraform.d/plugins/registry.terraform.io/cloudants/cloudstack/1.0.0/linux_amd64/terraform-provider-cloudstack_v1.0.0
```

You can also symlink the file locally to make testing a little easier:

```sh
$ ln -s $GOPATH/bin/terraform-provider-cloudstack ~/.terraform.d/plugins/registry.terraform.io/cloudants/cloudstack/1.0.0/linux_amd64/terraform-provider-cloudstack_v1.0.0
```

Testing the Provider
--------------------

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests you will need to run the CloudStack Simulator. Please note that the
official simulator is broken at this moment. The Docker image supplied by `ustcweizhou/cloudstack-simulator` seems a
good replacement, however the datacenter configuration created in the image doesn't seem to work with CloudStack
anymore.

```sh
$ docker pull ustcweizhou/cloudstack-simulator
$ docker run --name simulator -p 8080:5050 -d ustcweizhou/cloudstack-simulator
```

You can also set up a small Docker Compose project:

```yaml
version: "3.8"
services:
  simulator:
    container_name: cloudstack-simulator
    image: ustcweizhou/cloudstack-simulator
    ports:
      - "8080:5050"
```

When the container has started, you can go to http://localhost:8080/client and login to the CloudStack UI as user
`admin` with password `password`. It can take a few minutes for the container is fully ready, so you probably need to
wait and refresh the page for a few minutes before the login page is shown.

Once the login page is shown, and you can log in, you need to provision a simulated data-center:

```sh
$ docker exec -ti cloudstack python /root/tools/marvin/marvin/deployDataCenter.py -i /root/setup/dev/advanced.cfg
```

If you refresh the client or login again, you will now get passed the initial welcome screen and be able to go to your
account details and retrieve the API key and secret. Export those together with the URL:

```sh
$ export CLOUDSTACK_API_URL=http://localhost:8080/client/api
$ export CLOUDSTACK_API_KEY=<KEY>
$ export CLOUDSTACK_SECRET_KEY=<SECRET>
```

In order for all the tests to pass, you will need to create a new (empty) project in the UI called `terraform`.
When the project is created you can run the Acceptance tests against the CloudStack Simulator by simply running:

```sh
$ make testacc
```

Sample Terraform configuration
------------------------------
Below is an example configuration to initialize provider and create a Virtual Machine instance

```sh
$ cat provider.tf
terraform {
  required_providers {
    cloudstack = {
      source = "cloudants/cloudstack"
      version = "1.0.0"
    }
  }
}

provider "cloudstack" {
  # Configuration options
  api_url    = "${var.cloudstack_api_url}"
  api_key    = "${var.cloudstack_api_key}"
  secret_key = "${var.cloudstack_secret_key}"
}

resource "cloudstack_instance" "web" {
  name             = "server-1"
  service_offering = "Small Instance"
  network_id       = "df5fc279-86d5-4f5d-b7e9-b27f003ca3fc"
  template         = "616fe117-0c1c-11ec-aec4-1e00610002a9"
  zone             = "2b61ed5d-e8bd-431d-bf52-d127655dffab"
}
```

## History

This codebase has not been well maintained through 2022 Q1. As such, CloudAnts has forked the repository and made a few
necessary fixes in order to ensure that we can keep using CloudStack as intended.

## License

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at <http://www.apache.org/licenses/LICENSE-2.0>
