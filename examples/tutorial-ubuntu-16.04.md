# Compiling Apache on Ubuntu


On a clean 16.04 server machine I will build apache.  Desktop instructions should be no different.

```
>more /etc/lsb-release

DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=16.04
DISTRIB_CODENAME=xenial
DISTRIB_DESCRIPTION="Ubuntu 16.04.4 LTS"
```

## Step 0. 

Install the `go` language.
```
>sudo apt-get install golang
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
>sudo apt-get install llvm libclang-dev clang libapr1-dev libaprutil1-dev libpcre3-dev make

```

At this point, you could check your clang version with `which clang` and `ls -l /usr/bin/clang`.
It should be at least clang-3.8.

## Step 3.

  Configure the gllvm tool to be relatively quiet:

```
>export WLLVM_OUTPUT_LEVEL=WARNING
>export WLLVM_OUTPUT_FILE=/vagrant/apache-build.log
```

## Step 4.

 Fetch apache, untar, configure, then build:

```

>wget https://archive.apache.org/dist/httpd/httpd-2.4.33.tar.gz

>tar xfz httpd-2.4.33.tar.gz

>cd httpd-2.4.33

>CC=gllvm ./configure

>make
```

## Step 5.

Extract the bitcode.

```
>get-bc -m httpd

>ls -la httpd.bc
-rw-rw-r-- 1 vagrant vagrant 1172844 May  8 22:39 httpd.bc
```

The `-m` flag instructs the `get-bc` tool to write a manifest too, the manifest lists all the bitcode modules
that were linked together to create the `httpd.bc` module.

```
more httpd.bc.llvm.manifest

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

## Step 6.  

Turn the bitcode into a second executable binary. (optional -- just for fun and sanity checking)

```
llc -filetype=obj httpd.bc
clang httpd.o  -Wl,--export-dynamic -lpthread -lapr-1 -laprutil-1 -lpcre -o httpd_from_bc
```
See [here](http://tldp.org/HOWTO/Program-Library-HOWTO/shared-libraries.html) for an explanation of the
```
-Wl,--export-dynamic
```
incantation. The salient point being
```
"In some cases, the call to gcc to create the object file will also need to include the option 
'-Wl,-export-dynamic'. Normally, the dynamic symbol table contains only symbols which 
are used by a dynamic object. This option (when creating an ELF file) adds all symbols to the
dynamic symbol table (see ld(1) for more information). You need to use this option when there are 
'reverse dependencies', i.e., a DL library has unresolved symbols that by convention must be 
defined in the programs that intend to load these libraries."
```
