# Compiling Apache on Ubuntu


On a clean 14.04 machine I will build apache.

```
>pwd

/vagrant

>more /etc/lsb-release

DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=14.04
DISTRIB_CODENAME=trusty
DISTRIB_DESCRIPTION="Ubuntu 14.04.2 LTS"
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
>get-bc httpd

```

