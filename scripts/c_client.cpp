#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <cstdlib>
#include <ctime>
#include <chrono>
#include <sstream>
#include <iomanip>
#include <cstdio>
#include <cstring>
#include <unistd.h>
#include <numeric>
#include <random>

std::string GetTimeMs() {
    using namespace std::chrono;
    auto timepoint = system_clock::now();
    auto coarse = system_clock::to_time_t(timepoint);
    auto fine = time_point_cast<std::chrono::milliseconds>(timepoint);

    char buffer[sizeof "9999-12-31 23:59:59.999"];
    std::snprintf(buffer + std::strftime(buffer, sizeof buffer - 3,
                                         "%F %T.", std::localtime(&coarse)),
                  4, "%03lu", fine.time_since_epoch().count() % 1000);
    return buffer;
}

std::string CleanStr(std::string str) {
    str.erase(std::remove(str.begin(), str.end(), '\n'), str.end());
    return str;
}

int main() {
    char hostname[1024];
    if (gethostname(hostname, 1024) != 0) {
        std::cerr << "Error getting hostname" << std::endl;
        return 1;
    }
    const int n = 250;
    std::vector<int> iHost(n);
    std::vector<int> iDomain(n);
    std::iota(iHost.begin() + 1, iHost.end(), 1);
    std::iota(iDomain.begin(), iDomain.end(), 1);
    std::shuffle(std::begin(iHost), std::end(iHost), std::default_random_engine());
    std::shuffle(std::begin(iDomain), std::end(iDomain), std::default_random_engine());
    std::vector<std::string> list;
    list.reserve(1 + (n * n));
    for (auto d : iDomain) {
        for (auto h : iHost) {
            std::string url = "host" + std::to_string(h) + ".domain" + std::to_string(d) + ".rsx218-dox.cnam.fr\n";
            list.push_back(url);
        }
    }
    std::random_shuffle(list.begin(), list.end());
    std::cout << "client,reqID,url,timestamp" << std::endl;
    int counter = 1;
    for (auto url : list) {
        std::string clean_url = CleanStr(url);
        std::cout << hostname << "," << counter << "," << clean_url << "," << GetTimeMs() << std::endl;
        std::string command = "dnslookup " + clean_url + " 192.168.56.2:453  > /dev/null 2>&1 &";
        if (system(command.c_str()) != 0) {
            std::cerr << "Error running command: " << command << std::endl;
            return 1;
        }
        counter++;
    }
    return 0;
}
