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

go build

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
      "uuid": "f2cd9970-5904-4525-b7e4-14310aa98119",
      "version": "4.5.0-snapshot"
    }

And we need to add the image file:

    @ todo update me

Let's add it as a source:

    root@smartos ~> imgadm sources -a http://norbye.ddns.net/
    Added "imgapi" image source "http://norbye.ddns.net/"
    root@smartos ~> imgadm sources
    https://images.joyent.com
    https://datasets.joyent.com/datasets
    http://norbye.ddns.net/

At this time we may start creating a new image


Run server under SMF
--------------------

I've not created the files yet :-)



[image_api_link]: https://images.joyent.com/docs/#api-summary
[couchbase_link]: http://www.couchbase.com/
[smartos_serving_images_link]: https://wiki.smartos.org/display/DOC/Managing+Images#ManagingImages-ServingImages
[dsapid_link]: https://github.com/MerlinDMC/dsapid
[golang_link]: https://golang.org/
