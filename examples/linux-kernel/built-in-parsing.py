
## trying to get the whole parsing of any config

import sys,os
import subprocess

excluded_dirs=[]

if len(sys.argv) > 2:
    excluded_dirs=sys.argv[2:]
script_out= sys.argv[1]

excluded=[path.split('/') for path in excluded_dirs]

builtin0 = open("arbi",'r')
out = open(script_out,"w+")
link_args = open("/home/vagrant/standalone-build/link-args","w+")
out.writelines("export build_home=/home/vagrant/standalone-build\n")

standalone_objects = ["arch/x86/kernel/head_64.o","arch/x86/kernel/head64.o","arch/x86/kernel/ebda.o","arch/x86/kernel/platform-quirks.o"]#,"usr/initramfs_data.o"]
archbi=["lib","pci","video","power"]
def write_script(builtin, excluded_paths, depth, base_dir):
    out.writelines("mkdir -p $build_home/built-ins/"+base_dir+" \n")

    directories=[]
    for line in builtin.readlines():
        line=line.split('\n')[0]
        words=line.split('/')
        if words[0]=="arch": 
            words[0]=words[0]+'/'+words[1]
            if words[2] in archbi:
                words[0]=words[0]+'/'+words[2]
        if words[depth] not in directories and line not in standalone_objects:
            directories.append(words[depth])

    roots = [path[depth] for path in excluded_paths]

    for direc in directories:
        print base_dir+direc
        if direc in roots:
            if base_dir+direc in excluded_dirs:
                out.writelines("convert-thin-archive.sh "+base_dir+direc+"/built-in.o \n")
                out.writelines("cp "+ base_dir +direc+"/built-in.o.new $build_home/built-ins/"+base_dir+direc+"bi.o \n \n")
                link_args.writelines("built-ins/"+base_dir + direc +"bi.o ")
            else:
                direc_roots = [path for path in excluded_paths if (path[depth]==direc and len(path)>depth+1) ]
                if direc_roots:
                    arbi=open(base_dir+direc+"/arbi","w+")
                    subprocess.call(["ar", "-t", base_dir+direc+"/built-in.o"],stdout=arbi )
                    arbi.close()
                    x_builtin = open(base_dir+direc+"/arbi","r")
                    write_script(x_builtin,direc_roots,depth+1,base_dir+direc+"/")
        else:
            if direc[-2:] == ".o":
                out.writelines("get-bc -b "+base_dir+direc+"\n")
                out.writelines("mkdir -p $build_home/built-ins/"+base_dir+"objects \n")
                out.writelines("cp "+ base_dir+direc+".bc $build_home/built-ins/"+base_dir+"objects \n")
                out.writelines("clang -c -no-integrated-as -mcmodel=kernel -o $build_home/built-ins/"+base_dir+ direc + " $build_home/built-ins/"+base_dir+"objects/" + direc+".bc \n \n")

                link_args.writelines("built-ins/"+base_dir + direc +" ")

            else:
                if os.path.isfile("/home/vagrant/standalone-build/gllvm.log"):
                    #direc_no_slash=direc.split('/')[0]
                    #os.rename("/home/vagrant/standalone-build/gllvm.log","/home/vagrant/standalone-build/gllvm"+direc_no_slash+".log")
                    os.remove("/home/vagrant/standalone-build/gllvm.log")
                path = base_dir + direc +"/built-in.o"
                subprocess.call(["get-bc", "-b", path ])
                if not os.path.isfile("/home/vagrant/standalone-build/gllvm.log"): subprocess.call(["touch","/home/vagrant/standalone-build/gllvm.log"])
                llvm_log=open("/home/vagrant/standalone-build/gllvm.log","r")
                assembly_objects=[]
                for line in llvm_log.readlines():
                    if len(line)>=54 and line[:54]=="WARNING:Error reading the .llvm_bc section of ELF file":
                        assembly_objects.append(line[55:-2])
                if assembly_objects:
                    out.writelines("mkdir -p $build_home/built-ins/" + base_dir +direc+"\n")
                for asf in assembly_objects:
                    out.writelines("cp "+ asf + " $build_home/built-ins/" + base_dir +direc +"\n")
                    filename= asf.split('/')[-1]
                    link_args.writelines("built-ins/"+base_dir + direc +'/'+filename+" ")
                #out.writelines("convert-thin-archive.sh "+base_dir + direc +"/built-in.o \n")
                #out.writelines("get-bc -b "+ base_dir + direc +"/built-in.o \n")
                if os.path.isfile(base_dir + direc +"/built-in.o.a.bc"):
                    out.writelines("cp "+ base_dir + direc +"/built-in.o.a.bc $build_home/built-ins/"+ base_dir + direc+"bi.o.bc \n")
                    out.writelines("clang -c -no-integrated-as -mcmodel=kernel -o $build_home/built-ins/"+ base_dir + direc + "bibc.o $build_home/built-ins/" + base_dir+direc+"bi.o.bc \n \n")
                    link_args.writelines("built-ins/"+base_dir + direc +"bibc.o ")

                builtin.close()

write_script(builtin0,excluded,0,"")


out.writelines("get-bc -b $build_home/lib/lib.a \n ")
out.writelines("cp lib/lib.a.bc $build_home/lib \n")
out.writelines("clang -c -no-integrated-as -mcmodel=kernel -o $build_home/lib/lib.a.o $build_home/lib/lib.a.bc \n")

if os.path.isfile("/home/vagrant/standalone-build/gllvm.log"):
    #direc_no_slash=direc.split('/')[0]
    #os.rename("/home/vagrant/standalone-build/gllvm.log","/home/vagrant/standalone-build/gllvm"+direc_no_slash+".log")
    os.remove("/home/vagrant/standalone-build/gllvm.log")
subprocess.call(["get-bc", "-b", "arch/x86/lib/lib.a" ])
subprocess.call(["touch","/home/vagrant/standalone-build/gllvm.log"])
llvm_log=open("/home/vagrant/standalone-build/gllvm.log","r")
assembly_objects=[]
for line in llvm_log.readlines():
    if len(line)>=54 and line[:54]=="WARNING:Error reading the .llvm_bc section of ELF file":
        assembly_objects.append(line[55:-2])
if assembly_objects:
    out.writelines("mkdir -p $build_home/arch/x86/lib/objects \n")
for asf in assembly_objects:
    out.writelines("cp "+ asf + " $build_home/arch/x86/lib/objects/ \n")
    filename= asf.split('/')[-1]
    #link_args.writelines("arch/x86/lib/objects/" +filename+" ")
out.writelines("cp arch/x86/lib/lib.a.bc $build_home/arch/x86/lib/lib.a.bc \n")
out.writelines("clang -c -no-integrated-as -mcmodel=kernel -o $build_home/arch/x86/lib/lib.a.o $build_home/arch/x86/lib/lib.a.bc \n \n")

out.writelines("cp arch/x86/kernel/vmlinux.lds $build_home \n")
out.writelines("cp .tmp_kallsyms2.o $build_home \n")
for sto in standalone_objects:
    out.writelines("cp --parents "+ sto+" $build_home \n")
out.writelines("\n#linking command \n")
out.writelines("cd $build_home \n")
out.writelines("ld --build-id -T vmlinux.lds --whole-archive ")
for sto in standalone_objects:
    out.writelines(sto+" ")
out.writelines("\@link-args ")
out.writelines("--no-whole-archive --start-group lib/lib.a.o arch/x86/lib/lib.a.o arch/x86/lib/objects/*  .tmp_kallsyms2.o --end-group -o vmlinux")