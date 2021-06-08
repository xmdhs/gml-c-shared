cd c
gcc -c char.c -o char.o -O3
ar rcs libchar.a char.o

cd ..
go build -trimpath -ldflags "-w -s -linkmode \"external\" -extldflags \"-static -O3\"" -buildmode=c-shared -o cgo.dll