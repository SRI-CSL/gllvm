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

>mkdir -p ${GOPATH}

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
>get-bc httpd

>ls -la httpd.bc
-rw-r--r-- 1 vagrant vagrant 1119584 Aug  4 20:02 httpd.bc
```

## Step 6.  

Turn the bitcode into a second executable binary. (optional -- just for fun and sanity checking)

```
llc -filetype=obj httpd.bc
gcc httpd.o -lpthread -lapr-1 -laprutil-1 -lpcre -o httpd.new
```
