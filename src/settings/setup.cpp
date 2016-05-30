#include <iostream>
#include <fstream>
#include <string>

#include <sys/types.h>
#include <sys/stat.h>
#include <stdio.h>
#include <stdlib.h>
#include <direct.h>

// http://stackoverflow.com/questions/18100097/portable-way-to-check-if-directory-exists-windows-linux-c
int dirExists(const char *path)
{
    struct stat info;

    if(stat( path, &info ) != 0)
        return 0;
    else if(info.st_mode & S_IFDIR)
        return 1;
    else
        return 0;
}

int main(int argc, char *argv[]) {

	// Create directors
	// bin, log, pkg
    const char *pathBin = "./bin/";
    printf("%d\n", dirExists(pathBin));
    const char *pathLog = "./log/";
    printf("%d\n", dirExists(pathLog));
    const char *pathPkg = "./pkg/";
    printf("%d\n", dirExists(pathPkg));
	// mkdir("c:/myfolder");

	// std::ofstream config_file;
	// config_file.open ("./settings.json");

	// std::string settings = "{\n";

	// std::cout << "server port: ";
	// std::string port;
	// std::cin >> port;
	// settings += "\t\"port\": " + port + ",\n";

	// std::cout << "database: ";
	// std::string db;
	// std::cin >> db;
	// db += ".db";
	// settings += "\t\"db\": \"" + db + "\",\n";

	// std::cout << "authkey: ";
	// std::string authkey;
	// std::cin >> authkey;
	// settings += "\t\"authkey\": \"" + authkey + "\"\n";

	// settings += "}";

	// config_file << settings;

	// config_file.close();

	return 0;

}

