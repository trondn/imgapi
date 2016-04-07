imgapi
======

imgapi is a server that provides the [Image API][image_api_link] used
by `imgadm` to manage datasets. The current implementation only provides
a subset of the commands, but enough to be used for me at home.

Why I wrote my own
------------------

I work for [Couchbase][couchbase_link] and wanted to create my own images
with our server (for LX zones and SmartOS). Having my own image server
makes my life easier for my own testing.

So why didn't I just use an existing server? Well, I tried to follow
the documentation on the [SmartOS wiki][smartos_serving_images_link],
but it looks like the code (and instructions on how to build) is outdated.
At least I failed to build it, so I asked on the SmartOS mailing list for
alternatives (given that the server was using the "old" dataset API and
not the new API).

The response I got was to use [dsapid][dsapid_link]. So I installed that
and tried to figure out how I could add my own image there (I had no
intention of running a "mirror", just my own one). At this time I hadn't
looked at the [API][image_api_link]  and ehh I barely knew anything
about [golang][golang_link]. The comment: `definitely the upload machanism
is different` scared me a bit.

Wait a second! You said you didn't know [golang][golang_link], but this looks
like it is implemented in go! Good observation! I've had it on my long todo
list to take a look at go, but I've never found the time to do that (and
I haven't had a good reason to do so). I've written another tiny program
in go earlier (everything included in a single file), and with the size of
the [image api][image_api_link] I figured this would be a good time to
extend my experience with go :-)

So; Sorry for not contributing to an existing project by enhancing the
documentation (if it is already implemented).

This works for me, and I'm just happy if it may be of any use for others.

Limitations
-----------

As stated above I needed a server I could play around with testing out
my own images on my own servers. The API provides a lot of functionality
you can't utilize from the standard `imgadm` tool, and I just skipped
those.

The following features is not supported (there may be more):

 * All retrieval operations are public (but I haven't found a way to
   have `imgadm` provide credentials when adding a source anyway)
 * channels
 * acl
 * export / import
 * copy-remote
 * import-remote

Build
-----

If you've got your go build environment all set up you should be
able to get `imgapi` by simply executing:

    trond@ok ~> go get github.com/trondn/imgapi

 And you'll find the binary in `${GOPATH}/bin`

Run command
------------

    imgapi [-s] [-c configfile]


`-s`             - Start as a server

`-c configfile`  - Use `configfile` instead of `$HOME/.imgapi.json`

Configuration file
------------------

    {
        "datadir" : "/data/imgapi/files",
        "port" : 8080,
        "host" : "127.0.0.1",
        "userdb" : [
                {
                    "name" : "admin",
                    "password" : "secret"
                }
        ]
    }

`datadir` specifies the root directory where the server should store all
of the images to serve.

`port` specifies the port the server should listen to.

`host` specifies the hostname the server is listening on (used by the
client interface)

`userdb` is a list of credentials the user may provide in order to perform
operations that modifies the content on the server.


Example
-------

This example describes the setup I am using at home (the username/password
is of course different ;-)). I've configured my router to map the external
port 80 to port 8080 on the machine running my image server. This is why
you will see port 8080 in the configuration for the image server, but
on my smartos server you'll see that I don't specify that port. I could
of course connect _directly_ to my image server instead of going through
my router, but I normally don't do that because then I have to remember
to edit that if I end up writing a blog ;-)

In this example the host named `imgadmsrv` is the machine running our
image server (It may run any operating system as it don't have any external
dependencies). The host named `smartos` is my SmartOS server where I create
all my images.

So lets get started and start the image server:

    me@imgadmsrv ~> cat config.json
    {
        "datadir" : "/home/imgapi/files",
        "port" : 8080,
        "host" : "norbye.ddns.net",
        "userdb" : [
                {
                    "name" : "admin",
                    "password" : "secret"
                }
        ]
    }
    me@imgadmsrv ~> ls -l files
    total 0
    me@imgadmsrv ~> ./imgadmsrv -s -c config.json

Now that the image server is running lets test that it is working

    root@smartos ~> curl -i http://norbye.ddns.net/ping
    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Server: Norbye Public Images Repo
    Date: Thu, 31 Mar 2016 08:37:59 GMT
    Content-Length: 76

    {
      "imgapi": true,
      "pid": 93183,
      "ping": "pong",
      "version": "1.0.0"
    }

Time to create our first image.

    root@smartos ~> curl -X POST -u admin:secret --data-binary '
        {
           "name": "couchbase-server",
           "version": "4.5.0-snapshot",
           "description": "Snapshot of Couchbase Server 4.5.0 Community Edition",
           "type": "zone-dataset",
           "os": "smartos",
           "requirements": {
             "networks": [
               {
                 "name": "net0",
                 "description": "public"
               }
             ]
           }
        }' http://norbye.ddns.net/images

The server returns the newly created manifest:

    {
      "description": "Snapshot of Couchbase Server 4.5.0 Community Edition",
      "disabled": false,
      "name": "couchbase-server",
      "os": "smartos",
      "public": false,
      "requirements": {
        "networks": [
          {
            "description": "public",
            "name": "net0"
          }
        ]
      },
      "state": "unactivated",
      "type": "zone-dataset",
      "uuid": "6de01e97-d7ec-4906-bd8b-cb4eafdb7c8b",
      "v": 2,
      "version": "4.5.0-snapshot"
    }

And we need to add the image file to the generated UUID:

    root@smartos ~> curl --upload-file couchbase-watson-4.5.mg.gz \
                         -u admin:secret \
			 "http://norbye.ddns.net/images/f2cd9970-5904-4525-b7e4-14310aa98119/file?sha1=`digest -a sha1 couchbase-watson-4.5.img.gz`;compression=gzip"

At this time the server returns the full manifest:

    {
      "description": "Snapshot of Couchbase Server 4.5.0 Community Edition",
      "disabled": false,
      "files": [
        {
          "compression": "gzip",
          "sha1": "db8ef7fde4395a58b1a169e644bac3804cce62d9",
          "size": 379882274
        }
      ],
      "name": "couchbase-server",
      "os": "smartos",
      "public": false,
      "requirements": {
        "networks": [
          {
            "description": "public",
            "name": "net0"
          }
        ]
      },
      "state": "unactivated",
      "type": "zone-dataset",
      "uuid": "6de01e97-d7ec-4906-bd8b-cb4eafdb7c8b",
      "v": 2,
      "version": "4.5.0-snapshot"
    }

We can now activate this image with:

    root@smartos ~> curl -X POST -u admin:secret "http://norbye.ddns.net/images/6de01e97-d7ec-4906-bd8b-cb4eafdb7c8b?action=activate"

And the server respond with the updated manifest:

    {
      "description": "Snapshot of Couchbase Server 4.5.0 Community Edition",
      "disabled": false,
      "files": [
        {
          "compression": "gzip",
          "sha1": "db8ef7fde4395a58b1a169e644bac3804cce62d9",
          "size": 3.79882274e+08
        }
      ],
      "name": "couchbase-server",
      "os": "smartos",
      "public": false,
      "requirements": {
        "networks": [
          {
            "description": "public",
            "name": "net0"
          }
        ]
      },
      "state": "active",
      "type": "zone-dataset",
      "uuid": "6de01e97-d7ec-4906-bd8b-cb4eafdb7c8b",
      "v": 2,
      "version": "4.5.0-snapshot"
    }

Let's add it as a source:

    root@smartos ~> imgadm sources -a http://norbye.ddns.net/
    Added "imgapi" image source "http://norbye.ddns.net/"
    root@smartos ~> imgadm sources
    https://images.joyent.com
    https://datasets.joyent.com/datasets
    http://norbye.ddns.net/

Now look for our image:

    root@smartos ~> imgadm avail | grep couchbase
    6de01e97-d7ec-4906-bd8b-cb4eafdb7c8b  couchbase-server        4.5.0-snapshot  smartos  -

And import it:

    root@smartos ~> imgadm import 6de01e97-d7ec-4906-bd8b-cb4eafdb7c8b

Run server under SMF
--------------------

You may find a service manifest file in support-files/imgapid.xml. In order
to use that file you have to:

* Have a user named `imgapid` in the group `imgapid`
* Have a configuration file in `/opt/local/etc/imgapid/configuration.json`
* Install the binary as `/opt/local/sbin/imgapid`

The following commands creates a usable setup for you:

    groupadd imgapid
    roleadd -g imgapid -c "Image API Server" -d /data/imgapid -s /bin/false imgapid
    mkdir -p /data/imgapid
    chown imgapid:imgapid /data/imgapid
    mkdir /opt/local/etc/imgapid/
    cat > /opt/local/etc/imgapid/configuration.json <<EOF
    {
      "datadir" : "/data/imgapid",
      "port" : 8080,
      "host" : "127.0.0.1",
      "userdb" : [
        {
          "name" : "admin",
          "password" : "secret"
        }
      ]
    }
    EOF
    cp imgapi /opt/local/sbin/imgapid
    svccfg import support-files/imgapid.xml
    svcadm enable imgapid


<!-- links -->
[image_api_link]: https://images.joyent.com/docs/#api-summary
[couchbase_link]: http://www.couchbase.com/
[smartos_serving_images_link]: https://wiki.smartos.org/display/DOC/Managing+Images#ManagingImages-ServingImages
[dsapid_link]: https://github.com/MerlinDMC/dsapid
[golang_link]: https://golang.org/
