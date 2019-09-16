import random
fout=open("random.map", "w")
for line in range (0,1000):
    for col in range (0,1000):
        random_num=random.randint(0,100)
        fout.write(str(random_num))
        fout.write(" ")
    fout.write("\n")
fout.close()
