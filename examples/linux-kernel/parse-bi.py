import sys
arbi = sys.argv[1]
outs= sys.argv[2]
outlists = sys.argv[3]
bc_pos= int(sys.argv[4])

arg1 = arbi.split('/')
folder = arg1[0]
bi = open(arbi,"r")
out= open(outs,"w")
outlist= open (outlists,"w")

dir_set = []
for line in bi.readlines():
    line_words = line.split('/')
    if line_words[0] not in dir_set:
        dir_set.append(line_words[0])

out.writelines("export build_home=$HOME/standalone-build\n")

for direc in dir_set[2:-1]:
    if direc[-2:] == ".o":
        out.writelines("get-bc -b "+direc+"\n")
        out.writelines("cp "+ direc+".bc $build_home/built-ins/"+folder+"/objects \n")
        out.writelines("cp "+ direc+" $build_home/built-ins/"+folder+"/objects \n \n")
    else:
        out.writelines("convert-thin-archive.sh "+direc+"/built-in.o \n")
        out.writelines("get-bc -b "+direc+"/built-in.o \n")
        out.writelines("cp "+ direc+"/built-in.o.a.bc $build_home/built-ins/"+folder+"/"+direc+"bi.o.bc \n")
        out.writelines("cp "+ direc+"/built-in.o.new $build_home/built-ins/"+folder+"/"+direc+"bi.o \n \n")
n=0
for direc in dir_set[2:-1]:
    outlist.writelines("built-ins/"+folder+"/")

    if direc[-2:] == ".o":

        outlist.writelines("objects/")
        outlist.writelines(direc[:-2] + "bc.o ")
    else:
        if (bc_pos==-1 or n<=bc_pos) and n!=7:
            outlist.writelines(direc + "bibc.o ")
        else:
            outlist.writelines(direc + "bi.o ")
        n+=1
    
