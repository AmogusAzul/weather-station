#include <WiFiClient.h>
#include <ESP8266WiFi.h>

#include <tuple>
#include <string>

#include "../communication-types/struct-likes/struct-likes.h"
#include "../common/types/string-bytes-types/string-bytes-types.h"
#include "../communication-types/byte-types/communication-bytes-types.h"

class NetworkClient {
private:
    NetworkData networkData;  // Holds WiFi and server information

    std::string ip;
    std::string ssid;
    uint16_t port;

    WiFiClient client;        // ESP8266 WiFi client for TCP communication
    bool forcedBackup = false;  // Indicates if the system should store data instead of sending
    mutable bool lastConnection = false;  // Tracks if the connection was just recovered

public:    // Constructor that sets up the WiFi connection
    NetworkClient(std::string ip, std::string ssid, uint8_t port) : ip(ip), ssid(ssid), port(port) {
        connectWiFi();
    }

    // Connect to the WiFi network
    void connectWiFi() {
        WiFi.begin(ssid.c_str(), password.c_str());  // Use .c_str() for std::string
        while (WiFi.status() != WL_CONNECTED) {
          yield();
        }
    }


    // Check if we can send data (WiFi and TCP connected) returns if can be connected and if it just recovered
    std::tuple<bool, bool> canSend() {
        if (forcedBackup) {
            return std::make_tuple(false, false);
        }

        // Ping the server to check connectivity
        if (client.connect(ip.c_str(), port)) {
            client.stop();  // Close connection after ping

            auto result = std::make_tuple(true, lastConnection==false);
            lastConnection = true;  // Mark that connection was restored
            return result;
        }

        lastConnection = false;
        return std::make_tuple(false, false);
    }

    // Force the system into backup mode (e.g., after a server error)
    void forceBackup() {
        forcedBackup = true;
    }

    bool send(const Packet& packet, Packet& serverResponsePacket) {
        if (!client.connect(networkData.getIp().getString().c_str(), static_cast<uint16_t>(networkData.getPort().getNumber()))) {
            return false;  // Unable to connect
        }

        // Send the packet to the server
        client.write(packet.getBytes(), packet.getLength());

        // Wait for the server's response
        while (!client.available()) {
          yield();
        }

        // Dynamically read the response length
        std::vector<uint8_t> responseBuffer;

        // Read all available data from the server
        while (client.available()) {
            uint8_t byte = client.read();  // Read one byte at a time
            responseBuffer.push_back(byte);  // Append byte to the buffer
        }

        // If there's any response data, use the Packet class to handle it
        if (!responseBuffer.empty()) {

            
            // Extract the header bytes from the responseBuffer
            std::vector<uint8_t> headerBytes(responseBuffer.begin(), responseBuffer.begin() + getFieldEnum(HeaderField::COUNT));

            // Set the header for the serverResponsePacket
            HeaderBytes header(headerBytes.data());  // Assuming HeaderBytes has a constructor that takes a byte array
            serverResponsePacket.setHeader(header);
        } else {
            client.stop();
            return false;  // No response received
        }

        client.stop();  // Close the connection after sending and receiving

        return true;  // Return true if the packet was successfully sent and response received
    }
};