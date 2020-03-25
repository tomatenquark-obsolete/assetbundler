#include <iostream>
#include <string>
#include "libassetbundler.h"

int main() {
    GoString source = { "http://localhost:8000/packages/base/collide.cfg", 47 };
    const char* destination = DownloadMap(source);
    std::cout << destination << std::endl;
}
