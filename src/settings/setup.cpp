#include <iostream>
#include <fstream>
#include <string>

int main(int argc, char *argv[]) {

	// Create directors
	// bin, log, pkg

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

