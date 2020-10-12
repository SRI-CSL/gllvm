//https://llvm.org/docs/LinkTimeOptimization.html
#include <stdio.h>
#include "a.h"

void foo4(void) {
  printf("Hi\n");
}

int main() {
  return foo1();
}
