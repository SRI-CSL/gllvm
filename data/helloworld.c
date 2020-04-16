#include <stdio.h>
#include "testobj.h"


int main(int argc, char *argv[]) {

  fprintf(stdout, "hello world\n");

  int testval = testobj_func1(41);
  fprintf(stdout, "This value should be 42: %d\n", testval);
}
