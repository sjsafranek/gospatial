#include <iostream>
#include <fstream>
// #include <stdlib.h>
#include <string>

// class Person {
// 	private:
// 		int age;
// 	public:
// 		char name;
// 		void input_age() {
// 			std::cout << "enter int\n";
// 			std::cin >> age;
// 		}
// 		void output_name() {
// 			std::cout << "My name is: ";
// 			std::cout << name << "\n";
// 		}
// };

int main(int argc, char *argv[]) {

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

