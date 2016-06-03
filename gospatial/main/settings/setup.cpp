#include <iostream>
#include <fstream>
#include <string>

#include <sys/types.h>
#include <sys/stat.h>
#include <stdio.h>
#include <stdlib.h>
// #include <direct.h>

// http://stackoverflow.com/questions/18100097/portable-way-to-check-if-directory-exists-windows-linux-c
int dirExists(const char *path) {
	struct stat info;
	if(stat( path, &info ) != 0)
		return 0;
	else if(info.st_mode & S_IFDIR)
		return 1;
	else
		return 0;
}

int main(int argc, char *argv[]) {

	// Create required directories
	// bin, log, pkg
	// const char *pathBin = "./bin/";
	// const char *pathLog = "./log/";
	// const char *pathPkg = "./pkg/";
	// if (!dirExists(pathBin)) {
	// 	mkdir(pathBin,  S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);
	// }
	// if (!dirExists(pathLog)) {
	// 	mkdir(pathLog,  S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);
	// }
	// if (!dirExists(pathPkg)) {
	// 	mkdir(pathPkg,  S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);
	// }

	// Create settings file
	std::ofstream config_file;
	config_file.open ("./settings.json");

	std::string settings = "{\n";

	std::cout << "server port: ";
	std::string port;
	std::cin >> port;
	settings += "\t\"port\": " + port + ",\n";

	std::cout << "database: ";
	std::string db;
	std::cin >> db;
	db += ".db";
	settings += "\t\"db\": \"" + db + "\",\n";

	std::cout << "authkey: ";
	std::string authkey;
	std::cin >> authkey;
	settings += "\t\"authkey\": \"" + authkey + "\"\n";

	settings += "}";

	config_file << settings;

	config_file.close();

	return 0;

}

