#include <iostream>
#include <string>
#include "libassetbundler.h"

int main() {
    const char* destination = DownloadMap("http://localhost:8000", "curvedm");
    std::cout << destination << std::endl;
}
