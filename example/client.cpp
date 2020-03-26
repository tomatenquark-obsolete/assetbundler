#include <iostream>
#include <string>
#include "libassetbundler.h"

int main() {
    char* servercontent = "http://localhost:8000";
    char* map = "curvedm";
    int status = 1;
    char* zippath = StartDownload(servercontent, map);
    std::cout << zippath << std::endl;
    while (status == 1) {
        std::cout << "Do I get here?" << std::endl;
        status = GetStatus(zippath);
    }
    std::cout << "Downloaded everything" << std::endl;
}
