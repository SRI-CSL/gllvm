# Compiling Apache on Ubuntu


On a clean 14.04 machine I will build apache.

```
>pwd

/vagrant

>more /etc/lsb-release

DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=14.04
DISTRIB_CODENAME=trusty
DISTRIB_DESCRIPTION="Ubuntu 14.04.5 LTS"
```


## Step 0.

We want a recent version of `go` so we install `1.10`. Currently this requires some effort.

```
>sudo add-apt-repository ppa:gophers/archive
>sudo apt-get update
>sudo apt-get install golang-1.10-go


>export PATH=/usr/lib/go-1.10/bin:${PATH}
```



## Step 1.


Install `gllvm`.

```
>export GOPATH=/vagrant/go

>go get github.com/SRI-CSL/gllvm/cmd/...

>export PATH=${GOPATH}/bin:${PATH}
```

## Step 2.

I am only going to build apache, not apr, so I first install the prerequisites.

```
>sudo apt-get install llvm libclang-dev clang libapr1-dev libaprutil1-dev

``` 

Note `gclang` is agnostic with respect to llvm versions
so feel free to install a more recent version if you
wish. However, if you are going to use dragonegg the llvm version is
tightly coupled to the gcc and plugin versions you are using.


## Step 3.

  Configure the gllvm tool to be relatively quiet:

```
>export WLLVM_OUTPUT_LEVEL=WARNING
>export WLLVM_OUTPUT_FILE=/vagrant/apache-build.log
```


## Step 4.

 Fetch apache, untar, configure, then build:

```
>wget  https://archive.apache.org/dist/httpd/httpd-2.4.33.tar.gz

>tar xfz httpd-2.4.33.tar.gz

>cd httpd-2.4.33

>CC=gclang ./configure

>make
```

## Step 5.

Extract the bitcode.

```
>get-bc -m httpd

```

The `-m` flag tells the get-bc tool to also write a manifest listing all the bitcode modules
that were linked together to create the `httpd.bc` module.

```
 >more httpd.bc.llvm.manifest

/home/vagrant/httpd-2.4.33/.modules.o.bc
/home/vagrant/httpd-2.4.33/.buildmark.o.bc
/home/vagrant/httpd-2.4.33/server/.main.o.bc
/home/vagrant/httpd-2.4.33/server/.vhost.o.bc
/home/vagrant/httpd-2.4.33/server/.util.o.bc
/home/vagrant/httpd-2.4.33/server/.mpm_common.o.bc
/home/vagrant/httpd-2.4.33/server/.util_filter.o.bc
/home/vagrant/httpd-2.4.33/server/.util_pcre.o.bc
/home/vagrant/httpd-2.4.33/server/.exports.o.bc
/home/vagrant/httpd-2.4.33/server/.scoreboard.o.bc
/home/vagrant/httpd-2.4.33/server/.error_bucket.o.bc
/home/vagrant/httpd-2.4.33/server/.protocol.o.bc
/home/vagrant/httpd-2.4.33/server/.core.o.bc
/home/vagrant/httpd-2.4.33/server/.request.o.bc
/home/vagrant/httpd-2.4.33/server/.provider.o.bc
/home/vagrant/httpd-2.4.33/server/.eoc_bucket.o.bc
/home/vagrant/httpd-2.4.33/server/.eor_bucket.o.bc
/home/vagrant/httpd-2.4.33/server/.core_filters.o.bc
/home/vagrant/httpd-2.4.33/server/.util_expr_eval.o.bc
/home/vagrant/httpd-2.4.33/server/.config.o.bc
/home/vagrant/httpd-2.4.33/server/.log.o.bc
/home/vagrant/httpd-2.4.33/server/.util_fcgi.o.bc
/home/vagrant/httpd-2.4.33/server/.util_script.o.bc
/home/vagrant/httpd-2.4.33/server/.util_md5.o.bc
/home/vagrant/httpd-2.4.33/server/.util_cfgtree.o.bc
/home/vagrant/httpd-2.4.33/server/.util_time.o.bc
/home/vagrant/httpd-2.4.33/server/.connection.o.bc
/home/vagrant/httpd-2.4.33/server/.listen.o.bc
/home/vagrant/httpd-2.4.33/server/.util_mutex.o.bc
/home/vagrant/httpd-2.4.33/server/.mpm_unix.o.bc
/home/vagrant/httpd-2.4.33/server/.mpm_fdqueue.o.bc
/home/vagrant/httpd-2.4.33/server/.util_cookies.o.bc
/home/vagrant/httpd-2.4.33/server/.util_debug.o.bc
/home/vagrant/httpd-2.4.33/server/.util_xml.o.bc
/home/vagrant/httpd-2.4.33/server/.util_regex.o.bc
/home/vagrant/httpd-2.4.33/server/.util_expr_parse.o.bc
/home/vagrant/httpd-2.4.33/server/.util_expr_scan.o.bc
/home/vagrant/httpd-2.4.33/modules/core/.mod_so.o.bc
/home/vagrant/httpd-2.4.33/modules/http/.http_core.o.bc
/home/vagrant/httpd-2.4.33/modules/http/.http_protocol.o.bc
/home/vagrant/httpd-2.4.33/modules/http/.http_request.o.bc
/home/vagrant/httpd-2.4.33/modules/http/.http_filters.o.bc
/home/vagrant/httpd-2.4.33/modules/http/.chunk_filter.o.bc
/home/vagrant/httpd-2.4.33/modules/http/.byterange_filter.o.bc
/home/vagrant/httpd-2.4.33/modules/http/.http_etag.o.bc
/home/vagrant/httpd-2.4.33/server/mpm/event/.event.o.bc
/home/vagrant/httpd-2.4.33/os/unix/.unixd.o.bc
```
